// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"crypto/tls"
	"fmt"
	"net"

	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
)

// StandardResolver is for UDP/TCP resolvers.
type StandardResolver struct {
	opts util.Options
}

var _ Resolver = (*StandardResolver)(nil)

// LookUp performs a DNS query.
func (r *StandardResolver) LookUp(msg *dns.Msg) (util.Response, error) {
	var (
		resp util.Response
		err  error
	)

	dnsClient := new(dns.Client)
	dnsClient.Dialer = &net.Dialer{
		Timeout: r.opts.Request.Timeout,
	}

	if r.opts.TCP || r.opts.TLS {
		dnsClient.Net = tcp
	} else {
		dnsClient.Net = udp
	}

	switch {
	case r.opts.IPv4:
		dnsClient.Net += "4"
	case r.opts.IPv6:
		dnsClient.Net += "6"
	}

	if r.opts.TLS {
		dnsClient.Net += "-tls"
		dnsClient.TLSConfig = &tls.Config{
			//nolint:gosec // This is intentional if the user requests it
			InsecureSkipVerify: r.opts.TLSNoVerify,
			ServerName:         r.opts.TLSHost,
		}
	}

	r.opts.Logger.Debug("Using", dnsClient.Net, "for making the request")

	resp.DNS, resp.RTT, err = dnsClient.Exchange(msg, r.opts.Request.Server)
	if err != nil {
		return util.Response{}, fmt.Errorf("standard: DNS exchange: %w", err)
	}

	r.opts.Logger.Info("Request successful")

	if resp.DNS.MsgHdr.Truncated && !r.opts.Truncate {
		fmt.Printf(";; Truncated, retrying with TCP\n\n")

		dnsClient.Net = tcp

		switch {
		case r.opts.IPv4:
			dnsClient.Net += "4"
		case r.opts.IPv6:
			dnsClient.Net += "6"
		}

		resp.DNS, resp.RTT, err = dnsClient.Exchange(msg, r.opts.Request.Server)
	}

	if err != nil {
		return util.Response{}, fmt.Errorf("standard: DNS exchange: %w", err)
	}

	return resp, nil
}
