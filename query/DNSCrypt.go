// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"time"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/internal/helpers"

	"github.com/ameshkov/dnscrypt/v2"
	"github.com/miekg/dns"
)

type DNSCryptResolver struct {
	opts cli.Options
}

func (r *DNSCryptResolver) LookUp(msg *dns.Msg) (helpers.Response, error) {

	client := dnscrypt.Client{
		Timeout: r.opts.Request.Timeout,
		UDPSize: 1232,
	}

	if r.opts.TCP || r.opts.TLS {
		client.Net = "tcp"
	} else {
		client.Net = "udp"
	}

	switch {
	case r.opts.IPv4:
		client.Net += "4"
	case r.opts.IPv6:
		client.Net += "6"
	}
	r.opts.Logger.Debug("Using", client.Net, "for making the request")

	resolverInf, err := client.Dial(r.opts.Request.Server)
	if err != nil {
		return helpers.Response{}, err
	}

	now := time.Now()
	res, err := client.Exchange(msg, resolverInf)
	rtt := time.Since(now)

	if err != nil {
		return helpers.Response{}, err
	}
	r.opts.Logger.Info("Request successful")

	return helpers.Response{
		DNS: res,
		RTT: rtt,
	}, nil
}
