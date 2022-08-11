// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"fmt"
	"time"

	"git.froth.zone/sam/awl/util"
	"github.com/ameshkov/dnscrypt/v2"
	"github.com/miekg/dns"
)

// DNSCryptResolver is for making DNSCrypt queries.
type DNSCryptResolver struct {
	opts util.Options
}

var _ Resolver = (*DNSCryptResolver)(nil)

// LookUp performs a DNS query.
func (r *DNSCryptResolver) LookUp(msg *dns.Msg) (util.Response, error) {
	client := dnscrypt.Client{
		Timeout: r.opts.Request.Timeout,
		UDPSize: 1232,
	}

	if r.opts.TCP || r.opts.TLS {
		client.Net = tcp
	} else {
		client.Net = udp
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
		return util.Response{}, fmt.Errorf("dnscrypt: dial: %w", err)
	}

	now := time.Now()
	res, err := client.Exchange(msg, resolverInf)
	rtt := time.Since(now)

	if err != nil {
		return util.Response{}, fmt.Errorf("dnscrypt: exchange: %w", err)
	}

	r.opts.Logger.Info("Request successful")

	return util.Response{
		DNS: res,
		RTT: rtt,
	}, nil
}
