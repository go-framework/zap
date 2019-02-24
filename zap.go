package zap

import "go.uber.org/zap"

// Defined as zap Logger.
type Logger = zap.Logger

// Defined as zap Option.
type Option = zap.Option

// new zap logger with config.
func NewZapLogger(config *Config, opts ...zap.Option) *zap.Logger {
	if config == nil {
		config = DefaultZapConfig
	}
	return config.NewZapLogger(opts...)
}
