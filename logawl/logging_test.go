// SPDX-License-Identifier: BSD-3-Clause

package logawl_test

import (
	"bytes"
	"testing"
	"time"

	"git.froth.zone/sam/awl/logawl"

	"gotest.tools/v3/assert"
)

var logger = logawl.New()

func TestLogawl(t *testing.T) {
	t.Parallel()

	for i := range logawl.AllLevels {
		logger.SetLevel(logawl.Level(i))
		assert.Equal(t, logawl.Level(i), logger.GetLevel())
	}
}

func TestUnmarshalLevels(t *testing.T) {
	t.Parallel()
	m := make(map[int]string)
	var err error

	for i := range logawl.AllLevels {
		m[i], err = logger.UnMarshalLevel(logawl.Level(i))
		assert.NilError(t, err)
	}

	for i := range logawl.AllLevels {
		lv, err := logger.UnMarshalLevel(logawl.Level(i))
		assert.NilError(t, err)
		assert.Equal(t, m[i], lv)
	}

	lv, err := logger.UnMarshalLevel(logawl.Level(9001))
	assert.Equal(t, "", lv)
	assert.ErrorContains(t, err, "invalid log level")
}

func TestLogger(t *testing.T) {
	t.Parallel()

	for i := range logawl.AllLevels {
		// only test non-exiting log levels
		switch i {
		case 1:
			fn := func() {
				logger.Info("")
			}
			var buffer bytes.Buffer
			logger.Out = &buffer
			fn()
		case 2:
			fn := func() {
				logger.Info("Test")
			}
			var buffer bytes.Buffer
			logger.Out = &buffer
			fn()
		case 3:
			fn := func() {
				logger.Debug("Test")
			}
			var buffer bytes.Buffer
			logger.Out = &buffer
			fn()
		}
	}
}

func TestFmt(t *testing.T) {
	t.Parallel()
	ti := time.Now()
	test := []byte("test")
	assert.ErrorContains(t, logger.FormatHeader(&test, ti, 0, 9001), "invalid log level") //make sure error is error

}
