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
	opts.Logger.Info("Setting", arg)

	switch arg {
	// Set DNS query flags
	case "aa", "aaflag", "aaonly", "noaa", "noaaflag", "noaaonly":
		opts.AA = isNo
	case "ad", "adflag", "noad", "noadflag":
		opts.AD = isNo
	case "cd", "cdflag", "nocd", "nocdflag":
		opts.CD = isNo
	case "qr", "qrflag", "noqr", "noqrflag":
		opts.QR = isNo
	case "ra", "raflag", "nora", "noraflag":
		opts.RA = isNo
	case "rd", "rdflag", "recurse", "nord", "nordflag", "norecurse":
		opts.RD = isNo
	case "tc", "tcflag", "notc", "notcflag":
		opts.TC = isNo
	case "z", "zflag", "noz", "nozflag":
		opts.Z = isNo
	// End DNS query flags

	// DNS-over-X
	case "dnssec", "nodnssec":
		opts.DNSSEC = isNo
	case "tcp", "vc", "notcp", "novc":
		opts.TCP = isNo
	case "ignore", "noignore":
		opts.Truncate = isNo
	case "tls", "notls":
		opts.TLS = isNo
	case "dnscrypt", "nodnscrypt":
		opts.DNSCrypt = isNo
	case "https", "nohttps":
		opts.HTTPS = isNo
	case "quic", "noquic":
		opts.QUIC = isNo
	// End DNS-over-X

	// Formatting
	case "short", "noshort":
		opts.Short = isNo
	case "json", "nojson":
		opts.JSON = isNo
	case "xml", "noxml":
		opts.XML = isNo
	case "yaml", "noyaml":
		opts.YAML = isNo
	// End formatting

	// Output
	// TODO: get this to work
	// case "comments", "nocomments":
	// 	opts.Display.Comments = isNo
	case "question", "noquestion":
		opts.Display.Question = isNo
	case "answer", "noanswer":
		opts.Display.Answer = isNo
	case "authority", "noauthority":
		opts.Display.Authority = isNo
	case "additional", "noadditional":
		opts.Display.Additional = isNo
	case "stats", "nostats":
		opts.Display.Statistics = isNo

	case "all", "noall":
		opts.Display.Question = isNo
		opts.Display.Answer = isNo
		opts.Display.Authority = isNo
		opts.Display.Additional = isNo
		opts.Display.Statistics = isNo

	default:
		// Recursive switch statements WOO
		switch {
		case strings.HasPrefix(arg, "timeout"):
			timeout, err := strconv.Atoi(strings.Split(arg, "=")[1])

			if err != nil {
				return fmt.Errorf("dig: Invalid timeout value")
			}

			opts.Request.Timeout = time.Duration(timeout)

		case strings.HasPrefix(arg, "retry"), strings.HasPrefix(arg, "tries"):
			tries, err := strconv.Atoi(strings.Split(arg, "=")[1])
			if err != nil {
				return fmt.Errorf("dig: Invalid retry value")
			}

			if strings.HasPrefix(arg, "tries") {
				tries++
			}

			opts.Request.Retries = tries

		default:
			return fmt.Errorf("dig: unknown flag %s given", arg)
		}
	}
	return nil
}
