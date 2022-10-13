// SPDX-License-Identifier: BSD-3-Clause

package resolvers

import (
	"crypto/tls"
	"fmt"
	"io"
	"time"

	"git.froth.zone/sam/awl/pkg/util"
	"github.com/lucas-clemente/quic-go"
	"github.com/miekg/dns"
)

// QUICResolver is for DNS-over-QUIC queries.
type QUICResolver struct {
	opts *util.Options
}

var _ Resolver = (*QUICResolver)(nil)

// LookUp performs a DNS query.
func (resolver *QUICResolver) LookUp(msg *dns.Msg) (util.Response, error) {
	var resp util.Response

	tls := &tls.Config{
		//nolint:gosec // This is intentional if the user requests it
		InsecureSkipVerify: resolver.opts.TLSNoVerify,
		ServerName:         resolver.opts.TLSHost,
		MinVersion:         tls.VersionTLS12,
		NextProtos:         []string{"doq"},
	}

	conf := new(quic.Config)
	conf.HandshakeIdleTimeout = resolver.opts.Request.Timeout

	resolver.opts.Logger.Debug("quic: making query")

	connection, err := quic.DialAddr(resolver.opts.Request.Server, tls, conf)
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: dial: %w", err)
	}

	resolver.opts.Logger.Debug("quic: packing query")

	// Compress request to over-the-wire
	buf, err := msg.Pack()
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: pack: %w", err)
	}

	t := time.Now()

	resolver.opts.Logger.Debug("quic: creating stream")

	stream, err := connection.OpenStream()
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: quic stream creation: %w", err)
	}

	resolver.opts.Logger.Debug("quic: writing to stream")

	_, err = stream.Write(buf)
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: quic stream write: %w", err)
	}

	resolver.opts.Logger.Debug("quic: reading stream")

	fullRes, err := io.ReadAll(stream)
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: quic stream read: %w", err)
	}

	resp.RTT = time.Since(t)

	resolver.opts.Logger.Debug("quic: closing connection")
	// Close with error: no error
	err = connection.CloseWithError(0, "")
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: quic connection close: %w", err)
	}

	resolver.opts.Logger.Debug("quic: closing stream")

	err = stream.Close()
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: quic stream close: %w", err)
	}

	resp.DNS = &dns.Msg{}

	resolver.opts.Logger.Debug("quic: unpacking response")

	err = resp.DNS.Unpack(fullRes)
	if err != nil {
		return util.Response{}, fmt.Errorf("doq: unpack: %w", err)
	}

	return resp, nil
}
