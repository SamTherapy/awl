// SPDX-License-Identifier: BSD-3-Clause
//go:build plan9

package conf_test

import (
	"runtime"
	"testing"

	"git.froth.zone/sam/awl/conf"
	"gotest.tools/v3/assert"
)

func TestPlan9Config(t *testing.T) {
	t.Parallel()

	if runtime.GOOS != "plan9" {
		t.Skip("Not running Plan 9, skipping")
	}

	conf, err := conf.GetDNSConfig()

	assert.NilError(t, err)
	assert.Assert(t, len(conf.Servers) != 0)
}
