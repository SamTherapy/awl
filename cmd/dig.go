// SPDX-License-Identifier: BSD-3-Clause

package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"git.froth.zone/sam/awl/pkg/util"
)

// ParseDig parses commands from the popular DNS tool dig.
// All dig commands are taken from https://man.openbsd.org/dig.1 as the source of their functionality.
//
// [no]flags are supported just as flag are and are disabled as such.
func ParseDig(arg string, opts *util.Options) error {
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
		opts.Display.ShowQuery = isNo
	case "ttlunits":
		opts.Display.HumanTTL = isNo
	case "ttl", "ttlid":
		opts.Display.TTL = isNo
	case "class":
		opts.Display.ShowClass = isNo

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
	case "badcookie":
		opts.BadCookie = !isNo
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
	case "idnin", "idnout":
		opts.Display.UcodeTranslate = isNo

	default:
		if err := parseDigEq(isNo, arg, opts); err != nil {
			return err
		}
	}

	return nil
}

// For flags that contain "=".
func parseDigEq(startNo bool, arg string, opts *util.Options) error {
	// Recursive switch statements WOO
	arg, val, isSplit := strings.Cut(arg, "=")
	switch arg {
	case "time", "timeout":
		if isSplit && val != "" {
			timeout, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("digflags: timeout : %w", err)
			}

			opts.Request.Timeout = time.Duration(timeout)
		} else {
			return fmt.Errorf("digflags: timeout: %w", errNoArg)
		}

	case "retry", "tries":
		if isSplit && val != "" {
			tries, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("digflags: retry: %w", err)
			}

			opts.Request.Retries = tries

			// TODO: Is there a better way to do this?
			if arg == "tries" {
				opts.Request.Retries--
			}
		} else {
			return fmt.Errorf("digflags: retry: %w", errNoArg)
		}

	case "bufsize":
		if isSplit && val != "" {
			size, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("digflags: EDNS UDP: %w", err)
			}

			opts.EDNS.BufSize = uint16(size)
		} else {
			return fmt.Errorf("digflags: EDNS UDP: %w", errNoArg)
		}

	case "ednsflags":
		if isSplit && val != "" {
			ver, err := strconv.ParseInt(val, 0, 16)
			if err != nil {
				return fmt.Errorf("digflags: EDNS flag: %w", err)
			}

			// Ignore setting DO bit
			opts.EDNS.ZFlag = uint16(ver & 0x7FFF)
		} else {
			opts.EDNS.ZFlag = 0
		}

	case "edns":
		opts.EDNS.EnableEDNS = startNo

		if isSplit && val != "" {
			ver, err := strconv.Atoi(val)
			if err != nil {
				return fmt.Errorf("digflags: EDNS version: %w", err)
			}

			opts.EDNS.Version = uint8(ver)
		} else {
			opts.EDNS.Version = 0
		}

	case "subnet":
		if isSplit && val != "" {
			err := util.ParseSubnet(val, opts)
			if err != nil {
				return fmt.Errorf("digflags: EDNS Subnet: %w", err)
			}
		} else {
			return fmt.Errorf("digflags: EDNS Subnet: %w", errNoArg)
		}

	default:
		return &errInvalidArg{arg}
	}

	return nil
}
