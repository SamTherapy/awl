// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"
	"testing"

	"github.com/stefansundin/go-zflag"
	"gotest.tools/v3/assert"
)

func TestMain(t *testing.T) { //nolint: paralleltest // Race conditions
	os.Stdout = os.NewFile(0, os.DevNull)
	os.Stderr = os.NewFile(0, os.DevNull)

	args := []string{"awl", "+yaml", "@1.1.1.1"}

	_, code, err := run(args)
	assert.NilError(t, err)
	assert.Equal(t, code, 0)

	args = []string{"awl", "+short", "@1.1.1.1"}

	_, code, err = run(args)
	assert.NilError(t, err)
	assert.Equal(t, code, 0)
}

func TestHelp(t *testing.T) {
	args := []string{"awl", "-h"}

	_, code, err := run(args)
	assert.ErrorIs(t, err, zflag.ErrHelp)
	assert.Equal(t, code, 1)
}
