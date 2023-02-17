// SPDX-License-Identifier: BSD-3-Clause

package util_test

import (
	"testing"

	"dns.froth.zone/awl/pkg/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestPTR(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		in       string
		expected string
	}{
		{
			"IPv4",
			"8.8.4.4", "4.4.8.8.in-addr.arpa.",
		},
		{
			"IPv6",
			"2606:4700:4700::1111", "1.1.1.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.7.4.0.0.7.4.6.0.6.2.ip6.arpa.",
		},
		{
			"Inavlid value",
			"AAAAA", "",
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			act, err := util.ReverseDNS(test.in, dns.StringToType["PTR"])
			if err == nil {
				assert.NilError(t, err)
			} else {
				assert.ErrorContains(t, err, "unrecognized address")
			}
			assert.Equal(t, act, test.expected)
		})
	}
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
			act, err := util.ReverseDNS(test.in, dns.StringToType["NAPTR"])
			assert.NilError(t, err)
			assert.Equal(t, test.want, act)
		})
	}
}

func TestInvalidAll(t *testing.T) {
	_, err := util.ReverseDNS("q", 15236)
	assert.ErrorContains(t, err, "invalid value")
}
