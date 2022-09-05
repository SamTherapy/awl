// SPDX-License-Identifier: BSD-3-Clause

package cli_test

import (
	"os"
	"testing"
	"time"

	"git.froth.zone/sam/awl/cli"
	"gotest.tools/v3/assert"
)

func TestEmpty(t *testing.T) {
	args := os.Args
	os.Args = []string{"awl", "-4"}

	opts, err := cli.ParseCLI("TEST")

	assert.NilError(t, err)
	assert.Equal(t, opts.Port, 53)
	assert.Assert(t, opts.IPv4)

	os.Args = args
}

func TestTLSPort(t *testing.T) {
	args := os.Args
	os.Args = []string{"awl", "-T"}

	opts, err := cli.ParseCLI("TEST")

	assert.NilError(t, err)
	assert.Equal(t, opts.Port, 853)

	os.Args = args
}

func TestSubnet(t *testing.T) {
	args := os.Args
	os.Args = []string{"awl", "--subnet", "127.0.0.1/32"}

	opts, err := cli.ParseCLI("TEST")

	assert.NilError(t, err)
	assert.Equal(t, opts.EDNS.Subnet.Family, uint16(1))

	os.Args = args

	os.Args = []string{"awl", "--subnet", "0"}

	opts, err = cli.ParseCLI("TEST")
	assert.NilError(t, err)
	assert.Equal(t, opts.EDNS.Subnet.Family, uint16(1))

	os.Args = args

	os.Args = []string{"awl", "--subnet", "::/0"}

	opts, err = cli.ParseCLI("TEST")
	assert.NilError(t, err)
	assert.Equal(t, opts.EDNS.Subnet.Family, uint16(2))

	os.Args = args

	os.Args = []string{"awl", "--subnet", "/"}

	opts, err = cli.ParseCLI("TEST")
	assert.ErrorContains(t, err, "EDNS subnet")

	os.Args = args
}

func TestMBZ(t *testing.T) { //nolint: paralleltest // Race conditions
	args := os.Args
	os.Args = []string{"awl", "--zflag", "G"}

	_, err := cli.ParseCLI("TEST")

	assert.ErrorContains(t, err, "EDNS MBZ")

	os.Args = args
}

func TestInvalidFlag(t *testing.T) { //nolint: paralleltest // Race conditions
	args := os.Args
	stdout := os.Stdout
	stderr := os.Stderr

	os.Stdout = os.NewFile(0, os.DevNull)
	os.Stderr = os.NewFile(0, os.DevNull)

	os.Args = []string{"awl", "--treebug"}

	_, err := cli.ParseCLI("TEST")

	assert.ErrorContains(t, err, "unknown flag")

	os.Args = args
	os.Stdout = stdout
	os.Stderr = stderr
}

func TestInvalidDig(t *testing.T) { //nolint: paralleltest // Race conditions
	args := os.Args
	os.Args = []string{"awl", "+a"}

	_, err := cli.ParseCLI("TEST")

	assert.ErrorContains(t, err, "digflags: invalid argument")

	os.Args = args
}

func TestVersion(t *testing.T) { //nolint: paralleltest // Race conditions
	args := os.Args
	stdout := os.Stdout
	stderr := os.Stderr

	os.Args = []string{"awl", "--version"}

	_, err := cli.ParseCLI("test")

	assert.ErrorType(t, err, cli.ErrNotError)

	os.Args = args
	os.Stdout = stdout
	os.Stderr = stderr
}

func TestTimeout(t *testing.T) { //nolint: paralleltest // Race conditions
	args := [][]string{
		{"awl", "+timeout=0"},
		{"awl", "--timeout", "0"},
	}
	for _, test := range args {
		args := os.Args
		os.Args = test

		opt, err := cli.ParseCLI("TEST")

		assert.NilError(t, err)
		assert.Equal(t, opt.Request.Timeout, time.Second/2)

		os.Args = args
	}
}

func TestRetries(t *testing.T) { //nolint: paralleltest // Race conditions
	args := [][]string{
		{"awl", "+retry=-2"},
		{"awl", "+tries=-2"},
		{"awl", "--retries", "-2"},
	}
	for _, test := range args {
		args := os.Args
		os.Args = test

		opt, err := cli.ParseCLI("TEST")

		assert.NilError(t, err)
		assert.Equal(t, opt.Request.Retries, 0)

		os.Args = args
	}
}

func FuzzFlags(f *testing.F) {
	testcases := []string{"git.froth.zone", "", "!12345", "google.com.edu.org.fr"}

	for _, tc := range testcases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		// Get rid of outputs

		args := os.Args
		os.Args = []string{"awl", orig}
		//nolint:errcheck,gosec // Only make sure the program does not crash
		cli.ParseCLI("TEST")
		os.Args = args
	})
}
