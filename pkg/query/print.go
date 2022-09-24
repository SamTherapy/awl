// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"encoding/hex"
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

	return s, nil
}

func serverExtra(opts util.Options) string {
	// Add extra information to server string
	var extra string

	switch {
	case opts.TCP:
		extra = ":" + strconv.Itoa(opts.Request.Port) + " (TCP)"
	case opts.TLS:
		extra = ":" + strconv.Itoa(opts.Request.Port) + " (TLS)"
	case opts.HTTPS, opts.DNSCrypt:
		extra = ""
	case opts.QUIC:
		extra = ":" + strconv.Itoa(opts.Request.Port) + " (QUIC)"
	default:
		extra = ":" + strconv.Itoa(opts.Request.Port) + " (UDP)"
	}

	return extra
}

// stringParse edits the raw responses to user requests.
func stringParse(str string, isAns bool, opts util.Options) (string, error) {
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
func PrintSpecial(res util.Response, opts util.Options) (string, error) {
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
func MakePrintable(res util.Response, opts util.Options) (*Message, error) {
	var (
		err error
		msg = res.DNS
	)
	// The things I do for compatibility
	ret := Message{
		Header: Header{
			ID:                 msg.Id,
			Response:           msg.Response,
			Opcode:             dns.OpcodeToString[msg.Opcode],
			Authoritative:      msg.Authoritative,
			Truncated:          msg.Truncated,
			RecursionDesired:   msg.RecursionDesired,
			RecursionAvailable: msg.RecursionAvailable,
			Zero:               msg.Zero,
			AuthenticatedData:  msg.AuthenticatedData,
			CheckingDisabled:   msg.CheckingDisabled,
			Status:             dns.RcodeToString[msg.Rcode],
		},
	}

	opt := msg.IsEdns0()
	if opt != nil && opts.Display.Opt {
		ret.Opt, err = parseOpt(*opt)
		if err != nil {
			return nil, fmt.Errorf("edns print: %w", err)
		}
	}

	for _, question := range msg.Question {
		var name string
		if opts.Display.UcodeTranslate {
			name, err = idna.ToUnicode(question.Name)
			if err != nil {
				return nil, fmt.Errorf("punycode to unicode: %w", err)
			}
		} else {
			name = question.Name
		}

		ret.Question = append(ret.Question, Question{
			Name:  name,
			Type:  dns.TypeToString[question.Qtype],
			Class: dns.ClassToString[question.Qclass],
		})
	}

	for _, answer := range msg.Answer {
		temp := strings.Split(answer.String(), "\t")

		var (
			ttl  string
			name string
		)

		if opts.Display.TTL {
			if opts.Display.HumanTTL {
				ttl = (time.Duration(answer.Header().Ttl) * time.Second).String()
			} else {
				ttl = strconv.Itoa(int(answer.Header().Ttl))
			}
		}

		if opts.Display.UcodeTranslate {
			name, err = idna.ToUnicode(answer.Header().Name)
			if err != nil {
				return nil, fmt.Errorf("punycode to unicode: %w", err)
			}
		} else {
			name = answer.Header().Name
		}

		ret.Answer = append(ret.Answer, Answer{
			RRHeader: RRHeader{
				Name:     name,
				Type:     dns.TypeToString[answer.Header().Rrtype],
				Class:    dns.ClassToString[answer.Header().Class],
				Rdlength: answer.Header().Rdlength,
				TTL:      ttl,
			},
			Value: temp[len(temp)-1],
		})
	}

	for _, ns := range msg.Ns {
		temp := strings.Split(ns.String(), "\t")

		var (
			ttl  string
			name string
		)

		if opts.Display.TTL {
			if opts.Display.HumanTTL {
				ttl = (time.Duration(ns.Header().Ttl) * time.Second).String()
			} else {
				ttl = strconv.Itoa(int(ns.Header().Ttl))
			}
		}

		if opts.Display.UcodeTranslate {
			name, err = idna.ToUnicode(ns.Header().Name)
			if err != nil {
				return nil, fmt.Errorf("punycode to unicode: %w", err)
			}
		} else {
			name = ns.Header().Name
		}

		ret.Authority = append(ret.Authority, Answer{
			RRHeader: RRHeader{
				Name:     name,
				Type:     dns.TypeToString[ns.Header().Rrtype],
				Class:    dns.ClassToString[ns.Header().Class],
				Rdlength: ns.Header().Rdlength,
				TTL:      ttl,
			},
			Value: temp[len(temp)-1],
		})
	}

	for _, additional := range msg.Extra {
		if additional.Header().Rrtype == dns.StringToType["OPT"] {
			continue
		} else {
			temp := strings.Split(additional.String(), "\t")

			var (
				ttl  string
				name string
			)

			if opts.Display.TTL {
				if opts.Display.HumanTTL {
					ttl = (time.Duration(additional.Header().Ttl) * time.Second).String()
				} else {
					ttl = strconv.Itoa(int(additional.Header().Ttl))
				}
			}

			if opts.Display.UcodeTranslate {
				name, err = idna.ToUnicode(additional.Header().Name)
				if err != nil {
					return nil, fmt.Errorf("punycode to unicode: %w", err)
				}
			} else {
				name = additional.Header().Name
			}

			ret.Additional = append(ret.Additional, Answer{
				RRHeader: RRHeader{
					Name:     name,
					Type:     dns.TypeToString[additional.Header().Rrtype],
					Class:    dns.ClassToString[additional.Header().Class],
					Rdlength: additional.Header().Rdlength,
					TTL:      ttl,
				},
				Value: temp[len(temp)-1],
			})
		}
	}

	if opts.Display.Statistics {
		ret.Statistics = Statistics{
			RTT:     res.RTT.String(),
			Server:  opts.Request.Server + serverExtra(opts),
			When:    time.Now().Format(time.RFC1123Z),
			MsgSize: res.DNS.Len(),
		}
	} else {
		ret.Statistics = Statistics{}
	}

	return &ret, nil
}

func parseOpt(rr dns.OPT) ([]Opts, error) {
	ret := []Opts{}
	// Most of this is taken from https://github.com/miekg/dns/blob/master/edns.go#L76

	ret = append(ret, Opts{
		Name:  "Version",
		Value: strconv.Itoa(int(rr.Version())),
	})

	if rr.Do() {
		ret = append(ret, Opts{
			Name:  "Flags",
			Value: "do",
		})
	} else {
		ret = append(ret, Opts{
			Name:  "Flags",
			Value: "",
		})
	}

	if rr.Hdr.Ttl&0x7FFF != 0 {
		ret = append(ret, Opts{
			Name:  "MBZ",
			Value: fmt.Sprintf("0x%04x", rr.Hdr.Ttl&0x7FFF),
		})
	}

	ret = append(ret, Opts{
		Name:  "UDP Buffer Size",
		Value: strconv.Itoa(int(rr.UDPSize())),
	})

	for _, opt := range rr.Option {
		switch opt.(type) {
		case *dns.EDNS0_NSID:
			str := opt.String()

			hex, err := hex.DecodeString(str)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}

			ret = append(ret, Opts{
				Name:  "NSID",
				Value: fmt.Sprintf("%s (%s)", str, string(hex)),
			})
		case *dns.EDNS0_SUBNET:
			ret = append(ret, Opts{
				Name:  "Subnet",
				Value: opt.String(),
			})
		case *dns.EDNS0_COOKIE:
			ret = append(ret, Opts{
				Name:  "Cookie",
				Value: opt.String(),
			})
		case *dns.EDNS0_EXPIRE:
			ret = append(ret, Opts{
				Name:  "Expire",
				Value: opt.String(),
			})
		case *dns.EDNS0_TCP_KEEPALIVE:
			ret = append(ret, Opts{
				Name:  "TCP Keepalive",
				Value: opt.String(),
			})
		case *dns.EDNS0_UL:
			ret = append(ret, Opts{
				Name:  "Update Lease",
				Value: opt.String(),
			})
		case *dns.EDNS0_LLQ:
			ret = append(ret, Opts{
				Name:  "Long Lived Queries",
				Value: opt.String(),
			})
		case *dns.EDNS0_DAU:
			ret = append(ret, Opts{
				Name:  "DNSSEC Algorithm Understood",
				Value: opt.String(),
			})
		case *dns.EDNS0_DHU:
			ret = append(ret, Opts{
				Name:  "DS Hash Understood",
				Value: opt.String(),
			})
		case *dns.EDNS0_N3U:
			ret = append(ret, Opts{
				Name:  "NSEC3 Hash Understood",
				Value: opt.String(),
			})
		case *dns.EDNS0_LOCAL:
			ret = append(ret, Opts{
				Name:  "Local OPT",
				Value: opt.String(),
			})
		case *dns.EDNS0_PADDING:
			ret = append(ret, Opts{
				Name:  "Padding",
				Value: opt.String(),
			})
		case *dns.EDNS0_EDE:
			ret = append(ret, Opts{
				Name:  "EDE",
				Value: opt.String(),
			})
		case *dns.EDNS0_ESU:
			ret = append(ret, Opts{
				Name:  "ESU",
				Value: opt.String(),
			})
		}
	}

	return ret, nil
}

var errInvalidFormat = errors.New("this should never happen")
