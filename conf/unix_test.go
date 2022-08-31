// SPDX-License-Identifier: BSD-3-Clause
//go:build unix || (!windows && !plan9 && !js && !zos)

// FIXME: Can remove the or on the preprocessor when Go 1.18 becomes obsolete

package conf_test

import (
	"runtime"
	"testing"

	"git.froth.zone/sam/awl/conf"
	"gotest.tools/v3/assert"
)

func TestNonWinConfig(t *testing.T) {
	if runtime.GOOS == "windows" || runtime.GOOS == "plan9" {
		t.Skip("Not running Unix-like, skipping")
	}

	conf, err := conf.GetDNSConfig()
	assert.NilError(t, err)
	assert.Assert(t, len(conf.Servers) != 0)
}
