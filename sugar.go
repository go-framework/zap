package zap

import "go.uber.org/zap"

// Default sugared logger.
type SugaredLogger struct {
	*zap.SugaredLogger
}

// new sugared logger by default zap config.
func NewSugaredLogger(opts ...zap.Option) *SugaredLogger {
	return &SugaredLogger{
		SugaredLogger: DefaultZapConfig.NewZapLogger(opts...).Sugar(),
	}
}

// wrap args split by space
func wrap(args []interface{}) (a []interface{}) {
	for n, arg := range args {
		if len(args) > 0 && n < len(args)-1 {
			a = append(a, arg, " ")
		} else {
			a = append(a, arg)
		}
	}
	return
}

// Print uses fmt.Sprint to construct and log a message.
func (l *SugaredLogger) Print(args ...interface{}) {
	l.SugaredLogger.Info(args...)
}

// Printf uses fmt.Sprintf to log a templated message.
func (l *SugaredLogger) Printf(template string, args ...interface{}) {
	l.SugaredLogger.Infof(template, args...)
}

// Debugln uses fmt.Sprint to construct and log a message.
func (l *SugaredLogger) Debugln(args ...interface{}) {
	l.SugaredLogger.Debug(wrap(args)...)
}

// Infoln uses fmt.Sprint to construct and log a message.
func (l *SugaredLogger) Infoln(args ...interface{}) {
	l.SugaredLogger.Info(wrap(args)...)
}

// Println uses fmt.Sprint to construct and log a message.
func (l *SugaredLogger) Println(args ...interface{}) {
	l.SugaredLogger.Info(wrap(args)...)
}

// Warnln uses fmt.Sprint to construct and log a message.
func (l *SugaredLogger) Warnln(args ...interface{}) {
	l.SugaredLogger.Warn(wrap(args)...)
}

// Errorln uses fmt.Sprint to construct and log a message.
func (l *SugaredLogger) Errorln(args ...interface{}) {
	l.SugaredLogger.Error(wrap(args)...)
}

// DPanicln uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *SugaredLogger) DPanicln(args ...interface{}) {
	l.SugaredLogger.DPanic(wrap(args)...)
}

// Panicln uses fmt.Sprint to construct and log a message, then panics.
func (l *SugaredLogger) Panicln(args ...interface{}) {
	l.SugaredLogger.Panic(wrap(args)...)
}

// Fatalln uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *SugaredLogger) Fatalln(args ...interface{}) {
	l.SugaredLogger.Fatal(wrap(args)...)
}

// With adds a variadic number of fields to the logging context. It accepts a
// mix of strongly-typed Field objects and loosely-typed key-value pairs. When
// processing pairs, the first element of the pair is used as the field key
// and the second as the field value.
//
// For example,
//   sugaredLogger.With(
//     "hello", "world",
//     "failure", errors.New("oh no"),
//     Stack(),
//     "count", 42,
//     "user", User{Name: "alice"},
//  )
// is the equivalent of
//   unsugared.With(
//     String("hello", "world"),
//     String("failure", "oh no"),
//     Stack(),
//     Int("count", 42),
//     Object("user", User{Name: "alice"}),
//   )
//
// Note that the keys in key-value pairs should be strings. In development,
// passing a non-string key panics. In production, the logger is more
// forgiving: a separate error is logged, but the key-value pair is skipped
// and execution continues. Passing an orphaned key triggers similar behavior:
// panics in development and errors in production.
func (l *SugaredLogger) With(args ...interface{}) *SugaredLogger {
	return &SugaredLogger{SugaredLogger: l.SugaredLogger.With(args...)}
}
