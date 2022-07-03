// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"crypto/tls"
	"io"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/miekg/dns"
)

type QUICResolver struct {
	server string
	opts   Options
}

func (r *QUICResolver) LookUp(msg *dns.Msg) (*dns.Msg, time.Duration, error) {
	var resp Response
	tls := &tls.Config{
		NextProtos: []string{"doq"},
	}
	r.opts.Logger.Debug("making DoQ request")
	connection, err := quic.DialAddr(r.server, tls, nil)
	if err != nil {
		return nil, 0, err
	}

	// Compress request to over-the-wire
	buf, err := msg.Pack()
	if err != nil {
		return nil, 0, err
	}
	t := time.Now()
	stream, err := connection.OpenStream()
	if err != nil {
		return nil, 0, err
	}
	_, err = stream.Write(buf)
	if err != nil {
		return nil, 0, err
	}

	fullRes, err := io.ReadAll(stream)
	if err != nil {
		return nil, 0, err
	}
	resp.Answers.RTT = time.Since(t)

	// Close with error: no error
	err = connection.CloseWithError(0, "")
	if err != nil {
		return nil, 0, err
	}

	err = stream.Close()
	if err != nil {
		return nil, 0, err
	}

	resp.DNS = dns.Msg{}
	r.opts.Logger.Debug("unpacking DoQ response")
	err = resp.DNS.Unpack(fullRes)
	if err != nil {
		return nil, 0, err
	}
	return &resp.DNS, resp.Answers.RTT, nil
}
