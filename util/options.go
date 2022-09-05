// SPDX-License-Identifier: BSD-3-Clause

package util

import (
	"fmt"
	"net"

	"git.froth.zone/sam/awl/logawl"
	"github.com/miekg/dns"
)

// Options is the grand CLI options structure.
type Options struct {
	Logger  *logawl.Logger
	TLSHost string
	EDNS
	Request     Request
	Port        int
	Verbosity   int
	Display     Displays
	TC          bool
	ShowTTL     bool
	ShowClass   bool
	ShowQuery   bool
	AA          bool
	AD          bool
	CD          bool
	QR          bool
	RD          bool
	RA          bool
	IPv4        bool
	Z           bool
	Reverse     bool
	HumanTTL    bool
	Truncate    bool
	Short       bool
	Identify    bool
	JSON        bool
	XML         bool
	YAML        bool
	QUIC        bool
	HTTPS       bool
	TLSNoVerify bool
	TLS         bool
	DNSCrypt    bool
	TCP         bool
	IPv6        bool
}

// Displays contains toggles for what to (and not to) display.
type Displays struct {
	Comments       bool
	Question       bool // QUESTION SECTION
	Opt            bool // OPT PSEUDOSECTION
	Answer         bool // ANSWER SECTION
	Authority      bool // AUTHORITY SECTION
	Additional     bool // ADDITIONAL SECTION
	Statistics     bool // Query time, message size, etc.
	UcodeTranslate bool // Translate Punycode back to Unicode
}

// EDNS contains toggles for various EDNS options.
type EDNS struct {
	Subnet     dns.EDNS0_SUBNET
	ZFlag      uint16
	BufSize    uint16
	EnableEDNS bool
	Cookie     bool
	DNSSEC     bool
	Expire     bool
	KeepOpen   bool
	Nsid       bool
	Padding    bool
	Version    uint8
}

// ParseSubnet takes a subnet argument and makes it into one that the DNS library
// understands.
func ParseSubnet(subnet string, opts *Options) error {
	ip, inet, err := net.ParseCIDR(subnet)
	if err != nil {
		// TODO: make not a default?
		if subnet == "0" {
			opts.EDNS.Subnet = dns.EDNS0_SUBNET{
				Code:          dns.EDNS0SUBNET,
				Family:        1,
				SourceNetmask: 0,
				SourceScope:   0,
				Address:       net.IPv4(0, 0, 0, 0),
			}

			return nil
		}

		return fmt.Errorf("EDNS subnet parsing: %w", err)
	}

	sub, _ := inet.Mask.Size()
	opts.EDNS.Subnet = dns.EDNS0_SUBNET{}
	opts.EDNS.Subnet.Address = ip
	opts.EDNS.Subnet.SourceNetmask = uint8(sub)

	switch ip.To4() {
	case nil:
		// Not a valid IPv4 so assume IPv6
		opts.EDNS.Subnet.Family = 2
	default:
		// Valid IPv4
		opts.EDNS.Subnet.Family = 1
	}

	return nil
}
