// SPDX-License-Identifier: BSD-3-Clause

package resolvers

import (
	"fmt"
	"time"

	"dns.froth.zone/awl/pkg/util"
	"dns.froth.zone/dnscrypt"
	"github.com/miekg/dns"
)

// DNSCryptResolver is for making DNSCrypt queries.
type DNSCryptResolver struct {
	opts *util.Options
}

var _ Resolver = (*DNSCryptResolver)(nil)

// LookUp performs a DNS query.
func (resolver *DNSCryptResolver) LookUp(msg *dns.Msg) (resp util.Response, err error) {
	client := dnscrypt.Client{
		Timeout: resolver.opts.Request.Timeout,
		UDPSize: 1232,
	}

	if resolver.opts.TCP || resolver.opts.TLS {
		client.Net = tcp
	} else {
		client.Net = udp
	}

	switch {
	case resolver.opts.IPv4:
		client.Net += "4"
	case resolver.opts.IPv6:
		client.Net += "6"
	}

	resolver.opts.Logger.Debug("Using", client.Net, "for making the request")

	resolverInf, err := client.Dial(resolver.opts.Request.Server)
	if err != nil {
		return resp, fmt.Errorf("dnscrypt: dial: %w", err)
	}

	now := time.Now()
	res, err := client.Exchange(msg, resolverInf)
	rtt := time.Since(now)

	if err != nil {
		return resp, fmt.Errorf("dnscrypt: exchange: %w", err)
	}

	resp = util.Response{
		DNS: res,
		RTT: rtt,
	}

	resolver.opts.Logger.Info("Request successful")

	return
}
