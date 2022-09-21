// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"testing"

	"github.com/stefansundin/go-zflag"
	"gotest.tools/v3/assert"
)

func TestRun(t *testing.T) {
	t.Parallel()

	args := [][]string{
		{"awl", "+yaml", "@1.1.1.1"},
		{"awl", "+short", "@1.1.1.1"},
	}

	for _, test := range args {
		test := test

		t.Run("", func(t *testing.T) {
			t.Parallel()
			_, code, err := run(test)
			assert.NilError(t, err)
			assert.Equal(t, code, 0)
		})
	}
}

func TestHelp(t *testing.T) {
	t.Parallel()

	args := []string{"awl", "-h"}

	_, code, err := run(args)
	assert.ErrorIs(t, err, zflag.ErrHelp)
	assert.Equal(t, code, 1)
}
