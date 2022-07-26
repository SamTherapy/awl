// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"fmt"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/internal/helpers"

	"github.com/miekg/dns"
)

func CreateQuery(opts cli.Options) (helpers.Response, error) {
	var res helpers.Response
	res.DNS = new(dns.Msg)
	res.DNS.SetQuestion(opts.Request.Name, opts.Request.Type)
	res.DNS.Question[0].Qclass = opts.Request.Class

	res.DNS.MsgHdr.Response = opts.QR
	res.DNS.MsgHdr.Authoritative = opts.AA
	res.DNS.MsgHdr.Truncated = opts.TC
	res.DNS.MsgHdr.RecursionDesired = opts.RD
	res.DNS.MsgHdr.RecursionAvailable = opts.RA
	res.DNS.MsgHdr.Zero = opts.Z
	res.DNS.MsgHdr.AuthenticatedData = opts.AD
	res.DNS.MsgHdr.CheckingDisabled = opts.CD

	if opts.DNSSEC {
		res.DNS.SetEdns0(1232, true)
	}

	opts.Logger.Debug(fmt.Sprintf("%+v", res))

	resolver, err := LoadResolver(opts)
	if err != nil {
		return helpers.Response{}, err
	}
	opts.Logger.Info("Query successfully loaded")

	return resolver.LookUp(res.DNS)
}
