package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	zapConfig "github.com/go-framework/zap"
)

// default zap logger.
var defaultLogger = zapConfig.NewZapLogger(zapConfig.GetDebugConfig(), zap.AddCallerSkip(1))

// Set zap logger.
func Set(logger *zap.Logger) {
	if logger == nil {
		return
	}
	defaultLogger = logger.WithOptions(zap.AddCallerSkip(1))
}

func Get() *zap.Logger {
	return defaultLogger
}

// Sugar wraps the Logger to provide a more ergonomic, but slightly slower,
// API. Sugaring a Logger is quite inexpensive, so it's reasonable for a
// single application to use both Loggers and SugaredLoggers, converting
// between them on the boundaries of performance-sensitive code.
func Sugar() *zap.SugaredLogger {
	return defaultLogger.Sugar()
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func Named(s string) *zap.Logger {
	return defaultLogger.Named(s)
}

// WithOptions clones the current Logger, applies the supplied Options, and
// returns the resulting Logger. It's safe to use concurrently.
func WithOptions(opts ...zap.Option) *zap.Logger {
	return defaultLogger.WithOptions(opts...)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func With(fields ...zap.Field) *zap.Logger {
	return defaultLogger.With(fields...)
}

// Check returns a CheckedEntry if logging a message at the specified level
// is enabled. It's a completely optional optimization; in high-performance
// applications, Check can help avoid allocating a slice to hold fields.
func Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return defaultLogger.Check(lvl, msg)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.InfoLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.WarnLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.ErrorLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func DPanic(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.DPanicLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.PanicLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zap.Field) {
	if ce := defaultLogger.Check(zap.FatalLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func Sync() error {
	return defaultLogger.Sync()
}
