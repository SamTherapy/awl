package query

import (
	"net"
	"strconv"
	"strings"
	"time"

	"git.froth.zone/sam/awl/logawl"
	"github.com/miekg/dns"
)

// represent all CLI flags
type Options struct {
	Logger *logawl.Logger

	Port     int
	IPv4     bool
	IPv6     bool
	DNSSEC   bool
	Short    bool
	TCP      bool
	TLS      bool
	HTTPS    bool
	QUIC     bool
	Truncate bool
	AA       bool
	TC       bool
	Z        bool
	CD       bool
	NoRD     bool
	NoRA     bool
	Reverse  bool
	Debug    bool
	Answers  Answers
}
type Response struct {
	Answers Answers `json:"Response"` // These be DNS query answers
	DNS     dns.Msg
}

// The Answers struct is the basic structure of a DNS request
// to be returned to the user upon making a request
type Answers struct {
	Server  string `json:"Server"` // The server to make the DNS request from
	DNS     *dns.Msg
	Request uint16        `json:"Request"` // The type of request
	Name    string        `json:"Name"`    // The domain name to make a DNS request for
	RTT     time.Duration `json:"RTT"`     // The time it took to make the DNS query
}

type Resolver interface {
	LookUp(*dns.Msg) (*dns.Msg, time.Duration, error)
}

func LoadResolver(server string, opts Options) (Resolver, error) {
	if opts.HTTPS {
		opts.Logger.Debug("loading DoH resolver")
		if !strings.HasPrefix(server, "https://") {
			server = "https://" + server
		}
		return &HTTPSResolver{
			server: server,
			opts:   opts,
		}, nil
	} else if opts.QUIC {
		opts.Logger.Debug("loading DoQ resolver")
		server = net.JoinHostPort(opts.Answers.Server, strconv.Itoa(opts.Port))
		return &QUICResolver{
			server: server,
			opts:   opts,
		}, nil
	}

	return nil, nil
}
