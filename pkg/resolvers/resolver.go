// SPDX-License-Identifier: BSD-3-Clause

package resolvers

import (
	"net"
	"strconv"
	"strings"

	"dns.froth.zone/awl/pkg/util"
	"github.com/miekg/dns"
)

const (
	tcp = "tcp"
	udp = "udp"
)

// Resolver is the main resolver interface.
type Resolver interface {
	LookUp(*dns.Msg) (util.Response, error)
}

// LoadResolver loads the respective resolver for performing a DNS query.
func LoadResolver(opts *util.Options) (resolver Resolver, err error) {
	switch {
	case opts.HTTPS:
		opts.Logger.Info("loading DNS-over-HTTPS resolver")

		if !strings.HasPrefix(opts.Request.Server, "https://") {
			opts.Request.Server = "https://" + opts.Request.Server
		}

		// Make sure that the endpoint is defaulted to /dns-query
		if !strings.HasSuffix(opts.Request.Server, opts.HTTPSOptions.Endpoint) {
			opts.Request.Server += opts.HTTPSOptions.Endpoint
		}

		resolver = &HTTPSResolver{
			opts: opts,
		}

		return
	case opts.QUIC:
		opts.Logger.Info("loading DNS-over-QUIC resolver")

		if !strings.HasSuffix(opts.Request.Server, ":"+strconv.Itoa(opts.Request.Port)) {
			opts.Request.Server = net.JoinHostPort(opts.Request.Server, strconv.Itoa(opts.Request.Port))
		}

		resolver = &QUICResolver{
			opts: opts,
		}

		return
	case opts.DNSCrypt:
		opts.Logger.Info("loading DNSCrypt resolver")

		if !strings.HasPrefix(opts.Request.Server, "sdns://") {
			opts.Request.Server = "sdns://" + opts.Request.Server
		}

		resolver = &DNSCryptResolver{
			opts: opts,
		}

		return
	default:
		opts.Logger.Info("loading standard/DNS-over-TLS resolver")

		if !strings.HasSuffix(opts.Request.Server, ":"+strconv.Itoa(opts.Request.Port)) {
			opts.Request.Server = net.JoinHostPort(opts.Request.Server, strconv.Itoa(opts.Request.Port))
		}

		resolver = &StandardResolver{
			opts: opts,
		}

		return
	}
}
