// SPDX-License-Identifier: BSD-3-Clause
//go:build windows

package conf_test

import (
	"runtime"
	"testing"

	"dns.froth.zone/awl/conf"
	"gotest.tools/v3/assert"
)

func TestWinConfig(t *testing.T) {
	t.Parallel()

	if runtime.GOOS != "windows" {
		t.Skip("Not running Windows, skipping")
	}

	conf, err := conf.GetDNSConfig()

	assert.NilError(t, err)
	assert.Assert(t, len(conf.Servers) != 0)
}
