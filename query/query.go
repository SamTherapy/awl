// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"fmt"
	"strconv"

	"git.froth.zone/sam/awl/util"
	"github.com/dchest/uniuri"
	"github.com/miekg/dns"
)

const (
	tcp = "tcp"
	udp = "udp"
)

// CreateQuery creates a DNS query from the options given.
// It sets query flags and EDNS flags from the respective options.
func CreateQuery(opts util.Options) (util.Response, error) {
	req := new(dns.Msg)
	req.SetQuestion(opts.Request.Name, opts.Request.Type)
	req.Question[0].Qclass = opts.Request.Class

	// Set standard flags
	req.MsgHdr.Response = opts.QR
	req.MsgHdr.Authoritative = opts.AA
	req.MsgHdr.Truncated = opts.TC
	req.MsgHdr.RecursionDesired = opts.RD
	req.MsgHdr.RecursionAvailable = opts.RA
	req.MsgHdr.Zero = opts.Z
	req.MsgHdr.AuthenticatedData = opts.AD
	req.MsgHdr.CheckingDisabled = opts.CD

	// EDNS time :)
	if opts.EDNS.EnableEDNS {
		o := new(dns.OPT)
		o.Hdr.Name = "."
		o.Hdr.Rrtype = dns.TypeOPT

		o.SetVersion(opts.EDNS.Version)

		if opts.EDNS.Cookie {
			e := new(dns.EDNS0_COOKIE)
			e.Code = dns.EDNS0COOKIE
			e.Cookie = uniuri.NewLenChars(8, []byte("1234567890abcdef"))
			o.Option = append(o.Option, e)

			opts.Logger.Info("Setting EDNS cookie to", e.Cookie)
		}

		if opts.EDNS.Expire {
			o.Option = append(o.Option, new(dns.EDNS0_EXPIRE))

			opts.Logger.Info("Setting EDNS Expire option")
		}

		if opts.EDNS.KeepOpen {
			o.Option = append(o.Option, new(dns.EDNS0_TCP_KEEPALIVE))

			opts.Logger.Info("Setting EDNS TCP Keepalive option")
		}

		if opts.EDNS.Nsid {
			o.Option = append(o.Option, new(dns.EDNS0_NSID))

			opts.Logger.Info("Setting EDNS NSID option")
		}

		if opts.EDNS.Padding {
			o.Option = append(o.Option, new(dns.EDNS0_PADDING))

			opts.Logger.Info("Setting EDNS padding")
		}

		o.SetUDPSize(opts.BufSize)

		opts.Logger.Info("EDNS UDP buffer set to", opts.BufSize)

		o.SetZ(opts.EDNS.ZFlag)

		opts.Logger.Info("EDNS Z flag set to", opts.EDNS.ZFlag)

		if opts.EDNS.DNSSEC {
			o.SetDo()

			opts.Logger.Info("EDNS DNSSEC OK set")
		}

		if opts.EDNS.Subnet.Address != nil {
			o.Option = append(o.Option, &opts.EDNS.Subnet)
		}

		req.Extra = append(req.Extra, o)
	} else if opts.EDNS.DNSSEC {
		req.SetEdns0(1232, true)
		opts.Logger.Warn("DNSSEC implies EDNS, EDNS enabled")
		opts.Logger.Info("DNSSEC enabled, UDP buffer set to 1232")
	}

	opts.Logger.Debug(req)

	if !opts.Short {
		if opts.ShowQuery {
			opts.Logger.Info("Printing constructed query")

			var (
				str string
				err error
			)

			if opts.JSON || opts.XML || opts.YAML {
				str, err = PrintSpecial(req, opts)
				if err != nil {
					return util.Response{}, err
				}
			} else {
				temp := opts.Display.Statistics
				opts.Display.Statistics = false
				str = ToString(util.Response{
					DNS: req,
					RTT: 0,
				}, opts)
				opts.Display.Statistics = temp
				str += "\n;; QUERY SIZE: " + strconv.Itoa(req.Len()) + "\n"
			}

			fmt.Println(str)

			opts.ShowQuery = false
		}
	}

	resolver, err := LoadResolver(opts)
	if err != nil {
		return util.Response{}, err
	}

	opts.Logger.Info("Query successfully loaded")

	//nolint:wrapcheck // Error wrapping not needed here
	return resolver.LookUp(req)
}
