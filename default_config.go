package zap

import (
	"go.uber.org/zap"
)

var (
	// Default zap config.
	DefaultZapConfig = GetDefaultConfig()
)

// Get default zap config is in production mode, but enable console stdout.
func GetDefaultConfig() *Config {
	return &Config{
		Level: zap.NewAtomicLevel(),
		Development: false,
		Console:     true,
	}
}

// Get debug zap config is in debug mode and enable console stdout.
func GetDebugConfig() *Config {
	return &Config{
		Level: zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: true,
		Console:     true,
	}
}
