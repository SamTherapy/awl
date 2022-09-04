// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"strconv"
	"strings"
	"time"

	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
)

// Message is for overall DNS responses.
type Message struct {
	Question []Question `json:"question,omitempty" xml:"question,omitempty" yaml:",omitempty"`
	Answer   []Answer   `json:"answer,omitempty" xml:"answer,omitempty" yaml:",omitempty"`
	Ns       []Answer   `json:"ns,omitempty" xml:"ns,omitempty" yaml:",omitempty"`
	Extra    []Answer   `json:"extra,omitempty" xml:"extra,omitempty" yaml:",omitempty"`
	Header   dns.MsgHdr `json:"header,omitempty" xml:"header,omitempty" yaml:",omitempty"`
}

// Question is a DNS Query.
type Question struct {
	Name  string `json:"name,omitempty" xml:"name,omitempty" yaml:",omitempty"`
	Type  string `json:"type,omitempty" xml:"type,omitempty" yaml:",omitempty"`
	Class string `json:"class,omitempty" xml:"class,omitempty" yaml:",omitempty"`
}

// RRHeader is for DNS Resource Headers.
type RRHeader struct {
	Name     string `json:"name,omitempty" xml:"name,omitempty" yaml:",omitempty"`
	Type     string `json:"type,omitempty" xml:"type,omitempty" yaml:",omitempty"`
	Class    string `json:"class,omitempty" xml:"class,omitempty" yaml:",omitempty"`
	TTL      string `json:"ttl,omitempty" xml:"ttl,omitempty" yaml:",omitempty"`
	Rdlength uint16 `json:"-" xml:"-" yaml:"-"`
}

// Answer is for a DNS Response.
type Answer struct {
	Value    string `json:"response,omitempty" xml:"response,omitempty" yaml:"response,omitempty"`
	RRHeader `json:"header,omitempty" xml:"header,omitempty" yaml:"header,omitempty"`
}

// ToString turns the response into something that looks a lot like dig
//
// Much of this is taken from https://github.com/miekg/dns/blob/master/msg.go#L900
func ToString(res util.Response, opts util.Options) string {
	if res.DNS == nil {
		return "<nil> MsgHdr"
	}

	var (
		s   string
		opt *dns.OPT
	)

	if !opts.Short {
		if opts.Display.Comments {
			s += res.DNS.MsgHdr.String() + " "
			s += "QUERY: " + strconv.Itoa(len(res.DNS.Question)) + ", "
			s += "ANSWER: " + strconv.Itoa(len(res.DNS.Answer)) + ", "
			s += "AUTHORITY: " + strconv.Itoa(len(res.DNS.Ns)) + ", "
			s += "ADDITIONAL: " + strconv.Itoa(len(res.DNS.Extra)) + "\n"
			opt = res.DNS.IsEdns0()

			if opt != nil && opts.Display.Opt {
				// OPT PSEUDOSECTION
				s += opt.String() + "\n"
			}
		}

		if opts.Display.Question {
			if len(res.DNS.Question) > 0 {
				if opts.Display.Comments {
					s += "\n;; QUESTION SECTION:\n"
				}

				for _, r := range res.DNS.Question {
					s += r.String() + "\n"
				}
			}
		}

		if opts.Display.Answer {
			if len(res.DNS.Answer) > 0 {
				if opts.Display.Comments {
					s += "\n;; ANSWER SECTION:\n"
				}

				for _, r := range res.DNS.Answer {
					if r != nil {
						s += r.String() + "\n"
					}
				}
			}
		}

		if opts.Display.Authority {
			if len(res.DNS.Ns) > 0 {
				if opts.Display.Comments {
					s += "\n;; AUTHORITY SECTION:\n"
				}

				for _, r := range res.DNS.Ns {
					if r != nil {
						s += r.String() + "\n"
					}
				}
			}
		}

		if opts.Display.Additional {
			if len(res.DNS.Extra) > 0 && (opt == nil || len(res.DNS.Extra) > 1) {
				if opts.Display.Comments {
					s += "\n;; ADDITIONAL SECTION:\n"
				}

				for _, r := range res.DNS.Extra {
					if r != nil && r.Header().Rrtype != dns.TypeOPT {
						s += r.String() + "\n"
					}
				}
			}
		}

		if opts.Display.Statistics {
			s += "\n;; Query time: " + res.RTT.String()
			// Add extra information to server string
			var extra string

			switch {
			case opts.TCP:
				extra = ":" + strconv.Itoa(opts.Port) + " (TCP)"
			case opts.TLS:
				extra = ":" + strconv.Itoa(opts.Port) + " (TLS)"
			case opts.HTTPS, opts.DNSCrypt:
				extra = ""
			case opts.QUIC:
				extra = ":" + strconv.Itoa(opts.Port) + " (QUIC)"
			default:
				extra = ":" + strconv.Itoa(opts.Port) + " (UDP)"
			}

			s += "\n;; SERVER: " + opts.Request.Server + extra
			s += "\n;; WHEN: " + time.Now().Format(time.RFC1123Z)
			s += "\n;; MSG SIZE  rcvd: " + strconv.Itoa(res.DNS.Len()) + "\n"
		}
	} else {
		// Print just the responses, nothing else
		for i, resp := range res.DNS.Answer {
			temp := strings.Split(resp.String(), "\t")
			s += temp[len(temp)-1]

			if opts.Identify {
				s += " from server " + opts.Request.Server + " in " + res.RTT.String()
			}

			// Don't print newline on last line
			if i != len(res.DNS.Answer)-1 {
				s += "\n"
			}
		}
	}

	return s
}
