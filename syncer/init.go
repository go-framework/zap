package syncer

import (
	"github.com/go-framework/zap/syncer/lumberjack"
	"github.com/go-framework/zap/syncer/websocket"
)

// init for Register Writer.
func init() {
	// lumberjack
	RegisterWriter(lumberjack.Name, lumberjack.GetDefault())
	// websocket
	RegisterWriter(websocket.Name, websocket.GetDefault())
}
