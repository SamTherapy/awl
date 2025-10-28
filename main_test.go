// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"testing"

	"github.com/spf13/pflag"
	"gotest.tools/v3/assert"
)

func TestRun(t *testing.T) {
	//	t.Parallel()
	args := [][]string{
		{"awl", "+yaml", "@1.1.1.1"},
		{"awl", "+short", "@1.1.1.1"},
	}

	for _, test := range args {
		test := test

		t.Run("", func(t *testing.T) {
			_, code, err := run(test)
			assert.NilError(t, err)
			assert.Equal(t, code, 0)
		})
	}
}

func TestTrace(t *testing.T) {
	domains := []string{"git.froth.zone", "google.com", "amazon.com", "freecumextremist.com", "dns.froth.zone", "sleepy.cafe", "pkg.go.dev"}

	for i := range domains {
		args := []string{"awl", "+trace", domains[i], "@1.1.1.1"}
		_, code, err := run(args)
		assert.NilError(t, err)
		assert.Equal(t, code, 0)
	}
}

func TestHelp(t *testing.T) {
	// t.Parallel()
	args := []string{"awl", "-h"}

	_, code, err := run(args)
	assert.ErrorIs(t, err, pflag.ErrHelp)
	assert.Equal(t, code, 1)
}
