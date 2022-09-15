// SPDX-License-Identifier: BSD-3-Clause

package cli_test

import (
	"testing"
	"time"

	"git.froth.zone/sam/awl/cli"
	"gotest.tools/v3/assert"
)

func TestEmpty(t *testing.T) {
	args := []string{"awl", "-4"}

	opts, err := cli.ParseCLI(args, "TEST")

	assert.NilError(t, err)
	assert.Equal(t, opts.Request.Port, 53)
	assert.Assert(t, opts.IPv4)
}

func TestTLSPort(t *testing.T) {
	args := []string{"awl", "-T"}

	opts, err := cli.ParseCLI(args, "TEST")

	assert.NilError(t, err)
	assert.Equal(t, opts.Request.Port, 853)
}

func TestSubnet(t *testing.T) {
	args := []string{"awl", "--subnet", "127.0.0.1/32"}

	opts, err := cli.ParseCLI(args, "TEST")

	assert.NilError(t, err)
	assert.Equal(t, opts.EDNS.Subnet.Family, uint16(1))

	args = []string{"awl", "--subnet", "0"}

	opts, err = cli.ParseCLI(args, "TEST")
	assert.NilError(t, err)
	assert.Equal(t, opts.EDNS.Subnet.Family, uint16(1))

	args = []string{"awl", "--subnet", "::/0"}

	opts, err = cli.ParseCLI(args, "TEST")
	assert.NilError(t, err)
	assert.Equal(t, opts.EDNS.Subnet.Family, uint16(2))

	args = []string{"awl", "--subnet", "/"}

	opts, err = cli.ParseCLI(args, "TEST")
	assert.ErrorContains(t, err, "EDNS subnet")
}

func TestMBZ(t *testing.T) {
	args := []string{"awl", "--zflag", "G"}

	_, err := cli.ParseCLI(args, "TEST")

	assert.ErrorContains(t, err, "EDNS MBZ")
}

func TestInvalidFlag(t *testing.T) {
	args := []string{"awl", "--treebug"}

	_, err := cli.ParseCLI(args, "TEST")

	assert.ErrorContains(t, err, "unknown flag")
}

func TestInvalidDig(t *testing.T) {
	args := []string{"awl", "+a"}

	_, err := cli.ParseCLI(args, "TEST")

	assert.ErrorContains(t, err, "digflags: invalid argument")
}

func TestVersion(t *testing.T) {
	args := []string{"awl", "--version"}

	_, err := cli.ParseCLI(args, "test")

	assert.ErrorType(t, err, cli.ErrNotError)
}

func TestTimeout(t *testing.T) {
	args := [][]string{
		{"awl", "+timeout=0"},
		{"awl", "--timeout", "0"},
	}
	for _, test := range args {
		test := test

		opt, err := cli.ParseCLI(test, "TEST")

		assert.NilError(t, err)
		assert.Equal(t, opt.Request.Timeout, time.Second/2)
	}
}

func TestRetries(t *testing.T) {
	args := [][]string{
		{"awl", "+retry=-2"},
		{"awl", "+tries=-2"},
		{"awl", "--retries", "-2"},
	}
	for _, test := range args {
		test := test

		opt, err := cli.ParseCLI(test, "TEST")

		assert.NilError(t, err)
		assert.Equal(t, opt.Request.Retries, 0)
	}
}

func FuzzFlags(f *testing.F) {
	testcases := []string{"git.froth.zone", "", "!12345", "google.com.edu.org.fr"}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		// Get rid of outputs

		args := []string{"awl", orig}
		//nolint:errcheck,gosec // Only make sure the program does not crash
		cli.ParseCLI(args, "TEST")
	})
}
