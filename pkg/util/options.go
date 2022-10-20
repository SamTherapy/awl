// SPDX-License-Identifier: BSD-3-Clause

package util

import (
	"fmt"
	"net"

	"git.froth.zone/sam/awl/pkg/logawl"
	"github.com/miekg/dns"
)

// Options is the grand structure for all query options.
type Options struct {
	// The logger
	Logger *logawl.Logger `json:"-"`
	// Host to verify TLS cert with
	TLSHost string `json:"tlsHost" example:""`
	// EDNS Options
	EDNS

	// HTTPS options :)
	HTTPSOptions

	// DNS request :)
	Request

	// Verbosity levels, see [logawl.AllLevels]
	Verbosity int `json:"-" example:"0"`
	// Display options
	Display Display
	// Ignore Truncation
	Truncate bool `json:"ignoreTruncate" example:"false"`
	// Ignore BADCOOKIE
	BadCookie bool `json:"ignoreBadCookie" example:"false"`
	// Print only the answer
	Short bool `json:"short" example:"false"`
	// When Short is true, display where the query came from
	Identify bool `json:"identify" example:"false"`
	// Perform a reverse DNS query when true
	Reverse bool `json:"reverse" example:"false"`

	HeaderFlags

	// Display resposne as JSON
	JSON bool `json:"-" xml:"-" yaml:"-"`
	// Display response as XML
	XML bool `json:"-" xml:"-" yaml:"-"`
	// Display response as YAML
	YAML bool `json:"-" xml:"-" yaml:"-"`

	// Use TCP instead of UDP to make the query
	TCP bool `json:"tcp" example:"false"`
	// Use DNS-over-TLS to make the query
	TLS bool `json:"dnsOverTLS" example:"false"`
	// When using TLS, ignore certificates
	TLSNoVerify bool `json:"tlsNoVerify" example:"false"`
	// Use DNS-over-HTTPS to make the query
	HTTPS bool `json:"dnsOverHTTPS" example:"false"`
	// Use DNS-over-QUIC to make the query
	//nolint:tagliatelle // QUIC is an acronym
	QUIC bool `json:"dnsOverQUIC" example:"false"`
	// Use DNSCrypt to make the query
	DNSCrypt bool `json:"dnscrypt" example:"false"`

	// Force IPv4 only
	IPv4 bool `json:"forceIPv4" example:"false"`
	// Force IPv6 only
	IPv6 bool `json:"forceIPv6" example:"false"`
}

// HTTPSOptions are options exclusively for DNS-over-HTTPS queries.
type HTTPSOptions struct {
	// URL endpoint
	Endpoint string `json:"endpoint" example:"/dns-query"`

	// True, make GET request.
	// False, make POST request.
	Get bool `json:"get" example:"false"`
}

// HeaderFlags are the flags that are in DNS headers.
type HeaderFlags struct {
	// Authoritative Answer DNS query flag
	AA bool `json:"authoritative" example:"false"`
	// Authenticated Data DNS query flag
	AD bool `json:"authenticatedData" example:"false"`
	// Checking Disabled DNS query flag
	CD bool `json:"checkingDisabled" example:"false"`
	// QueRy DNS query flag
	QR bool `json:"query" example:"false"`
	// Recursion Desired DNS query flag
	RD bool `json:"recursionDesired" example:"true"`
	// Recursion Available DNS query flag
	RA bool `json:"recursionAvailable" example:"false"`
	// TrunCated DNS query flag
	TC bool `json:"truncated" example:"false"`
	// Zero DNS query flag
	Z bool `json:"zero" example:"false"`
}

// Display contains toggles for what to (and not to) display.
type Display struct {
	/* Section displaying */

	// Comments?
	Comments bool `json:"comments" example:"true"`
	// QUESTION SECTION
	Question bool `json:"question" example:"true"`
	// OPT PSEUDOSECTION
	Opt bool `json:"opt" example:"true"`
	// ANSWER SECTION
	Answer bool `json:"answer" example:"true"`
	// AUTHORITY SECTION
	Authority bool `json:"authority" example:"true"`
	// ADDITIONAL SECTION
	Additional bool `json:"additional" example:"true"`
	// Query time, message size, etc.
	Statistics bool `json:"statistics" example:"true"`
	// Display TTL in response
	TTL bool `json:"ttl" example:"true"`

	/* Answer formatting */

	// Display Class in response
	ShowClass bool `json:"showClass" example:"true"`
	// Display query before it is sent
	ShowQuery bool `json:"showQuery" example:"false"`
	// Display TTL as human-readable
	HumanTTL bool `json:"humanTTL" example:"false"`
	// Translate Punycode back to Unicode
	UcodeTranslate bool `json:"unicode" example:"true"`
}

// EDNS contains toggles for various EDNS options.
type EDNS struct {
	// Subnet to originate query from.
	Subnet dns.EDNS0_SUBNET `json:"subnet"`
	// Must Be Zero flag
	ZFlag uint16 `json:"zflag" example:"0"`
	// UDP buffer size
	BufSize uint16 `json:"bufSize" example:"1232"`
	// Enable/Disable EDNS entirely
	EnableEDNS bool `json:"edns" example:"false"`
	// Sending EDNS cookie
	Cookie bool `json:"cookie" example:"true"`
	// Enabling DNSSEC
	DNSSEC bool `json:"dnssec" example:"false"`
	// Sending EDNS Expire
	Expire bool `json:"expire" example:"false"`
	// Sending EDNS TCP keepopen
	KeepOpen bool `json:"keepOpen" example:"false"`
	// Sending EDNS NSID
	Nsid bool `json:"nsid" example:"false"`
	// Send EDNS Padding
	Padding bool `json:"padding" example:"false"`
	// Set EDNS version (default: 0)
	Version uint8 `json:"version" example:"0"`
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
