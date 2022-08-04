// SPDX-License-Identifier: BSD-3-Clause

package cli

import (
	"fmt"
	"math/rand"
	"strings"

	"git.froth.zone/sam/awl/conf"
	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"golang.org/x/net/idna"
)

// ParseMiscArgs parses the wildcard arguments, drill style.
// Only one command is supported at a time, so any extra information overrides previous.
func ParseMiscArgs(args []string, opts *Options) error {
	var err error

	for _, arg := range args {
		r, ok := dns.StringToType[strings.ToUpper(arg)]
		switch {
		// If it starts with @, it's a DNS server
		case strings.HasPrefix(arg, "@"):
			arg = arg[1:]
			// Automatically set flags based on URI header
			opts.Logger.Info(arg, "detected as a server")
			switch {
			case strings.HasPrefix(arg, "tls://"):
				opts.TLS = true
				opts.Request.Server = arg[6:]
				opts.Logger.Info("DNS-over-TLS implicitly set")
			case strings.HasPrefix(arg, "https://"):
				opts.HTTPS = true
				opts.Request.Server = arg
				opts.Logger.Info("DNS-over-HTTPS implicitly set")
			case strings.HasPrefix(arg, "quic://"):
				opts.QUIC = true
				opts.Request.Server = arg[7:]
				opts.Logger.Info("DNS-over-QUIC implicitly set.")
			case strings.HasPrefix(arg, "sdns://"):
				opts.DNSCrypt = true
				opts.Request.Server = arg
				opts.Logger.Info("DNSCrypt implicitly set")
			default:
				opts.Request.Server = arg
			}

		// Dig-style +queries
		case strings.HasPrefix(arg, "+"):
			opts.Logger.Info(arg, "detected as a dig query")
			err = ParseDig(strings.ToLower(arg[1:]), opts)
			if err != nil {
				return err
			}

		// Domain names
		case strings.Contains(arg, "."):
			opts.Logger.Info(arg, "detected as a domain name")
			opts.Request.Name, err = idna.ToASCII(arg)
			if err != nil {
				return fmt.Errorf("punycode translate error: %w", err)
			}

		// DNS query type
		case ok:
			opts.Logger.Info(arg, "detected as a type")
			opts.Request.Type = r

		// Domain?
		default:
			opts.Logger.Info(arg, "is unknown. Assuming domain")
			opts.Request.Name, err = idna.ToASCII(arg)
			if err != nil {
				return fmt.Errorf("punycode translate error: %w", err)
			}
		}
	}

	// If nothing was set, set a default
	if opts.Request.Name == "" {
		opts.Logger.Info("Domain not specified, making a default")
		opts.Request.Name = "."
		if opts.Request.Type == 0 {
			opts.Request.Type = dns.StringToType["NS"]
		}
	} else {
		opts.Logger.Info("Query not specified, making an \"A\" query")
		if opts.Request.Type == 0 {
			opts.Request.Type = dns.StringToType["A"]
		}
	}
	//
	if opts.Request.Server == "" {
		opts.Logger.Info("Server not specified, selecting a default")
		// Set "defaults" for each if there is no input
		switch {
		case opts.DNSCrypt:
			// This is adguard
			opts.Request.Server = "sdns://AQMAAAAAAAAAETk0LjE0MC4xNC4xNDo1NDQzINErR_JS3PLCu_iZEIbq95zkSV2LFsigxDIuUso_OQhzIjIuZG5zY3J5cHQuZGVmYXVsdC5uczEuYWRndWFyZC5jb20"
		case opts.TLS:
			opts.Request.Server = "dns.google"
		case opts.HTTPS:
			opts.Request.Server = "https://dns.cloudflare.com/dns-query"
		case opts.QUIC:
			opts.Request.Server = "dns.adguard.com"
		default:
			resolv, err := conf.GetDNSConfig()
			if err != nil {
				// :^)
				opts.Logger.Warn("Could not query system for server. Using default")
				opts.Request.Server = "95.216.99.249"
			} else {
				// Make sure that if IPv4 or IPv6 is asked for it actually uses it
			harmful:
				for _, srv := range resolv.Servers {
					switch {
					case opts.IPv4:
						if strings.Contains(srv, ".") {
							opts.Request.Server = srv
							break harmful
						}
					case opts.IPv6:
						if strings.Contains(srv, ":") {
							opts.Request.Server = srv
							break harmful
						}
					default:
						//#nosec -- This isn't used for anything secure
						opts.Request.Server = resolv.Servers[rand.Intn(len(resolv.Servers))]
						break harmful
					}
				}
			}
		}
	}
	opts.Logger.Info("DNS server set to", opts.Request.Server)

	// Make reverse adresses proper addresses
	if opts.Reverse {
		opts.Logger.Info("Making reverse DNS query proper *.arpa domain")
		if dns.TypeToString[opts.Request.Type] == "A" {
			opts.Request.Type = dns.StringToType["PTR"]
		}
		opts.Request.Name, err = util.ReverseDNS(opts.Request.Name, opts.Request.Type)
		if err != nil {
			return fmt.Errorf("reverse DNS error: %w", err)
		}
	}

	// if the domain is not canonical, make it canonical
	if !strings.HasSuffix(opts.Request.Name, ".") {
		opts.Request.Name = fmt.Sprintf("%s.", opts.Request.Name)
		opts.Logger.Debug("Domain made canonical")
	}
	return nil
}
