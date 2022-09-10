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

	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"golang.org/x/net/idna"
	"gopkg.in/yaml.v3"
)

// PrintSpecial is for printing as JSON, XML or YAML.
// As of now JSON and XML use the stdlib version.
func PrintSpecial(msg *dns.Msg, opts util.Options) (string, error) {
	formatted, err := MakePrintable(msg, opts)
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
func MakePrintable(msg *dns.Msg, opts util.Options) (*Message, error) {
	var err error

	ret := Message{
		Header: msg.MsgHdr,
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

		if opts.ShowTTL {
			if opts.HumanTTL {
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

		if opts.ShowTTL {
			if opts.HumanTTL {
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

		ret.Ns = append(ret.Ns, Answer{
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

			if opts.ShowTTL {
				if opts.HumanTTL {
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

	opt := msg.IsEdns0()
	if opt != nil && opts.Display.Opt {
		ret.Opt, err = parseOpt(*opt)
		if err != nil {
			return nil, fmt.Errorf("edns print: %w", err)
		}
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
