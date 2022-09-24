// SPDX-License-Identifier: BSD-3-Clause

package logawl

import (
	"errors"
	"io"
	"sync"
	"sync/atomic"
)

type (
	// Level is the logging level.
	Level int32

	// Logger is the overall logger.
	Logger struct {
		Out       io.Writer
		Prefix    string
		buf       []byte
		Mu        sync.Mutex
		Level     Level
		isDiscard int32
	}
)

// SetLevel stores whatever input value is in mem address of l.level.
func (logger *Logger) SetLevel(level Level) {
	atomic.StoreInt32((*int32)(&logger.Level), int32(level))
}

// GetLevel gets the logger level.
func (logger *Logger) GetLevel() Level {
	return logger.level()
}

// Retrieves whatever was stored in mem address of l.level.
func (logger *Logger) level() Level {
	return Level(atomic.LoadInt32((*int32)(&logger.Level)))
}

// UnMarshalLevel unmarshalls the int value of level for writing the header.
func (logger *Logger) UnMarshalLevel(lv Level) (string, error) {
	switch lv {
	case ErrLevel:
		return "ERROR ", nil
	case WarnLevel:
		return "WARN ", nil
	case InfoLevel:
		return "INFO ", nil
	case DebugLevel:
		return "DEBUG ", nil
	}

	return "", errInvalidLevel
}

// IsLevel returns true if the logger level is above the level given.
func (logger *Logger) IsLevel(level Level) bool {
	return logger.level() >= level
}

// AllLevels is an array of all valid log levels.
var AllLevels = []Level{
	ErrLevel,
	WarnLevel,
	InfoLevel,
	DebugLevel,
}

const (
	// ErrLevel is the fatal (error) log level.
	ErrLevel Level = iota

	// WarnLevel is for warning logs.
	//
	// Example: when one setting implies another, when a request fails but is retried.
	WarnLevel

	// InfoLevel is for saying what is going on when.
	// This is essentially the "verbose" option.
	//
	// When in doubt, use info.
	InfoLevel

	// DebugLevel is for spewing debug structs/interfaces.
	DebugLevel
)

var errInvalidLevel = errors.New("invalid log level")
