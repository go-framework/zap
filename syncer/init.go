package syncer

import (
	"github.com/go-framework/zap/syncer/lumberjack"
)

// init for Register Writer.
func init() {
	RegisterWriter(lumberjack.Name, lumberjack.GetDefault())
}
