// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"crypto/tls"
	"fmt"
	"io"
	"time"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/internal/helpers"
	"github.com/lucas-clemente/quic-go"
	"github.com/miekg/dns"
)

type QUICResolver struct {
	opts cli.Options
}

// LookUp performs a DNS query.
func (r *QUICResolver) LookUp(msg *dns.Msg) (helpers.Response, error) {
	var resp helpers.Response
	tls := &tls.Config{
		MinVersion: tls.VersionTLS12,
		NextProtos: []string{"doq"},
	}

	conf := new(quic.Config)
	conf.HandshakeIdleTimeout = r.opts.Request.Timeout

	r.opts.Logger.Debug("making DoQ request")
	connection, err := quic.DialAddr(r.opts.Request.Server, tls, conf)
	if err != nil {
		return helpers.Response{}, fmt.Errorf("doq: dial error: %w", err)
	}

	// Compress request to over-the-wire
	buf, err := msg.Pack()
	if err != nil {
		return helpers.Response{}, fmt.Errorf("doq: pack error: %w", err)
	}
	t := time.Now()
	stream, err := connection.OpenStream()
	if err != nil {
		return helpers.Response{}, fmt.Errorf("doq: quic stream creation error: %w", err)
	}
	_, err = stream.Write(buf)
	if err != nil {
		return helpers.Response{}, fmt.Errorf("doq: quic stream write error: %w", err)
	}

	fullRes, err := io.ReadAll(stream)
	if err != nil {
		return helpers.Response{}, fmt.Errorf("doq: quic stream read error: %w", err)
	}
	resp.RTT = time.Since(t)

	// Close with error: no error
	err = connection.CloseWithError(0, "")
	if err != nil {
		return helpers.Response{}, fmt.Errorf("doq: quic connection close error: %w", err)
	}

	err = stream.Close()
	if err != nil {
		return helpers.Response{}, fmt.Errorf("doq: quic stream close error: %w", err)
	}

	resp.DNS = &dns.Msg{}
	r.opts.Logger.Debug("unpacking DoQ response")
	err = resp.DNS.Unpack(fullRes)
	if err != nil {
		return helpers.Response{}, fmt.Errorf("doq: upack error: %w", err)
	}
	return resp, nil
}
