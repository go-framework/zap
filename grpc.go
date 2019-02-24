package zap

import (
	"context"

	"go.uber.org/zap/zapcore"
)

// Implement LevelServiceServer interface.

// Set logger atomic Level.
func (c *Config) SetLevel(ctx context.Context, level *AtomicLevel) (*Empty, error) {
	// set level.
	c.Level.SetLevel(zapcore.Level(level.Level - 1))
	return nil, nil
}

// Get logger atomic Level.
func (c *Config) GetLevel(context.Context, *Empty) (*AtomicLevel, error) {
	// get level.
	return &AtomicLevel{Level: Level(c.Level.Level() + 1)}, nil
}
