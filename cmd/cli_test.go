// SPDX-License-Identifier: BSD-3-Clause

package cli_test

import (
	"testing"
	"time"

	cli "git.froth.zone/sam/awl/cmd"
	"gotest.tools/v3/assert"
)

func TestEmpty(t *testing.T) {
	t.Parallel()

	args := []string{"awl", "-4"}

	opts, err := cli.ParseCLI(args, "TEST")
	assert.NilError(t, err)
	assert.Assert(t, opts.IPv4)
	assert.Equal(t, opts.Request.Port, 53)
}

func TestTLSPort(t *testing.T) {
	t.Parallel()

	args := []string{"awl", "-T"}

	opts, err := cli.ParseCLI(args, "TEST")
	assert.NilError(t, err)
	assert.Equal(t, opts.Request.Port, 853)
}

func TestValidSubnet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args []string
		want uint16
	}{
		{[]string{"awl", "--subnet", "127.0.0.1/32"}, uint16(1)},
		{[]string{"awl", "--subnet", "0"}, uint16(1)},
		{[]string{"awl", "--subnet", "::/0"}, uint16(2)},
	}

	for _, test := range tests {
		test := test

		t.Run(test.args[2], func(t *testing.T) {
			t.Parallel()

			opts, err := cli.ParseCLI(test.args, "TEST")

			assert.NilError(t, err)
			assert.Equal(t, opts.EDNS.Subnet.Family, test.want)
		})
	}
}

func TestInvalidSubnet(t *testing.T) {
	t.Parallel()

	args := []string{"awl", "--subnet", "/"}

	_, err := cli.ParseCLI(args, "TEST")
	assert.ErrorContains(t, err, "EDNS subnet")
}

func TestMBZ(t *testing.T) {
	t.Parallel()

	args := []string{"awl", "--zflag", "G"}

	_, err := cli.ParseCLI(args, "TEST")

	assert.ErrorContains(t, err, "EDNS MBZ")
}

func TestInvalidFlag(t *testing.T) {
	t.Parallel()

	args := []string{"awl", "--treebug"}

	_, err := cli.ParseCLI(args, "TEST")

	assert.ErrorContains(t, err, "unknown flag")
}

func TestInvalidDig(t *testing.T) {
	t.Parallel()

	args := []string{"awl", "+a"}

	_, err := cli.ParseCLI(args, "TEST")

	assert.ErrorContains(t, err, "digflags: invalid argument")
}

func TestVersion(t *testing.T) {
	t.Parallel()

	args := []string{"awl", "--version"}

	_, err := cli.ParseCLI(args, "test")

	assert.ErrorType(t, err, cli.ErrNotError)
}

func TestTimeout(t *testing.T) {
	t.Parallel()

	args := [][]string{
		{"awl", "+timeout=0"},
		{"awl", "--timeout", "0"},
	}
	for _, test := range args {
		test := test

		t.Run(test[1], func(t *testing.T) {
			t.Parallel()

			opt, err := cli.ParseCLI(test, "TEST")

			assert.NilError(t, err)
			assert.Equal(t, opt.Request.Timeout, time.Second/2)
		})
	}
}

func TestRetries(t *testing.T) {
	t.Parallel()

	args := [][]string{
		{"awl", "+retry=-2"},
		{"awl", "+tries=-2"},
		{"awl", "--retries", "-2"},
	}
	for _, test := range args {
		test := test

		t.Run(test[1], func(t *testing.T) {
			t.Parallel()

			opt, err := cli.ParseCLI(test, "TEST")

			assert.NilError(t, err)
			assert.Equal(t, opt.Request.Retries, 0)
		})
	}
}

func TestSetHTTPS(t *testing.T) {
	t.Parallel()

	args := [][]string{
		{"awl", "-H", "@dns.froth.zone/dns-query"},
		{"awl", "+https", "@dns.froth.zone"},
	}
	for _, test := range args {
		test := test

		t.Run(test[1], func(t *testing.T) {
			t.Parallel()

			opt, err := cli.ParseCLI(test, "TEST")

			assert.NilError(t, err)
			assert.Equal(t, opt.Request.Server, "dns.froth.zone")
			assert.Equal(t, opt.HTTPSOptions.Endpoint, "/dns-query")
		})
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
