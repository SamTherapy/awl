// SPDX-License-Identifier: BSD-3-Clause

package resolvers

import (
	"crypto/tls"
	"fmt"
	"net"

	"git.froth.zone/sam/awl/pkg/util"
	"github.com/miekg/dns"
)

// StandardResolver is for UDP/TCP resolvers.
type StandardResolver struct {
	opts util.Options
}

var _ Resolver = (*StandardResolver)(nil)

// LookUp performs a DNS query.
func (resolver *StandardResolver) LookUp(msg *dns.Msg) (util.Response, error) {
	var (
		resp util.Response
		err  error
	)

	dnsClient := new(dns.Client)
	dnsClient.Dialer = &net.Dialer{
		Timeout: resolver.opts.Request.Timeout,
	}

	if resolver.opts.TCP || resolver.opts.TLS {
		dnsClient.Net = tcp
	} else {
		dnsClient.Net = udp
	}

	switch {
	case resolver.opts.IPv4:
		dnsClient.Net += "4"
	case resolver.opts.IPv6:
		dnsClient.Net += "6"
	}

	if resolver.opts.TLS {
		dnsClient.Net += "-tls"
		dnsClient.TLSConfig = &tls.Config{
			//nolint:gosec // This is intentional if the user requests it
			InsecureSkipVerify: resolver.opts.TLSNoVerify,
			ServerName:         resolver.opts.TLSHost,
		}
	}

	resolver.opts.Logger.Info("Using", dnsClient.Net, "for making the request")

	resp.DNS, resp.RTT, err = dnsClient.Exchange(msg, resolver.opts.Request.Server)
	if err != nil {
		return util.Response{}, fmt.Errorf("standard: DNS exchange: %w", err)
	}

	switch dns.RcodeToString[resp.DNS.MsgHdr.Rcode] {
	case "BADCOOKIE":
		if !resolver.opts.BadCookie {
			fmt.Printf(";; BADCOOKIE, retrying.\n\n")

			msg.Extra = resp.DNS.Extra

			resp.DNS, resp.RTT, err = dnsClient.Exchange(msg, resolver.opts.Request.Server)

			if err != nil {
				return util.Response{}, fmt.Errorf("badcookie: DNS exchange: %w", err)
			}
		}

	case "NOERR":
		break
	}

	resolver.opts.Logger.Info("Request successful")

	if resp.DNS.MsgHdr.Truncated && !resolver.opts.Truncate {
		fmt.Printf(";; Truncated, retrying with TCP\n\n")

		dnsClient.Net = tcp

		switch {
		case resolver.opts.IPv4:
			dnsClient.Net += "4"
		case resolver.opts.IPv6:
			dnsClient.Net += "6"
		}

		resp.DNS, resp.RTT, err = dnsClient.Exchange(msg, resolver.opts.Request.Server)
	}

	if err != nil {
		return util.Response{}, fmt.Errorf("standard: DNS exchange: %w", err)
	}

	return resp, nil
}
