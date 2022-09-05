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

	old := os.Args

	os.Args = []string{"awl", "+yaml", "@1.1.1.1"}

	_, code, err := run()
	assert.NilError(t, err)
	assert.Equal(t, code, 0)

	os.Args = []string{"awl", "+short", "@1.1.1.1"}

	_, code, err = run()
	assert.NilError(t, err)
	assert.Equal(t, code, 0)

	os.Args = old
}

func TestHelp(t *testing.T) {
	old := os.Args
	os.Stdout = os.NewFile(0, os.DevNull)
	os.Stderr = os.NewFile(0, os.DevNull)

	os.Args = []string{"awl", "-h"}

	_, code, err := run()
	assert.ErrorIs(t, err, zflag.ErrHelp)
	assert.Equal(t, code, 1)

	os.Args = old
}
