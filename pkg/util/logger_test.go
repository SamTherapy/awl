// SPDX-License-Identifier: BSD-3-Clause

package util_test

import (
	"testing"

	"dns.froth.zone/awl/pkg/logawl"
	"dns.froth.zone/awl/pkg/util"
	"gotest.tools/v3/assert"
)

func TestInitLogger(t *testing.T) {
	t.Parallel()

	logger := util.InitLogger(0)
	assert.Equal(t, logger.Level, logawl.Level(0))
}
