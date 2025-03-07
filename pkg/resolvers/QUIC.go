// SPDX-License-Identifier: BSD-3-Clause

package resolvers

import (
	"context"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"

	"dns.froth.zone/awl/pkg/util"
	"github.com/miekg/dns"
	"github.com/quic-go/quic-go"
)

// QUICResolver is for DNS-over-QUIC queries.
type QUICResolver struct {
	opts *util.Options
}

var _ Resolver = (*QUICResolver)(nil)

// LookUp performs a DNS query.
func (resolver *QUICResolver) LookUp(msg *dns.Msg) (resp util.Response, err error) {
	tls := &tls.Config{
		//nolint:gosec // This is intentional if the user requests it
		InsecureSkipVerify: resolver.opts.TLSNoVerify,
		ServerName:         resolver.opts.TLSHost,
		MinVersion:         tls.VersionTLS12,
		NextProtos:         []string{"doq"},
	}

	// Make sure that TLSHost is ALWAYS set
	if resolver.opts.TLSHost == "" {
		tls.ServerName = strings.Split(resolver.opts.Request.Server, ":")[0]
	}

	conf := new(quic.Config)
	conf.HandshakeIdleTimeout = resolver.opts.Request.Timeout

	resolver.opts.Logger.Debug("quic: making query")

	ctx, cancel := context.WithTimeout(context.Background(), resolver.opts.Request.Timeout)
	defer cancel()

	connection, err := quic.DialAddr(ctx, resolver.opts.Request.Server, tls, conf)
	if err != nil {
		return resp, fmt.Errorf("doq: dial: %w", err)
	}

	resolver.opts.Logger.Debug("quic: packing query")

	msg.Id = 0
	// Compress request to over-the-wire
	buf, err := msg.Pack()
	if err != nil {
		return resp, fmt.Errorf("doq: pack: %w", err)
	}

	t := time.Now()

	resolver.opts.Logger.Debug("quic: creating stream")

	stream, err := connection.OpenStream()
	if err != nil {
		return resp, fmt.Errorf("doq: quic stream creation: %w", err)
	}

	resolver.opts.Logger.Debug("quic: writing to stream")

	_, err = stream.Write(rfc9250prefix(buf))
	if err != nil {
		return resp, fmt.Errorf("doq: quic stream write: %w", err)
	}

	err = stream.Close()
	if err != nil {
		return resp, fmt.Errorf("doq: quic stream close: %w", err)
	}

	resolver.opts.Logger.Debug("quic: reading stream")

	fullRes, err := io.ReadAll(stream)
	if err != nil {
		return resp, fmt.Errorf("doq: quic stream read: %w", err)
	}

	resp.RTT = time.Since(t)

	resolver.opts.Logger.Debug("quic: closing connection")
	// Close with error: no error
	err = connection.CloseWithError(0, "")
	if err != nil {
		return resp, fmt.Errorf("doq: quic connection close: %w", err)
	}

	resolver.opts.Logger.Debug("quic: closing stream")

	resp.DNS = &dns.Msg{}

	resolver.opts.Logger.Debug("quic: unpacking response")

	// Unpack response and lop off the first two bytes (RFC 9250 moment)
	err = resp.DNS.Unpack(fullRes[2:])
	if err != nil {
		return resp, fmt.Errorf("doq: unpack: %w", err)
	}

	return
}

// rfc9250prefix adds a two-byte prefix to the input data as per RFC 9250.
func rfc9250prefix(in []byte) []byte {
	out := make([]byte, 2+len(in))
	binary.BigEndian.PutUint16(out, uint16(len(in)))
	copy(out[2:], in)
	return out
}
