package websocket

import (
	"github.com/json-iterator/go"
)

// message
type message struct {
	// Type
	// 0 - Command | 1 - Data
	Type int `json:"type"`
	// Command
	Command int `json:"command,omitempty"`
	// Data
	Data interface{} `json:"data,omitempty"`
}

// new message.
func newMessage(data []byte) (m *message, err error) {
	m = &message{}
	err = jsoniter.Unmarshal(data, m)

	return
}
