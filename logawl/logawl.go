// SPDX-License-Identifier: BSD-3-Clause

package logawl

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

type (
	Level  int32
	Logger struct {
		Mu        sync.Mutex
		Level     Level
		Prefix    string
		Out       io.Writer
		buf       []byte
		isDiscard int32
	}
)

// Stores whatever input value is in mem address of l.level.
func (l *Logger) SetLevel(level Level) {
	atomic.StoreInt32((*int32)(&l.Level), int32(level))
}

// Mostly nothing.
func (l *Logger) GetLevel() Level {
	return l.level()
}

// Retrieves whatever was stored in mem address of l.level.
func (l *Logger) level() Level {
	return Level(atomic.LoadInt32((*int32)(&l.Level)))
}

// Unmarshalls the int value of level for writing the header.
func (l *Logger) UnMarshalLevel(lv Level) (string, error) {
	switch lv {
	case 0:
		return "ERROR ", nil
	case 1:
		return "WARN ", nil
	case 2:
		return "INFO ", nil
	case 3:
		return "DEBUG ", nil
	}
	return "", fmt.Errorf("invalid log level")
}

func (l *Logger) IsLevel(level Level) bool {
	return l.level() >= level
}

var AllLevels = []Level{
	ErrLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
}

const (
	// Fatal logs (will call exit(1)).
	ErrLevel Level = iota

	// Error logs.
	WarnLevel

	// What is going on level.
	InfoLevel
	// Verbose log level.
	DebugLevel
)
