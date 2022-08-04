// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"fmt"
	"net"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/internal/helpers"
	"github.com/miekg/dns"
)

type StandardResolver struct {
	opts cli.Options
}

// LookUp performs a DNS query
func (r *StandardResolver) LookUp(msg *dns.Msg) (helpers.Response, error) {
	var (
		resp helpers.Response
		err  error
	)
	dnsClient := new(dns.Client)
	dnsClient.Dialer = &net.Dialer{
		Timeout: r.opts.Request.Timeout,
	}
	if r.opts.TCP || r.opts.TLS {
		dnsClient.Net = "tcp"
	} else {
		dnsClient.Net = "udp"
	}

	switch {
	case r.opts.IPv4:
		dnsClient.Net += "4"
	case r.opts.IPv6:
		dnsClient.Net += "6"
	}

	if r.opts.TLS {
		dnsClient.Net += "-tls"
	}
	r.opts.Logger.Debug("Using", dnsClient.Net, "for making the request")

	resp.DNS, resp.RTT, err = dnsClient.Exchange(msg, r.opts.Request.Server)
	if err != nil {
		return helpers.Response{}, fmt.Errorf("standard: DNS exchange error: %w", err)
	}
	r.opts.Logger.Info("Request successful")

	if resp.DNS.MsgHdr.Truncated && !r.opts.Truncate {
		fmt.Printf(";; Truncated, retrying with TCP\n\n")
		dnsClient.Net = "tcp"
		switch {
		case r.opts.IPv4:
			dnsClient.Net += "4"
		case r.opts.IPv6:
			dnsClient.Net += "6"
		}
		resp.DNS, resp.RTT, err = dnsClient.Exchange(msg, r.opts.Request.Server)
	}
	if err != nil {
		return helpers.Response{}, fmt.Errorf("standard: DNS exchange error: %w", err)
	}

	return resp, nil
}
