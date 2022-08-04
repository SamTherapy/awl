// SPDX-License-Identifier: BSD-3-Clause

package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDig parses commands from the popular DNS tool dig.
// All dig commands are taken from https://man.openbsd.org/dig.1 as the source of their functionality.
//
// [no]flags are supported just as flag are and are disabled as such.
func ParseDig(arg string, opts *Options) error {
	// returns true if the flag starts with a no
	isNo := !strings.HasPrefix(arg, "no")
	if !isNo {
		arg = strings.TrimPrefix(arg, "no")
	}
	opts.Logger.Info("Setting", arg)

	switch arg {
	// Set DNS query flags
	case "aa", "aaflag", "aaonly":
		opts.AA = isNo
	case "ad", "adflag":
		opts.AD = isNo
	case "cd", "cdflag":
		opts.CD = isNo
	case "qrflag":
		opts.QR = isNo
	case "ra", "raflag":
		opts.RA = isNo
	case "rd", "rdflag", "recurse":
		opts.RD = isNo
	case "tc", "tcflag":
		opts.TC = isNo
	case "z", "zflag":
		opts.Z = isNo
	// End DNS query flags

	case "qr":
		opts.ShowQuery = isNo
	case "ttlunits":
		opts.HumanTTL = isNo
	case "ttlid":
		opts.ShowTTL = isNo

	// EDNS queries
	case "dnssec":
		opts.EDNS.DNSSEC = isNo
	case "expire":
		opts.EDNS.Expire = isNo
	case "cookie":
		opts.EDNS.Cookie = isNo
	case "keepopen", "keepalive":
		opts.EDNS.KeepOpen = isNo
	case "nsid":
		opts.EDNS.Nsid = isNo
	case "padding":
		opts.EDNS.Padding = isNo
	// End EDNS queries

	// DNS-over-X
	case "tcp", "vc":
		opts.TCP = isNo
	case "ignore":
		opts.Truncate = isNo
	case "tls":
		opts.TLS = isNo
	case "dnscrypt":
		opts.DNSCrypt = isNo
	case "https":
		opts.HTTPS = isNo
	case "quic":
		opts.QUIC = isNo
	// End DNS-over-X

	// Formatting
	case "short":
		opts.Short = isNo
	case "identify":
		opts.Identify = isNo
	case "json":
		opts.JSON = isNo
	case "xml":
		opts.XML = isNo
	case "yaml":
		opts.YAML = isNo
	// End formatting

	// Output
	case "comments":
		opts.Display.Comments = isNo
	case "question":
		opts.Display.Question = isNo
	case "opt":
		opts.Display.Opt = isNo
	case "answer":
		opts.Display.Answer = isNo
	case "authority":
		opts.Display.Authority = isNo
	case "additional":
		opts.Display.Additional = isNo
	case "stats":
		opts.Display.Statistics = isNo
	case "all":
		opts.Display.Comments = isNo
		opts.Display.Question = isNo
		opts.Display.Opt = isNo
		opts.Display.Answer = isNo
		opts.Display.Authority = isNo
		opts.Display.Additional = isNo
		opts.Display.Statistics = isNo
	case "idnout":
		opts.Display.UcodeTranslate = isNo

	default:
		// Recursive switch statements WOO
		arg := strings.Split(arg, "=")
		switch arg[0] {
		case "time", "timeout":
			if len(arg) > 1 && arg[1] != "" {
				timeout, err := strconv.Atoi(arg[1])
				if err != nil {
					return fmt.Errorf("digflags: Invalid timeout value: %w", err)
				}

				opts.Request.Timeout = time.Duration(timeout)
			} else {
				return fmt.Errorf("digflags: Invalid timeout value: %w", errNoArg)
			}

		case "retry", "tries":
			if len(arg) > 1 && arg[1] != "" {
				tries, err := strconv.Atoi(arg[1])
				if err != nil {
					return fmt.Errorf("digflags: Invalid retry value: %w", err)
				}
				opts.Request.Retries = tries

				// TODO: Is there a better way to do this?
				if arg[0] == "tries" {
					opts.Request.Retries++
				}
			} else {
				return fmt.Errorf("digflags: Invalid retry value: %w", errNoArg)
			}

		case "bufsize":
			if len(arg) > 1 && arg[1] != "" {
				size, err := strconv.Atoi(arg[1])
				if err != nil {
					return fmt.Errorf("digflags: Invalid UDP buffer size value: %w", err)
				}
				opts.EDNS.BufSize = uint16(size)
			} else {
				return fmt.Errorf("digflags: Invalid UDP buffer size value: %w", errNoArg)
			}

		case "ednsflags":
			if len(arg) > 1 && arg[1] != "" {
				ver, err := strconv.ParseInt(arg[1], 0, 16)
				if err != nil {
					return fmt.Errorf("digflags: Invalid EDNS flag: %w", err)
				}
				// Ignore setting DO bit
				opts.EDNS.ZFlag = uint16(ver & 0x7FFF)
			} else {
				opts.EDNS.ZFlag = 0
			}

		case "edns":
			opts.EDNS.EnableEDNS = isNo
			if len(arg) > 1 && arg[1] != "" {
				ver, err := strconv.Atoi(arg[1])
				if err != nil {
					return fmt.Errorf("digflags: Invalid EDNS version: %w", err)
				}
				opts.EDNS.Version = uint8(ver)
			} else {
				opts.EDNS.Version = 0
			}

		case "subnet":
			if len(arg) > 1 && arg[1] != "" {
				err := parseSubnet(arg[1], opts)
				if err != nil {
					return fmt.Errorf("digflags: Invalid EDNS Subnet: %w", err)
				}
			} else {
				return fmt.Errorf("digflags: Invalid EDNS Subnet: %w", errNoArg)
			}

		default:
			return &errInvalidArg{arg[0]}
		}
	}
	return nil
}
