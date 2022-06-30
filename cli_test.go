// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"testing"

	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApp(t *testing.T) {
	app := prepareCLI()
	// What more can even be done lmao
	require.NotNil(t, app)
}

func TestArgParse(t *testing.T) {
	tests := []struct {
		in   []string
		want util.Answers
	}{
		{
			[]string{"@::1", "localhost", "AAAA"},
			util.Answers{Server: "::1", Request: dns.TypeAAAA, Name: "localhost"},
		},
		{
			[]string{"@1.0.0.1", "google.com"},
			util.Answers{Server: "1.0.0.1", Request: dns.TypeA, Name: "google.com"},
		},
		{
			[]string{"@8.8.4.4"},
			util.Answers{Server: "8.8.4.4", Request: dns.TypeNS, Name: "."},
		},
	}
	for _, test := range tests {
		act, err := parseArgs(test.in)
		assert.Nil(t, err)
		assert.Equal(t, test.want, act)
	}
}
