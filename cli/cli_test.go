// SPDX-License-Identifier: BSD-3-Clause

package cli_test

import (
	"os"
	"testing"
	"time"

	"git.froth.zone/sam/awl/cli"
	"gotest.tools/v3/assert"
)

// nolint: paralleltest
func TestEmpty(t *testing.T) {
	old := os.Args
	os.Args = []string{"awl", "-4"}
	opts, err := cli.ParseCLI("TEST")
	assert.NilError(t, err)
	assert.Equal(t, opts.Port, 53)
	assert.Assert(t, opts.IPv4)
	os.Args = old
}

// nolint: paralleltest
func TestTLSPort(t *testing.T) {
	old := os.Args
	os.Args = []string{"awl", "-T"}
	opts, err := cli.ParseCLI("TEST")
	assert.NilError(t, err)
	assert.Equal(t, opts.Port, 853)
	os.Args = old
}

// nolint: paralleltest
func TestSubnet(t *testing.T) {
	old := os.Args
	os.Args = []string{"awl", "--subnet", "127.0.0.1/32"}
	opts, err := cli.ParseCLI("TEST")
	assert.NilError(t, err)
	assert.Equal(t, opts.EDNS.Subnet.Family, uint16(1))

	os.Args = []string{"awl", "--subnet", "0"}
	opts, err = cli.ParseCLI("TEST")
	assert.NilError(t, err)
	assert.Equal(t, opts.EDNS.Subnet.Family, uint16(1))
	os.Args = old

	os.Args = []string{"awl", "--subnet", "::/0"}
	opts, err = cli.ParseCLI("TEST")
	assert.NilError(t, err)
	assert.Equal(t, opts.EDNS.Subnet.Family, uint16(2))
	os.Args = old
}

// nolint: paralleltest
func TestInvalidFlag(t *testing.T) {
	old := os.Args
	os.Args = []string{"awl", "--treebug"}
	_, err := cli.ParseCLI("TEST")
	assert.ErrorContains(t, err, "unknown flag")
	os.Args = old
}

// nolint: paralleltest
func TestInvalidDig(t *testing.T) {
	old := os.Args
	os.Args = []string{"awl", "+a"}
	_, err := cli.ParseCLI("TEST")
	assert.ErrorContains(t, err, "digflags: invalid argument")
	os.Args = old
}

// nolint: paralleltest
func TestVersion(t *testing.T) {
	old := os.Args
	os.Args = []string{"awl", "--version"}
	_, err := cli.ParseCLI("TEST")
	assert.ErrorType(t, err, cli.ErrNotError)
	os.Args = old
}

// nolint: paralleltest
func TestTimeout(t *testing.T) {
	args := [][]string{
		{"awl", "+timeout=0"},
		{"awl", "--timeout", "0"},
	}
	for _, test := range args {
		old := os.Args
		os.Args = test
		opt, err := cli.ParseCLI("TEST")
		assert.NilError(t, err)
		assert.Equal(t, opt.Request.Timeout, time.Second/2)
		os.Args = old
	}
}

// nolint: paralleltest
func TestRetries(t *testing.T) {
	args := [][]string{
		{"awl", "+retry=-2"},
		{"awl", "+tries=-2"},
		{"awl", "--retries", "-2"},
	}
	for _, test := range args {
		old := os.Args
		os.Args = test
		opt, err := cli.ParseCLI("TEST")
		assert.NilError(t, err)
		assert.Equal(t, opt.Request.Retries, 0)
		os.Args = old
	}
}

// nolint: paralleltest
func FuzzFlags(f *testing.F) {
	testcases := []string{"git.froth.zone", "", "!12345", "google.com.edu.org.fr"}
	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		os.Args = []string{"awl", orig}
		//nolint:errcheck // Only make sure the program does not crash
		cli.ParseCLI("TEST")
	})
}
