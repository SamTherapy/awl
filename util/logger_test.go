// SPDX-License-Identifier: BSD-3-Clause

package util_test

import (
	"testing"

	"git.froth.zone/sam/awl/logawl"
	"git.froth.zone/sam/awl/util"
	"gotest.tools/v3/assert"
)

func TestInitLogger(t *testing.T) {
	logger := util.InitLogger(0)
	assert.Equal(t, logger.Level, logawl.Level(0))
}
