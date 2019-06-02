package zap

import (
	"fmt"
	"os"
	"sort"
	"time"
	"unsafe"

	"github.com/json-iterator/go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/go-framework/zap/syncer"
)

// Zap logger config.
// Level as AtomicLevel is an atomically changeable, dynamic logging level.
type Config struct {
	// Logger outs.
	// Writes: lumberjack
	Writes []*syncer.Write `json:"writes" yaml:"writes"`
	// Logger text level.
	// level: debug, info, warn, error, dpanic, panic, and fatal.
	Level zap.AtomicLevel `json:"level" yaml:"level"`
	// Logger development mode.
	Development bool `json:"development" yaml:"development"`
	// Enable console logger.
	Console bool `json:"console" yaml:"console"`
}

// Implement Stringer.
func (c *Config) String() string {
	data, err := jsoniter.Marshal(c)
	if err == nil {
		return *(*string)(unsafe.Pointer(&data))
	}
	return fmt.Sprintf("level: %s development: %t console: %t", c.Level.Level(), c.Development, c.Console)
}

// Clone.
func (c *Config) Clone() *Config {
	config := *c
	return &config
}

// New zap logger.
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

	// enable stdout.
	if c.Console {
		ws = append(ws, os.Stdout)
	}

	// enable Writes.
	if len(c.Writes) != 0 {
		// append writer.
		for _, writer := range c.Writes {
			ws = append(ws, zapcore.AddSync(writer.GetWriter()))
		}
	}

	// new zap core.
	core := zapcore.NewCore(
		enc,
		zapcore.NewMultiWriteSyncer(ws...),
		c.Level,
	)

	// new zap logger.
	logger := zap.New(core)

	return logger.WithOptions(opts...)
}

// Add syncer write.
func (c *Config) AddSyncerWrite(write *syncer.Write) {
	c.Writes = append(c.Writes, write)
}
