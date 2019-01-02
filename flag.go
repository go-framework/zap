package zap

import "flag"

// Register config flag.
func (c *Config) RegisterFlag() {
	flag.BoolVar(&c.Development, "log-dev", false, "logger development")
	flag.BoolVar(&c.Console, "log-console", false, "logger console")
}
