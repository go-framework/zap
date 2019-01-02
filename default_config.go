package zap

import (
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
)

var (
	// Get zap config function.
	GetZapConfigFunc = GetDefaultZapConfig
	// Default zap config.
	DefaultZapConfig = GetDefaultZapConfig()
	// Get lumberjack logger function.
	GetLumberjackLoggerFunc = GetDefaultLumberjackLogger
	// Default lumberjack logger.
	DefaultLumberjackLogger = GetDefaultLumberjackLogger()
)

// Get default lumberjack logger.
func GetDefaultLumberjackLogger() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s.log", os.Args[0]),
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     30,
		Compress:   true,
	}
}

// Get default zap config is in production mode, but enable console stdout.
func GetDefaultZapConfig() *Config {
	return &Config{
		Rotation:    nil,
		AtomicLevel: zap.NewAtomicLevel(),
		Development: false,
		Console:     true,
	}
}
