// SPDX-License-Identifier: BSD-3-Clause

package query_test

import (
	"testing"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/internal/helpers"
	"git.froth.zone/sam/awl/query"
	"git.froth.zone/sam/awl/util"

	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestBadFormat(t *testing.T) {
	t.Parallel()
	_, err := query.PrintSpecial(new(dns.Msg), cli.Options{})
	assert.ErrorContains(t, err, "never happen")
}

func TestPrinting(t *testing.T) {
	t.Parallel()
	opts := cli.Options{
		Logger:    util.InitLogger(0),
		Port:      53,
		IPv4:      false,
		IPv6:      false,
		TCP:       true,
		DNSCrypt:  false,
		TLS:       false,
		HTTPS:     false,
		QUIC:      false,
		Truncate:  false,
		ShowQuery: true,
		AA:        false,
		AD:        false,
		CD:        false,
		QR:        false,
		RD:        true,
		RA:        false,
		TC:        false,
		Z:         false,
		Reverse:   false,
		Verbosity: 0,
		HumanTTL:  false,
		ShowTTL:   true,
		Short:     false,
		Identify:  false,
		JSON:      true,
		XML:       false,
		YAML:      false,
		Display: cli.Displays{
			Comments:       true,
			Question:       true,
			Answer:         true,
			Authority:      true,
			Additional:     true,
			Statistics:     true,
			UcodeTranslate: true,
		},
		Request: helpers.Request{
			Server:  "a.gtld-servers.net",
			Type:    dns.StringToType["NS"],
			Class:   1,
			Name:    "google.com.",
			Timeout: 0,
			Retries: 0,
		},
		EDNS: cli.EDNS{
			EnableEDNS: false,
		},
	}

	resp, err := query.CreateQuery(opts)
	assert.NilError(t, err)

	str, err := query.PrintSpecial(resp.DNS, opts)
	assert.NilError(t, err)
	assert.Assert(t, str != "")
}

func TestPrinting2(t *testing.T) {
	t.Parallel()
	opts := cli.Options{
		Logger:    util.InitLogger(0),
		Port:      53,
		IPv4:      false,
		IPv6:      false,
		TCP:       true,
		DNSCrypt:  false,
		TLS:       false,
		HTTPS:     false,
		QUIC:      false,
		Truncate:  false,
		ShowQuery: true,
		AA:        false,
		AD:        false,
		CD:        false,
		QR:        false,
		RD:        true,
		RA:        false,
		TC:        false,
		Z:         false,
		Reverse:   false,
		Verbosity: 0,
		HumanTTL:  false,
		ShowTTL:   true,
		Short:     true,
		Identify:  true,
		JSON:      false,
		XML:       false,
		YAML:      true,
		Display: cli.Displays{
			Comments:       true,
			Question:       true,
			Answer:         true,
			Authority:      true,
			Additional:     true,
			Statistics:     true,
			UcodeTranslate: true,
		},
		Request: helpers.Request{
			Server:  "ns1.google.com",
			Type:    dns.StringToType["NS"],
			Class:   1,
			Name:    "google.com.",
			Timeout: 0,
			Retries: 0,
		},
		EDNS: cli.EDNS{
			EnableEDNS: false,
		},
	}

	resp, err := query.CreateQuery(opts)
	assert.NilError(t, err)

	str, err := query.PrintSpecial(resp.DNS, opts)
	assert.NilError(t, err)
	assert.Assert(t, str != "")

	str = query.ToString(resp, opts)
	assert.Assert(t, str != "")
}

func TestPrinting3(t *testing.T) {
	t.Parallel()
	opts := cli.Options{
		Logger:    util.InitLogger(0),
		Port:      53,
		IPv4:      false,
		IPv6:      false,
		TCP:       false,
		DNSCrypt:  false,
		TLS:       false,
		HTTPS:     true,
		QUIC:      false,
		Truncate:  false,
		ShowQuery: true,
		AA:        false,
		AD:        false,
		CD:        false,
		QR:        false,
		RD:        true,
		RA:        false,
		TC:        false,
		Z:         false,
		Reverse:   false,
		Verbosity: 0,
		HumanTTL:  false,
		ShowTTL:   true,
		Short:     false,
		Identify:  true,
		JSON:      false,
		XML:       false,
		YAML:      true,
		Display: cli.Displays{
			Comments:       true,
			Question:       true,
			Answer:         true,
			Authority:      true,
			Additional:     true,
			Statistics:     true,
			UcodeTranslate: true,
		},
		Request: helpers.Request{
			Server:  "https://dns.froth.zone/dns-query",
			Type:    dns.StringToType["NS"],
			Class:   1,
			Name:    "freecumextremist.com.",
			Timeout: 0,
			Retries: 0,
		},
		EDNS: cli.EDNS{
			EnableEDNS: false,
		},
	}

	resp, err := query.CreateQuery(opts)
	assert.NilError(t, err)

	str, err := query.PrintSpecial(resp.DNS, opts)
	assert.NilError(t, err)
	assert.Assert(t, str != "")

	str = query.ToString(resp, opts)
	assert.Assert(t, str != "")
}

func TestPrinting4(t *testing.T) {
	t.Parallel()
	opts := cli.Options{
		Logger:    util.InitLogger(0),
		Port:      853,
		IPv4:      false,
		IPv6:      false,
		TCP:       false,
		DNSCrypt:  false,
		TLS:       true,
		HTTPS:     false,
		QUIC:      false,
		Truncate:  false,
		ShowQuery: true,
		AA:        false,
		AD:        false,
		CD:        false,
		QR:        false,
		RD:        true,
		RA:        false,
		TC:        false,
		Z:         false,
		Reverse:   false,
		Verbosity: 0,
		HumanTTL:  false,
		ShowTTL:   true,
		Short:     false,
		Identify:  true,
		JSON:      false,
		XML:       false,
		YAML:      true,
		Display: cli.Displays{
			Comments:       true,
			Question:       true,
			Answer:         true,
			Authority:      true,
			Additional:     true,
			Statistics:     true,
			UcodeTranslate: true,
		},
		Request: helpers.Request{
			Server:  "dns.google",
			Type:    dns.StringToType["NS"],
			Class:   1,
			Name:    "freecumextremist.com.",
			Timeout: 0,
			Retries: 0,
		},
		EDNS: cli.EDNS{
			EnableEDNS: false,
		},
	}

	resp, err := query.CreateQuery(opts)
	assert.NilError(t, err)

	str, err := query.PrintSpecial(resp.DNS, opts)
	assert.NilError(t, err)
	assert.Assert(t, str != "")

	str = query.ToString(resp, opts)
	assert.Assert(t, str != "")
}

func TestPrinting5(t *testing.T) {
	t.Parallel()
	opts := cli.Options{
		Logger:    util.InitLogger(0),
		Port:      53,
		IPv4:      false,
		IPv6:      false,
		TCP:       true,
		DNSCrypt:  false,
		TLS:       false,
		HTTPS:     false,
		QUIC:      false,
		Truncate:  false,
		ShowQuery: true,
		AA:        true,
		AD:        false,
		CD:        false,
		QR:        false,
		RD:        true,
		RA:        false,
		TC:        false,
		Z:         false,
		Reverse:   false,
		Verbosity: 0,
		HumanTTL:  false,
		ShowTTL:   true,
		Short:     false,
		Identify:  false,
		JSON:      false,
		XML:       false,
		YAML:      true,
		Display: cli.Displays{
			Comments:       true,
			Question:       true,
			Answer:         true,
			Authority:      true,
			Additional:     true,
			Statistics:     true,
			UcodeTranslate: true,
		},
		Request: helpers.Request{
			Server:  "rin.froth.zone",
			Type:    dns.StringToType["A"],
			Class:   1,
			Name:    "froth.zone.",
			Timeout: 0,
			Retries: 0,
		},
		EDNS: cli.EDNS{
			EnableEDNS: true,
			Cookie:     true,
		},
	}

	resp, err := query.CreateQuery(opts)
	assert.NilError(t, err)

	str, err := query.PrintSpecial(resp.DNS, opts)
	assert.NilError(t, err)
	assert.Assert(t, str != "")

	str = query.ToString(resp, opts)
	assert.Assert(t, str != "")
}

func TestToString6(t *testing.T) {
	assert.Assert(t, query.ToString(*new(helpers.Response), *new(cli.Options)) == "<nil> MsgHdr")
}
