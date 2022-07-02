package logawl

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var logger = New()

func TestLogawl(t *testing.T) {

	assert.Equal(t, Level(2), logger.Level) //cast 2 (int) to 2 (level)

	//Validate setting and getting levels from memory works
	for i := range AllLevels {
		logger.SetLevel(Level(i))
		assert.Equal(t, Level(i), logger.GetLevel())
	}

}

func TestUnmarshalLevels(t *testing.T) {
	m := make(map[int]string)
	var err error
	//Fill map with unmarshalled level info
	for i := range AllLevels {
		m[i], err = logger.UnMarshalLevel(Level(i))
		assert.Nil(t, err)
	}

	//iterate over map and assert equal
	for i := range AllLevels {
		lv, err := logger.UnMarshalLevel(Level(i))
		assert.Nil(t, err)
		assert.Equal(t, m[i], lv)
	}

	lv, err := logger.UnMarshalLevel(Level(9001))
	assert.NotNil(t, err)
	assert.Equal(t, "", lv)
	assert.ErrorContains(t, err, "invalid log level choice")
}

func TestLogger(t *testing.T) {

	for i := range AllLevels {
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
	ti := time.Now()
	test := []byte("test")
	assert.NotNil(t, logger.formatHeader(&test, ti, 0, Level(9001))) //make sure error is error

}
