// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"
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

func TestQuery(t *testing.T) {
	app := prepareCLI()
	args := os.Args[0:1]
	args = append(args, "--Treebug")
	err := app.Run(args)
	assert.NotNil(t, err)
}
func TestHTTPS(t *testing.T) {
	app := prepareCLI()
	args := os.Args[0:1]
	args = append(args, "-H")
	args = append(args, "@https://cloudflare-dns.com/dns-query")
	args = append(args, "git.froth.zone")
	err := app.Run(args)
	assert.Nil(t, err)
}

func FuzzCli(f *testing.F) {
	testcases := []string{"git.froth.zone", "", "!12345", "google.com.edu.org.fr"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		app := prepareCLI()
		args := os.Args[0:1]
		args = append(args, orig)
		err := app.Run(args)
		if err != nil {
			assert.ErrorContains(t, err, "domain must be fully qualified")
		}
		assert.Nil(t, err)
	})
}
