package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIPv4(t *testing.T) {
	act, err := ReverseDNS("8.8.4.4", "PTR")
	assert.Nil(t, err)
	assert.Equal(t, act, "4.4.8.8.in-addr.arpa.", "IPv4 reverse")
}

// WIP
// func FuzzIPv4(f *testing.F) {
// 	f.Add("1.1.1.1", "1.1.1.1.in-addr.arpa")
// 	f.Fuzz(func(t *testing.T, a string, b string) {
// 		_, err := ReverseDNS(a, b)
// 		assert.Nil(t, err)
// 	})
// }

func TestIPv6(t *testing.T) {
	act, err := ReverseDNS("2606:4700:4700::1111", "PTR")
	assert.Nil(t, err)
	assert.Equal(t, act, "1.1.1.1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.7.4.0.0.7.4.6.0.6.2.ip6.arpa.", "IPv6 reverse")
}

func TestNAPTR(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"1-800-555-1234", "4.3.2.1.5.5.5.0.0.8.1.e164.arpa."},
		{"+1 800555  1234", "4.3.2.1.5.5.5.0.0.8.1.e164.arpa."},
		{"18005551234", "4.3.2.1.5.5.5.0.0.8.1.e164.arpa."},
		{"+1-770-555-1212", "2.1.2.1.5.5.5.0.7.7.1.e164.arpa."},
	}
	for _, test := range tests {
		act, err := ReverseDNS(test.in, "NAPTR")
		assert.Nil(t, err)
		assert.Equal(t, test.want, act)
	}
}

func TestInvalid(t *testing.T) {
	_, err := ReverseDNS("AAAAA", "A")
	assert.NotNil(t, err)
}
