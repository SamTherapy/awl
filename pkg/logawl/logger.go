// SPDX-License-Identifier: BSD-3-Clause

package logawl

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// New instantiates Logger
//
// Level can be changed to one of the other log levels (ErrorLevel, WarnLevel, InfoLevel, DebugLevel).
func New() *Logger {
	return &Logger{
		Out:   os.Stderr,
		Level: WarnLevel, // Default value is WarnLevel
	}
}

// Println takes any and prints it out to Logger -> Out (io.Writer (default is std.Err)).
func (logger *Logger) Println(level Level, in ...any) {
	if atomic.LoadInt32(&logger.isDiscard) != 0 {
		return
	}
	// If verbose is not set --debug etc print _nothing_
	if logger.IsLevel(level) {
		switch level { // Goes through log levels and does stuff based on them (currently nothing)
		case ErrLevel:
			if err := logger.Printer(ErrLevel, fmt.Sprintln(in...)); err != nil {
				fmt.Fprintln(logger.Out, "Logger failed: ", err)
			}
		case WarnLevel:
			if err := logger.Printer(WarnLevel, fmt.Sprintln(in...)); err != nil {
				fmt.Fprintln(logger.Out, "Logger failed: ", err)
			}
		case InfoLevel:
			if err := logger.Printer(InfoLevel, fmt.Sprintln(in...)); err != nil {
				fmt.Fprintln(logger.Out, "Logger failed: ", err)
			}
		case DebugLevel:
			if err := logger.Printer(DebugLevel, fmt.Sprintln(in...)); err != nil {
				fmt.Fprintln(logger.Out, "Logger failed: ", err)
			}
		default:
			break
		}
	}
}

// FormatHeader formats the log header as such <LogLevel> YYYY/MM/DD HH:MM:SS (local time) <the message to log>.
func (logger *Logger) FormatHeader(buf *[]byte, t time.Time, line int, level Level) error {
	if lvl, err := logger.UnMarshalLevel(level); err == nil {
		// This is ugly but functional
		// maybe there can be an append func or something in the future
		*buf = append(*buf, lvl...)
		year, month, day := t.Date()

		*buf = append(*buf, '[')
		formatter(buf, year, 4)
		*buf = append(*buf, '/')
		formatter(buf, int(month), 2)
		*buf = append(*buf, '/')
		formatter(buf, day, 2)
		*buf = append(*buf, ' ')
		hour, min, sec := t.Clock()
		formatter(buf, hour, 2)
		*buf = append(*buf, ':')
		formatter(buf, min, 2)
		*buf = append(*buf, ':')
		formatter(buf, sec, 2)
		*buf = append(*buf, ']')
		*buf = append(*buf, ':')
		*buf = append(*buf, ' ')
	} else {
		return errInvalidLevel
	}

	return nil
}

// Printer prints the formatted message directly to stdErr.
func (logger *Logger) Printer(level Level, s string) error {
	now := time.Now()

	var line int

	logger.Mu.Lock()
	defer logger.Mu.Unlock()

	logger.buf = logger.buf[:0]

	if err := logger.FormatHeader(&logger.buf, now, line, level); err != nil {
		return err
	}

	logger.buf = append(logger.buf, s...)

	if len(s) == 0 || s[len(s)-1] != '\n' {
		logger.buf = append(logger.buf, '\n')
	}

	_, err := logger.Out.Write(logger.buf)
	if err != nil {
		return fmt.Errorf("logger printing: %w", err)
	}

	return nil
}

// Some line formatting stuff from Golang log stdlib file
//
// Please view
// https://cs.opensource.google/go/go/+/refs/tags/go1.19:src/log/log.go;drc=41e1d9075e428c2fc32d966b3752a3029b620e2c;l=96
//
// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func formatter(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1

	for i >= 10 || wid > 1 {
		wid--

		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--

		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// Debug calls print directly with Debug level.
func (logger *Logger) Debug(in ...any) {
	logger.Println(DebugLevel, in...)
}

// Debugf calls print after formatting the string with Debug level.
func (logger *Logger) Debugf(format string, in ...any) {
	logger.Println(ErrLevel, fmt.Sprintf(format, in...))
}

// Info calls print directly with Info level.
func (logger *Logger) Info(in ...any) {
	logger.Println(InfoLevel, in...)
}

// Infof calls print after formatting the string with Info level.
func (logger *Logger) Infof(format string, in ...any) {
	logger.Println(ErrLevel, fmt.Sprintf(format, in...))
}

// Warn calls print directly with Warn level.
func (logger *Logger) Warn(in ...any) {
	logger.Println(WarnLevel, in...)
}

// Warnf calls print after formatting the string with Warn level.
func (logger *Logger) Warnf(format string, in ...any) {
	logger.Println(WarnLevel, fmt.Sprintf(format, in...))
}

// Error calls print directly with Error level.
func (logger *Logger) Error(in ...any) {
	logger.Println(ErrLevel, in...)
}

// Errorf calls print after formatting the string with Error level.
func (logger *Logger) Errorf(format string, in ...any) {
	logger.Println(ErrLevel, fmt.Sprintf(format, in...))
}
