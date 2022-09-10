// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"golang.org/x/net/idna"
)

// Message is for overall DNS responses.
//
//nolint:govet // Better output is worth 32 bytes.
type Message struct {
	Header     dns.MsgHdr `json:"header,omitempty" xml:"header,omitempty" yaml:",omitempty"`
	Opt        []Opts     `json:"opt,omitempty" xml:"opt,omitempty" yaml:"opt,omitempty"`
	Question   []Question `json:"question,omitempty" xml:"question,omitempty" yaml:",omitempty"`
	Answer     []Answer   `json:"answer,omitempty" xml:"answer,omitempty" yaml:",omitempty"`
	Ns         []Answer   `json:"ns,omitempty" xml:"ns,omitempty" yaml:",omitempty"`
	Additional []Answer   `json:"additional,omitempty" xml:"additional,omitempty" yaml:",omitempty"`
}

// Question is a DNS Query.
type Question struct {
	Name  string `json:"name,omitempty" xml:"name,omitempty" yaml:",omitempty"`
	Class string `json:"class,omitempty" xml:"class,omitempty" yaml:",omitempty"`
	Type  string `json:"type,omitempty" xml:"type,omitempty" yaml:",omitempty"`
}

// RRHeader is for DNS Resource Headers.
type RRHeader struct {
	Name     string `json:"name,omitempty" xml:"name,omitempty" yaml:",omitempty"`
	TTL      string `json:"ttl,omitempty" xml:"ttl,omitempty" yaml:",omitempty"`
	Class    string `json:"class,omitempty" xml:"class,omitempty" yaml:",omitempty"`
	Type     string `json:"type,omitempty" xml:"type,omitempty" yaml:",omitempty"`
	Rdlength uint16 `json:"-" xml:"-" yaml:"-"`
}

// Opts is for the OPT pseudosection, nearly exclusively for EDNS.
type Opts struct {
	Name  string `json:"name,omitempty" xml:"name,omitempty" yaml:",omitempty"`
	Value string `json:"value" xml:"value" yaml:""`
}

// Answer is for a DNS Response.
type Answer struct {
	Value    string `json:"response,omitempty" xml:"response,omitempty" yaml:"response,omitempty"`
	RRHeader `json:"header,omitempty" xml:"header,omitempty" yaml:"header,omitempty"`
}

// ToString turns the response into something that looks a lot like dig
//
// Much of this is taken from https://github.com/miekg/dns/blob/master/msg.go#L900
func ToString(res util.Response, opts util.Options) (string, error) {
	if res.DNS == nil {
		return "<nil> MsgHdr", errNoMessage
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
					str, err := stringParse(r.String(), false, opts)
					if err != nil {
						return "", fmt.Errorf("%w", err)
					}

					s += str + "\n"
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
						str, err := stringParse(r.String(), true, opts)
						if err != nil {
							return "", fmt.Errorf("%w", err)
						}

						s += str + "\n"
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
						str, err := stringParse(r.String(), true, opts)
						if err != nil {
							return "", fmt.Errorf("%w", err)
						}

						s += str + "\n"
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
						str, err := stringParse(r.String(), true, opts)
						if err != nil {
							return "", fmt.Errorf("%w", err)
						}

						s += str + "\n"
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

	return s, nil
}

func stringParse(str string, isAns bool, opts util.Options) (string, error) {
	split := strings.Split(str, "\t")

	// Make edits if so requested

	// TODO: make less ew?
	// This exists because the question section should be left alone EXCEPT for punycode.

	if isAns {
		if !opts.ShowTTL {
			// Remove from existence
			split = append(split[:1], split[2:]...)
		}

		if !opts.ShowClass {
			// Position depends on if the TTL is there or not.
			if opts.ShowTTL {
				split = append(split[:2], split[3:]...)
			} else {
				split = append(split[:1], split[2:]...)
			}
		}

		if opts.ShowTTL && opts.HumanTTL {
			ttl, _ := strconv.Atoi(split[1])
			split[1] = (time.Duration(ttl) * time.Second).String()
		}
	}

	if opts.Display.UcodeTranslate {
		var (
			err  error
			semi string
		)

		if strings.HasPrefix(split[0], ";") {
			split[0] = strings.TrimPrefix(split[0], ";")
			semi = ";"
		}

		split[0], err = idna.ToUnicode(split[0])
		if err != nil {
			return "", fmt.Errorf("punycode: %w", err)
		}

		split[0] = semi + split[0]
	}

	return strings.Join(split, "\t"), nil
}

var errNoMessage = errors.New("no message")
