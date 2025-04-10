package sockets

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	WriteWait = 5 * time.Second

	// Time allowed to read the next pong message from the peer.
	PongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	PingPeriod = (PongWait * 9) / 10

	// Maximum message size allowed from peer.
	MaxMessageSize = 1024

	BufferSize = MaxMessageSize * 32

	//Tick Per Second
	TickRate = 1
)

var (
	Newline = []byte{'\n'}
	Space   = []byte{' '}
)

var SocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  BufferSize,
	WriteBufferSize: BufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
