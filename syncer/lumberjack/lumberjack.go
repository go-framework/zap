package lumberjack

import (
	"fmt"
	"io"
	"os"

	"github.com/natefinch/lumberjack"
)

const (
	Name = "lumberjack"
)

// Redefined lumberjack Logger.
type Logger struct {
	lumberjack.Logger `json:",inline" yaml:",inline" mapstructure:",squash"`
}

// Implement Cloner interface.
func (l *Logger) Clone() io.Writer {
	n := *l

	return &n
}

// Get default logger.
func GetDefault() *Logger {
	l := lumberjack.Logger{
		Filename:   fmt.Sprintf("%s.log", os.Args[0]),
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}

	return &Logger{
		Logger: l,
	}
}

// New logger with filename.
func New(filename string) *Logger {
	l := lumberjack.Logger{
		Filename:   filename,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}

	return &Logger{
		Logger: l,
	}
}
