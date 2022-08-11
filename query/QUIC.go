// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"crypto/tls"
	"fmt"
	"io"
	"time"

	"git.froth.zone/sam/awl/util"
	"github.com/lucas-clemente/quic-go"
	"github.com/miekg/dns"
)

// QUICResolver is for DNS-over-QUIC queries.
type QUICResolver struct {
	opts util.Options
}

var _ Resolver = (*QUICResolver)(nil)

// LookUp performs a DNS query.
func (r *QUICResolver) LookUp(msg *dns.Msg) (util.Response, error) {
	var resp util.Response

	tls := &tls.Config{
		//nolint:gosec // This is intentional if the user requests it
		InsecureSkipVerify: r.opts.TLSNoVerify,
		ServerName:         r.opts.TLSHost,
		MinVersion:         tls.VersionTLS12,
		NextProtos:         []string{"doq"},
	}

	conf := new(quic.Config)
	conf.HandshakeIdleTimeout = r.opts.Request.Timeout

	r.opts.Logger.Debug("quic: making query")

	connection, err := quic.DialAddr(r.opts.Request.Server, tls, conf)
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: dial: %w", err)
	}

	r.opts.Logger.Debug("quic: packing query")

	// Compress request to over-the-wire
	buf, err := msg.Pack()
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: pack: %w", err)
	}

	t := time.Now()

	r.opts.Logger.Debug("quic: creating stream")

	stream, err := connection.OpenStream()
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: quic stream creation: %w", err)
	}

	r.opts.Logger.Debug("quic: writing to stream")

	_, err = stream.Write(buf)
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: quic stream write: %w", err)
	}

	r.opts.Logger.Debug("quic: reading stream")

	fullRes, err := io.ReadAll(stream)
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: quic stream read: %w", err)
	}

	resp.RTT = time.Since(t)

	r.opts.Logger.Debug("quic: closing connection")
	// Close with error: no error
	err = connection.CloseWithError(0, "")
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: quic connection close: %w", err)
	}

	r.opts.Logger.Debug("quic: closing stream")

	err = stream.Close()
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: quic stream close: %w", err)
	}

	resp.DNS = &dns.Msg{}

	r.opts.Logger.Debug("quic: unpacking response")

	err = resp.DNS.Unpack(fullRes)
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: unpack: %w", err)
	}

	return resp, nil
}
