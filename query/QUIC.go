// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"crypto/tls"
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
		return helpers.Response{}, err
	}

	// Compress request to over-the-wire
	buf, err := msg.Pack()
	if err != nil {
		return helpers.Response{}, err
	}
	t := time.Now()
	stream, err := connection.OpenStream()
	if err != nil {
		return helpers.Response{}, err
	}
	_, err = stream.Write(buf)
	if err != nil {
		return helpers.Response{}, err
	}

	fullRes, err := io.ReadAll(stream)
	if err != nil {
		return helpers.Response{}, err
	}
	resp.RTT = time.Since(t)

	// Close with error: no error
	err = connection.CloseWithError(0, "")
	if err != nil {
		return helpers.Response{}, err
	}

	err = stream.Close()
	if err != nil {
		return helpers.Response{}, err
	}

	resp.DNS = &dns.Msg{}
	r.opts.Logger.Debug("unpacking DoQ response")
	err = resp.DNS.Unpack(fullRes)
	if err != nil {
		return helpers.Response{}, err
	}
	return resp, nil
}
