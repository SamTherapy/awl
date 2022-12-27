// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.froth.zone/sam/awl/pkg/util"
	"github.com/miekg/dns"
	"golang.org/x/net/idna"
	"gopkg.in/yaml.v3"
)

// ToString turns the response into something that looks a lot like dig
//
// Much of this is taken from https://github.com/miekg/dns/blob/master/msg.go#L900
func ToString(res util.Response, opts *util.Options) (s string, err error) {
	if res.DNS == nil {
		return "<nil> MsgHdr", errNoMessage
	}

	var opt *dns.OPT

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
			s += "\n;; SERVER: " + opts.Request.Server + serverExtra(opts)
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

	return
}

func serverExtra(opts *util.Options) string {
	switch {
	case opts.TCP:
		return " (TCP)"
	case opts.TLS:
		return " (TLS)"
	case opts.HTTPS, opts.DNSCrypt:
		return ""
	case opts.QUIC:
		return " (QUIC)"
	default:
		return " (UDP)"
	}
}

// stringParse edits the raw responses to user requests.
func stringParse(str string, isAns bool, opts *util.Options) (string, error) {
	split := strings.Split(str, "\t")

	// Make edits if so requested

	// TODO: make less ew?
	// This exists because the question section should be left alone EXCEPT for punycode.

	if isAns {
		if !opts.Display.TTL {
			// Remove from existence
			split = append(split[:1], split[2:]...)
		}

		if !opts.Display.ShowClass {
			// Position depends on if the TTL is there or not.
			if opts.Display.TTL {
				split = append(split[:2], split[3:]...)
			} else {
				split = append(split[:1], split[2:]...)
			}
		}

		if opts.Display.TTL && opts.Display.HumanTTL {
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

// PrintSpecial is for printing as JSON, XML or YAML.
// As of now JSON and XML use the stdlib version.
func PrintSpecial(res util.Response, opts *util.Options) (string, error) {
	formatted, err := MakePrintable(res, opts)
	if err != nil {
		return "", err
	}

	switch {
	case opts.JSON:
		opts.Logger.Info("Printing as JSON")

		json, err := json.MarshalIndent(formatted, " ", "  ")

		return string(json), err
	case opts.XML:
		opts.Logger.Info("Printing as XML")

		xml, err := xml.MarshalIndent(formatted, " ", "  ")

		return string(xml), err
	case opts.YAML:
		opts.Logger.Info("Printing as YAML")

		yaml, err := yaml.Marshal(formatted)

		return string(yaml), err
	default:
		return "", errInvalidFormat
	}
}

// MakePrintable takes a DNS message and makes it nicer to be printed as JSON,YAML,
// and XML. Little is changed beyond naming.
func MakePrintable(res util.Response, opts *util.Options) (*Message, error) {
	var (
		err error
		msg = res.DNS
	)
	// The things I do for compatibility
	ret := &Message{
		DateString:  time.Now().Format(time.RFC3339),
		DateSeconds: time.Now().Unix(),
		MsgSize:     res.DNS.Len(),
		ID:          msg.Id,
		Opcode:      msg.Opcode,
		Response:    msg.Response,

		Authoritative:      msg.Authoritative,
		Truncated:          msg.Truncated,
		RecursionDesired:   msg.RecursionDesired,
		RecursionAvailable: msg.RecursionAvailable,
		AuthenticatedData:  msg.AuthenticatedData,
		CheckingDisabled:   msg.CheckingDisabled,
		Zero:               msg.Zero,

		QdCount: len(msg.Question),
		AnCount: len(msg.Answer),
		NsCount: len(msg.Ns),
		ArCount: len(msg.Extra),
	}

	opt := msg.IsEdns0()
	if opt != nil && opts.Display.Opt {
		ret.EDNS0, err = ret.ParseOpt(msg.Rcode, *opt)
		if err != nil {
			return nil, fmt.Errorf("edns print: %w", err)
		}
	}

	if opts.Display.Question {
		err = ret.displayQuestion(msg, opts, opt)
		if err != nil {
			return nil, fmt.Errorf("unable to display questions: %w", err)
		}
	}

	if opts.Display.Answer {
		err = ret.displayAnswers(msg, opts, opt)
		if err != nil {
			return nil, fmt.Errorf("unable to display answers: %w", err)
		}
	}

	if opts.Display.Authority {
		err = ret.displayAuthority(msg, opts, opt)
		if err != nil {
			return nil, fmt.Errorf("unable to display authority: %w", err)
		}
	}

	if opts.Display.Additional {
		err = ret.displayAdditional(msg, opts, opt)
		if err != nil {
			return nil, fmt.Errorf("unable to display additional: %w", err)
		}
	}

	return ret, nil
}

var errInvalidFormat = errors.New("this should never happen")
