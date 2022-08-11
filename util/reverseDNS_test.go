// SPDX-License-Identifier: BSD-3-Clause

package util_test

import (
	"testing"

	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

var (
	PTR   = dns.StringToType["PTR"]
	NAPTR = dns.StringToType["NAPTR"]
)

func TestIPv4(t *testing.T) {
	t.Parallel()

	act, err := util.ReverseDNS("8.8.4.4", PTR)
	assert.NilError(t, err)
	assert.Equal(t, act, "4.4.8.8.in-addr.arpa.", "IPv4 reverse")
}

func TestIPv6(t *testing.T) {
	t.Parallel()

	act, err := util.ReverseDNS("2606:4700:4700::1111", PTR)
	assert.NilError(t, err)
	assert.Equal(t, act, "1.1.1.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.7.4.0.0.7.4.6.0.6.2.ip6.arpa.", "IPv6 reverse")
}

func TestNAPTR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"1-800-555-1234", "4.3.2.1.5.5.5.0.0.8.1.e164.arpa."},
		{"+1 800555  1234", "4.3.2.1.5.5.5.0.0.8.1.e164.arpa."},
		{"+46766861004", "4.0.0.1.6.8.6.6.7.6.4.e164.arpa."},
		{"17705551212", "2.1.2.1.5.5.5.0.7.7.1.e164.arpa."},
	}
	for _, test := range tests {
		// Thanks Goroutines, very cool!
		test := test
		t.Run(test.in, func(t *testing.T) {
			t.Parallel()
			act, err := util.ReverseDNS(test.in, NAPTR)
			assert.NilError(t, err)
			assert.Equal(t, test.want, act)
		})
	}
}

func TestInvalid(t *testing.T) {
	t.Parallel()

	_, err := util.ReverseDNS("AAAAA", 1)
	assert.ErrorContains(t, err, "invalid value AAAAA given")
}

func TestInvalid2(t *testing.T) {
	t.Parallel()

	_, err := util.ReverseDNS("1.0", PTR)
	assert.ErrorContains(t, err, "PTR reverse")
}
