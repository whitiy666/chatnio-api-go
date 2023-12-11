package utils

import (
	"testing"
)

func TestNewWebsocket(t *testing.T) {
	ws, err := NewWebsocket("wss://api.chatnio.net/chat")
	if err != nil || ws == nil {
		t.Error("websocket is nil")
	}
}
