package zap

import "go.uber.org/zap"

// new zap logger with config.
func NewZapLogger(config *Config, opts ...zap.Option) *zap.Logger {
	if config == nil {
		config = DefaultZapConfig
	}
	return config.NewZapLogger(opts...)
}
