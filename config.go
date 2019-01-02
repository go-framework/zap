package zap

import (
	"os"
	"sort"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)


// Zap logger config.
// use lumberjack writing logs to rolling files.
// Level as AtomicLevel is an atomically changeable, dynamic logging level.
type Config struct {
	// lumberjack rotation config.
	Rotation *lumberjack.Logger `json:"rotation" yaml:"rotation"`
	// Logger text level.
	// level: debug, info, warn, error, dpanic, panic, and fatal.
	zap.AtomicLevel `json:"level" yaml:"level"`
	// logger development mode.
	Development bool `json:"development" yaml:"development"`
	// enable console logger.
	Console bool `json:"console" yaml:"console"`
}

// new zap logger.
func (c *Config) NewZapLogger(opts ...zap.Option) *zap.Logger {
	// zap config.
	var config zap.Config
	// zap encoder.
	var enc zapcore.Encoder

	// default stack level.
	stackLevel := zap.ErrorLevel

	// zap mode.
	if c.Development {
		opts = append(opts, zap.Development())
		config = zap.NewDevelopmentConfig()
		enc = zapcore.NewConsoleEncoder(config.EncoderConfig)
		stackLevel = zap.WarnLevel
	} else {
		config = zap.NewProductionConfig()
		enc = zapcore.NewJSONEncoder(config.EncoderConfig)
	}

	// enable caller.
	if !config.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}

	// enable stack trace .
	if !config.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}

	// config sampling.
	if config.Sampling != nil {
		opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewSampler(core, time.Second, int(config.Sampling.Initial), int(config.Sampling.Thereafter))
		}))
	}

	// initial fields.
	if len(config.InitialFields) > 0 {
		fs := make([]zap.Field, 0, len(config.InitialFields))
		keys := make([]string, 0, len(config.InitialFields))
		for k := range config.InitialFields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fs = append(fs, zap.Any(k, config.InitialFields[k]))
		}
		opts = append(opts, zap.Fields(fs...))
	}

	// multiple write syncer.
	var ws []zapcore.WriteSyncer

	// default is stdout.
	if c.Console || c.Rotation == nil {
		ws = append(ws, os.Stdout)
	}

	// enable rotation.
	if c.Rotation != nil {
		// append rotation syncer.
		ws = append(ws, zapcore.AddSync(c.Rotation))
	}

	// new zap core.
	core := zapcore.NewCore(
		enc,
		zapcore.NewMultiWriteSyncer(ws...),
		c.AtomicLevel,
	)

	// new zap logger.
	logger := zap.New(core)

	return logger.WithOptions(opts...)
}
