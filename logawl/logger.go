// SPDX-License-Identifier: BSD-3-Clause

package logawl

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

// Calling New instantiates Logger
//
// Level can be changed to one of the other log levels (ErrorLevel, WarnLevel, InfoLevel, DebugLevel).
func New() *Logger {
	return &Logger{
		Out:   os.Stderr,
		Level: WarnLevel, // Default value is WarnLevel
	}
}

// Takes any and prints it out to Logger -> Out (io.Writer (default is std.Err)).
func (l *Logger) Println(level Level, v ...any) {
	if atomic.LoadInt32(&l.isDiscard) != 0 {
		return
	}
	// If verbose is not set --debug etc print _nothing_
	if l.IsLevel(level) {
		switch level { // Goes through log levels and does stuff based on them (currently nothing)
		case 0:
			err := l.Printer(0, fmt.Sprintln(v...)) // Error level
			if err != nil {
				fmt.Fprintln(os.Stderr, "Logger failed: ", err)
			}
		case 1:
			err := l.Printer(1, fmt.Sprintln(v...)) // Warn level
			if err != nil {
				fmt.Fprintln(os.Stderr, "Logger failed: ", err)
			}
		case 2:
			err := l.Printer(2, fmt.Sprintln(v...)) // Info level
			if err != nil {
				fmt.Fprintln(os.Stderr, "Logger failed: ", err)
			}
		case 3:
			err := l.Printer(3, fmt.Sprintln(v...)) // Debug level
			if err != nil {
				fmt.Fprintln(os.Stderr, "Logger failed: ", err)
			}
		default:
			break
		}
	}
}

// Formats the log header as such <LogLevel> YYYY/MM/DD HH:MM:SS (local time) <the message to log>.
func (l *Logger) FormatHeader(buf *[]byte, t time.Time, line int, level Level) error {
	if lvl, err := l.UnMarshalLevel(level); err == nil {
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
func (l *Logger) Printer(level Level, s string) error {
	now := time.Now()
	var line int
	l.Mu.Lock()
	defer l.Mu.Unlock()

	l.buf = l.buf[:0]
	err := l.FormatHeader(&l.buf, now, line, level)
	if err != nil {
		return err
	}
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err = l.Out.Write(l.buf)
	if err != nil {
		return fmt.Errorf("logger printing error %w", err)
	}
	return nil
}

// Some line formatting stuff from Golang log stdlib file
//
// Please view https://cs.opensource.google/go/go/+/refs/tags/go1.18.3:src/log/log.go;drc=41e1d9075e428c2fc32d966b3752a3029b620e2c;l=96
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

// Call print directly with Debug level.
func (l *Logger) Debug(v ...any) {
	l.Println(DebugLevel, v...)
}

// Call print directly with Info level.
func (l *Logger) Info(v ...any) {
	l.Println(InfoLevel, v...)
}

// Call print directly with Warn level.
func (l *Logger) Warn(v ...any) {
	l.Println(WarnLevel, v...)
}

// Call print directly with Error level.
func (l *Logger) Error(v ...any) {
	l.Println(ErrLevel, v...)
}
