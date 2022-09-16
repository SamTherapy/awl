// SPDX-License-Identifier: BSD-3-Clause

package cli_test

import (
	"strconv"
	"testing"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestParseArgs(t *testing.T) {
	t.Parallel()

	args := []string{
		"go.dev",
		"AAAA",
		"@1.1.1.1",
		"+ignore",
	}
	opts := new(util.Options)
	opts.Logger = util.InitLogger(0)
	err := cli.ParseMiscArgs(args, opts)
	assert.NilError(t, err)
	assert.Equal(t, opts.Request.Name, "go.dev.")
	assert.Equal(t, opts.Request.Type, dns.StringToType["AAAA"])
	assert.Equal(t, opts.Request.Server, "1.1.1.1")
	assert.Equal(t, opts.Truncate, true)
}

func TestParseNoInput(t *testing.T) {
	t.Parallel()

	args := []string{}
	opts := new(util.Options)
	opts.Logger = util.InitLogger(0)
	err := cli.ParseMiscArgs(args, opts)
	assert.NilError(t, err)
	assert.Equal(t, opts.Request.Name, ".")
	assert.Equal(t, opts.Request.Type, dns.StringToType["NS"])
}

func TestParseA(t *testing.T) {
	t.Parallel()

	args := []string{
		"golang.org.",
	}
	opts := new(util.Options)
	opts.Logger = util.InitLogger(0)
	err := cli.ParseMiscArgs(args, opts)
	assert.NilError(t, err)
	assert.Equal(t, opts.Request.Name, "golang.org.")
	assert.Equal(t, opts.Request.Type, dns.StringToType["A"])
}

func TestParsePTR(t *testing.T) {
	t.Parallel()

	args := []string{"8.8.8.8"}
	opts := new(util.Options)
	opts.Logger = util.InitLogger(0)
	opts.Reverse = true
	err := cli.ParseMiscArgs(args, opts)
	assert.NilError(t, err)
	assert.Equal(t, opts.Request.Type, dns.StringToType["PTR"])
}

func TestParseInvalidPTR(t *testing.T) {
	t.Parallel()

	args := []string{"8.88.8"}
	opts := new(util.Options)
	opts.Logger = util.InitLogger(0)
	opts.Reverse = true
	err := cli.ParseMiscArgs(args, opts)
	assert.ErrorContains(t, err, "unrecognized address")
}

func TestDefaultServer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"DNSCRYPT", "sdns://AQMAAAAAAAAAETk0LjE0MC4xNC4xNDo1NDQzINErR_JS3PLCu_iZEIbq95zkSV2LFsigxDIuUso_OQhzIjIuZG5zY3J5cHQuZGVmYXVsdC5uczEuYWRndWFyZC5jb20"},
		{"TLS", "dns.google"},
		{"HTTPS", "https://dns.cloudflare.com/dns-query"},
		{"QUIC", "dns.adguard.com"},
	}

	for _, test := range tests {
		test := test
		t.Run(test.in, func(t *testing.T) {
			t.Parallel()
			args := []string{}
			opts := new(util.Options)
			opts.Logger = util.InitLogger(0)
			switch test.in {
			case "DNSCRYPT":
				opts.DNSCrypt = true
			case "TLS":
				opts.TLS = true
			case "HTTPS":
				opts.HTTPS = true
			case "QUIC":
				opts.QUIC = true
			}
			err := cli.ParseMiscArgs(args, opts)
			assert.NilError(t, err)
			assert.Equal(t, opts.Request.Server, test.want)
		})
	}
}

func TestFlagSetting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in []string
	}{
		{[]string{"@sdns://AQMAAAAAAAAAETk0LjE0MC4xNC4xNDo1NDQzINErR_JS3PLCu_iZEIbq95zkSV2LFsigxDIuUso_OQhzIjIuZG5zY3J5cHQuZGVmYXVsdC5uczEuYWRndWFyZC5jb20"}},
		{[]string{"@tls://dns.google"}},
		{[]string{"@https://dns.cloudflare.com/dns-query"}},
		{[]string{"@quic://dns.adguard.com"}},
		{[]string{"@tcp://dns.froth.zone"}},
		{[]string{"@udp://dns.example.com"}},
	}

	for i, test := range tests {
		test := test
		i := i
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			opts := new(util.Options)
			opts.Logger = util.InitLogger(0)
			t.Parallel()
			err := cli.ParseMiscArgs(test.in, opts)
			assert.NilError(t, err)
			switch i {
			case 0:
				assert.Assert(t, opts.DNSCrypt)
			case 1:
				assert.Assert(t, opts.TLS)
			case 2:
				assert.Assert(t, opts.HTTPS)
			case 3:
				assert.Assert(t, opts.QUIC)
			case 4:
				assert.Assert(t, opts.TCP)
			case 5:
				assert.Assert(t, true)
			}
		})
	}
}

func FuzzParseArgs(f *testing.F) {
	cases := []string{
		"go.dev",
		"AAAA",
		"@1.1.1.1",
		"+ignore",
		"e",
	}

	for _, tc := range cases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, arg string) {
		// Get rid of outputs

		args := []string{arg}
		opts := new(util.Options)
		opts.Logger = util.InitLogger(0)
		//nolint:errcheck,gosec // Only make sure the program does not crash
		cli.ParseMiscArgs(args, opts)
	})
}
