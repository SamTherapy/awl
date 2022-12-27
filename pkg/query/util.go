package query

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"git.froth.zone/sam/awl/pkg/util"
	"github.com/miekg/dns"
	"golang.org/x/net/idna"
)

func (message *Message) displayQuestion(msg *dns.Msg, opts *util.Options, opt *dns.OPT) error {
	var (
		name string
		err  error
	)

	for _, question := range msg.Question {
		if opts.Display.UcodeTranslate {
			name, err = idna.ToUnicode(question.Name)
			if err != nil {
				return fmt.Errorf("punycode to unicode: %w", err)
			}
		} else {
			name = question.Name
		}

		message.Name = name
		message.Type = question.Qtype
		message.TypeName = dns.TypeToString[question.Qtype]
		message.Class = question.Qclass
		message.ClassName = dns.ClassToString[question.Qclass]
	}

	return nil
}

func (message *Message) displayAnswers(msg *dns.Msg, opts *util.Options, opt *dns.OPT) error {
	var (
		ttl  any
		name string
		err  error
	)

	for _, answer := range msg.Answer {
		temp := strings.Split(answer.String(), "\t")

		if opts.Display.TTL {
			if opts.Display.HumanTTL {
				ttl = (time.Duration(answer.Header().Ttl) * time.Second).String()
			} else {
				ttl = answer.Header().Ttl
			}
		}

		if opts.Display.UcodeTranslate {
			name, err = idna.ToUnicode(answer.Header().Name)
			if err != nil {
				return fmt.Errorf("punycode to unicode: %w", err)
			}
		} else {
			name = answer.Header().Name
		}

		message.AnswerRRs = append(message.AnswerRRs, Answer{
			Name:      name,
			ClassName: dns.ClassToString[answer.Header().Class],
			Class:     answer.Header().Class,
			TypeName:  dns.TypeToString[answer.Header().Rrtype],
			Type:      answer.Header().Rrtype,
			Rdlength:  answer.Header().Rdlength,
			TTL:       ttl,

			Value: temp[len(temp)-1],
		})
	}

	return nil
}

func (message *Message) displayAuthority(msg *dns.Msg, opts *util.Options, opt *dns.OPT) error {
	var (
		ttl  any
		name string
		err  error
	)

	for _, ns := range msg.Ns {
		temp := strings.Split(ns.String(), "\t")

		if opts.Display.TTL {
			if opts.Display.HumanTTL {
				ttl = (time.Duration(ns.Header().Ttl) * time.Second).String()
			} else {
				ttl = ns.Header().Ttl
			}
		}

		if opts.Display.UcodeTranslate {
			name, err = idna.ToUnicode(ns.Header().Name)
			if err != nil {
				return fmt.Errorf("punycode to unicode: %w", err)
			}
		} else {
			name = ns.Header().Name
		}

		message.AuthoritativeRRs = append(message.AuthoritativeRRs, Answer{
			Name:      name,
			TypeName:  dns.TypeToString[ns.Header().Rrtype],
			Type:      ns.Header().Rrtype,
			Class:     ns.Header().Class,
			ClassName: dns.ClassToString[ns.Header().Class],
			Rdlength:  ns.Header().Rdlength,
			TTL:       ttl,

			Value: temp[len(temp)-1],
		})
	}

	return nil
}

func (message *Message) displayAdditional(msg *dns.Msg, opts *util.Options, opt *dns.OPT) error {
	var (
		ttl  any
		name string
		err  error
	)

	for _, additional := range msg.Extra {
		if additional.Header().Rrtype == dns.StringToType["OPT"] {
			continue
		} else {
			temp := strings.Split(additional.String(), "\t")

			if opts.Display.TTL {
				if opts.Display.HumanTTL {
					ttl = (time.Duration(additional.Header().Ttl) * time.Second).String()
				} else {
					ttl = additional.Header().Ttl
				}
			}

			if opts.Display.UcodeTranslate {
				name, err = idna.ToUnicode(additional.Header().Name)
				if err != nil {
					return fmt.Errorf("punycode to unicode: %w", err)
				}
			} else {
				name = additional.Header().Name
			}
			message.AdditionalRRs = append(message.AdditionalRRs, Answer{
				Name:      name,
				TypeName:  dns.TypeToString[additional.Header().Rrtype],
				Type:      additional.Header().Rrtype,
				Class:     additional.Header().Class,
				ClassName: dns.ClassToString[additional.Header().Class],
				Rdlength:  additional.Header().Rdlength,
				TTL:       ttl,
				Value:     temp[len(temp)-1],
			})
		}
	}

	return nil
}

// ParseOpt parses opts.
func (message *Message) ParseOpt(rcode int, rr dns.OPT) (ret EDNS0, err error) {
	ret.Rcode = dns.RcodeToString[rcode]

	// Most of this is taken from https://github.com/miekg/dns/blob/master/edns.go#L76
	if rr.Do() {
		ret.Flags = append(ret.Flags, "DO")
	}

	for i := uint32(1); i <= 0x7FFF; i <<= 1 {
		if rr.Hdr.Ttl&i != 0 {
			ret.Flags = append(ret.Flags, fmt.Sprintf("BIT%d", i))
		}
	}

	ret.PayloadSize = rr.UDPSize()

	for _, opt := range rr.Option {
		switch opt := opt.(type) {
		case *dns.EDNS0_NSID:
			str := opt.String()

			hex, err := hex.DecodeString(str)
			if err != nil {
				return ret, fmt.Errorf("%w", err)
			}

			ret.NsidHex = string(hex)
			ret.Nsid = str

		case *dns.EDNS0_SUBNET:
			ret.Subnet = &EDNSSubnet{
				Source: opt.SourceNetmask,
				Family: opt.Family,
			}

			// 1: IPv4 2: IPv6
			if ret.Subnet.Family <= 2 {
				ret.Subnet.IP = opt.Address.String()
			} else {
				ret.Subnet.IP = hex.EncodeToString([]byte(opt.Address))
			}

			if opt.SourceScope != 0 {
				ret.Subnet.Scope = opt.SourceScope
			}

		case *dns.EDNS0_COOKIE:
			ret.Cookie = append(ret.Cookie, opt.String())

		case *dns.EDNS0_EXPIRE:
			ret.Expire = opt.Expire

		case *dns.EDNS0_TCP_KEEPALIVE:
			ret.KeepAlive = opt.Timeout

		case *dns.EDNS0_LLQ:
			ret.LLQ = &EdnsLLQ{
				Version: opt.Version,
				Opcode:  opt.Opcode,
				Error:   opt.Error,
				ID:      opt.Id,
				Lease:   opt.LeaseLife,
			}

		case *dns.EDNS0_DAU:
			ret.Dau = opt.AlgCode

		case *dns.EDNS0_DHU:
			ret.Dhu = opt.AlgCode

		case *dns.EDNS0_N3U:
			ret.N3u = opt.AlgCode

		case *dns.EDNS0_PADDING:
			ret.Padding = string(opt.Padding)

		case *dns.EDNS0_EDE:
			ret.EDE = &EDNSErr{
				Code:    opt.InfoCode,
				Purpose: dns.ExtendedErrorCodeToString[opt.InfoCode],
				Text:    opt.ExtraText,
			}
		}
	}

	return ret, nil
}
