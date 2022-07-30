// SPDX-License-Identifier: BSD-3-Clause

package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Parse dig-like commands and set the options as such.
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
		opts.Display.Answer = isNo
		opts.Display.Authority = isNo
		opts.Display.Additional = isNo
		opts.Display.Statistics = isNo
	case "idnout":
		opts.Display.UcodeTranslate = isNo

	default:
		// Recursive switch statements WOO
		switch {
		case strings.HasPrefix(arg, "time"), strings.HasPrefix(arg, "timeout"):
			timeout, err := strconv.Atoi(strings.Split(arg, "=")[1])

			if err != nil {
				return fmt.Errorf("digflags: Invalid timeout value")
			}

			opts.Request.Timeout = time.Duration(timeout)

		case strings.HasPrefix(arg, "retry"), strings.HasPrefix(arg, "tries"):
			tries, err := strconv.Atoi(strings.Split(arg, "=")[1])
			if err != nil {
				return fmt.Errorf("digflags: Invalid retry value")
			}

			if strings.HasPrefix(arg, "tries") {
				tries++
			}

			opts.Request.Retries = tries

		case strings.HasPrefix(arg, "bufsize"):
			size, err := strconv.Atoi(strings.Split(arg, "=")[1])
			if err != nil {
				return fmt.Errorf("digflags: Invalid UDP buffer size value")
			}
			opts.EDNS.BufSize = uint16(size)

		case strings.HasPrefix(arg, "ednsflags"):
			split := strings.Split(arg, "=")
			if len(split) > 1 {
				ver, err := strconv.ParseInt(split[1], 0, 16)
				if err != nil {
					return fmt.Errorf("digflags: Invalid EDNS flag")
				}
				// Ignore setting DO bit
				opts.EDNS.ZFlag = uint16(ver & 0x7FFF)
			} else {
				opts.EDNS.ZFlag = 0
			}

		case strings.HasPrefix(arg, "edns"):
			opts.EDNS.EnableEDNS = isNo
			split := strings.Split(arg, "=")
			if len(split) > 1 {
				ver, err := strconv.Atoi(split[1])
				if err != nil {
					return fmt.Errorf("digflags: Invalid EDNS version")
				}
				opts.EDNS.Version = uint8(ver)
			} else {
				opts.EDNS.Version = 0
			}

		case strings.HasPrefix(arg, "subnet"):

		default:
			return fmt.Errorf("digflags: unknown flag %s given", arg)
		}
	}
	return nil
}
