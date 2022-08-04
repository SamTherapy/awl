// SPDX-License-Identifier: BSD-3-Clause

package cli

import (
	"errors"
	"fmt"
	"net"

	"git.froth.zone/sam/awl/internal/helpers"
	"git.froth.zone/sam/awl/logawl"
	"github.com/miekg/dns"
)

// CLI options structure.
type Options struct {
	Logger    *logawl.Logger // Logger
	Port      int            // DNS port
	IPv4      bool           // Force IPv4
	IPv6      bool           // Force IPv6
	TCP       bool           // Query with TCP
	DNSCrypt  bool           // Query over DNSCrypt
	TLS       bool           // Query over TLS
	HTTPS     bool           // Query over HTTPS
	QUIC      bool           // Query over QUIC
	Truncate  bool           // Ignore truncation
	ShowQuery bool           // Show query before being sent
	AA        bool           // Set Authoratative Answer
	AD        bool           // Set Authenticated Data
	CD        bool           // Set CD
	QR        bool           // Set QueRy
	RD        bool           // Set Recursion Desired
	RA        bool           // Set Recursion Available
	TC        bool           // Set TC (TrunCated)
	Z         bool           // Set Z (Zero)
	Reverse   bool           // Make reverse query
	Verbosity int            // Set logawl verbosity
	HumanTTL  bool           // Make TTL human readable
	ShowTTL   bool           // Display TTL
	Short     bool           // Short output
	Identify  bool           // If short, add identity stuffs
	JSON      bool           // Outout as JSON
	XML       bool           // Output as XML
	YAML      bool           // Output at YAML

	Display Displays        // Display options
	Request helpers.Request // DNS reuqest
	EDNS                    // EDNS
}

// What to (and not to) display.
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

type EDNS struct {
	EnableEDNS bool             // Enable EDNS
	Cookie     bool             // Enable EDNS cookie
	DNSSEC     bool             // Enable DNSSEC
	BufSize    uint16           // Set UDP buffer size
	Version    uint8            // Set EDNS version
	Expire     bool             // Set EDNS expiration
	KeepOpen   bool             // TCP keep alive
	Nsid       bool             // Show EDNS nsid
	ZFlag      uint16           // EDNS flags
	Padding    bool             // EDNS padding
	Subnet     dns.EDNS0_SUBNET // EDNS Subnet (duh)
}

// parseSubnet takes a subnet argument and makes it into one that the DNS library
// understands.
func parseSubnet(subnet string, opts *Options) error {
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
		return fmt.Errorf("subnet parsing error %w", err)
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

var ErrNotError = errors.New("not an error")

var errNoArg = errors.New("no argument given")

type errInvalidArg struct {
	arg string
}

func (e *errInvalidArg) Error() string {
	return fmt.Sprintf("digflags: invalid argument %s", e.arg)
}
