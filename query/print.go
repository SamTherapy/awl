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

	"git.froth.zone/sam/awl/cli"
	"github.com/miekg/dns"
	"golang.org/x/net/idna"
	"gopkg.in/yaml.v3"
)

func PrintSpecial(msg *dns.Msg, opts cli.Options) (string, error) {
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
// and XML. Little is changed beyond naming
func MakePrintable(msg *dns.Msg, opts cli.Options) (*Message, error) {
	var err error
	ret := Message{
		Header: msg.MsgHdr,
	}

	for _, question := range msg.Question {
		var name string
		if opts.Display.UcodeTranslate {
			name, err = idna.ToUnicode(question.Name)
			if err != nil {
				return nil, fmt.Errorf("punycode: error translating to unicode: %w", err)
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
				return nil, fmt.Errorf("punycode: error translating to unicode: %w", err)
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
				return nil, fmt.Errorf("punycode: error translating to unicode: %w", err)
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
					return nil, fmt.Errorf("punycode: error translating to unicode: %w", err)
				}
			} else {
				name = additional.Header().Name
			}
			ret.Extra = append(ret.Extra, Answer{
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

	return &ret, nil
}

var errInvalidFormat = errors.New("this should never happen")
