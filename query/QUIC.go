package query

import (
	"crypto/tls"
	"io"
	"time"

	"github.com/lucas-clemente/quic-go"
	"github.com/miekg/dns"
)

// Resolve DNS over QUIC, the hip new standard (for privacy I think, IDK)
func ResolveQUIC(msg *dns.Msg, server string) (*dns.Msg, time.Duration, error) {
	tls := &tls.Config{
		NextProtos: []string{"doq"},
	}
	connection, err := quic.DialAddr(server, tls, nil)
	if err != nil {
		return nil, 0, err
	}

	// Close with error: no error
	defer connection.CloseWithError(0, "")

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
	rtt := time.Since(t)

	err = stream.Close()
	if err != nil {
		return nil, 0, err
	}

	response := dns.Msg{}
	err = response.Unpack(fullRes)
	if err != nil {
		return nil, 0, err
	}

	return &response, rtt, nil
}
