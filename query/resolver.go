// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"net"
	"strconv"
	"strings"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/internal/helpers"
	"github.com/miekg/dns"
)

type Resolver interface {
	LookUp(*dns.Msg) (helpers.Response, error)
}

func LoadResolver(opts cli.Options) (Resolver, error) {
	switch {
	case opts.HTTPS:
		opts.Logger.Info("loading DNS-over-HTTPS resolver")
		if !strings.HasPrefix(opts.Request.Server, "https://") {
			opts.Request.Server = "https://" + opts.Request.Server
		}
		return &HTTPSResolver{
			opts: opts,
		}, nil
	case opts.QUIC:
		opts.Logger.Info("loading DNS-over-QUIC resolver")
		opts.Request.Server = net.JoinHostPort(opts.Request.Server, strconv.Itoa(opts.Port))
		return &QUICResolver{
			opts: opts,
		}, nil
	case opts.DNSCrypt:
		opts.Logger.Info("loading DNSCrypt resolver")
		if !strings.HasPrefix(opts.Request.Server, "sdns://") {
			opts.Request.Server = "sdns://" + opts.Request.Server
		}
		return &DNSCryptResolver{
			opts: opts,
		}, nil
	default:
		opts.Logger.Info("loading standard/DNS-over-TLS resolver")
		opts.Request.Server = net.JoinHostPort(opts.Request.Server, strconv.Itoa(opts.Port))
		return &StandardResolver{
			opts: opts,
		}, nil
	}
}
