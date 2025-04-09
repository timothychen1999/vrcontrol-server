package sockets

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	DeiviceID  string
	Connection *websocket.Conn
	Room       *Room
	Stage      int
	InChannel  chan []byte
	OutChannel chan []byte
	Sequence   int
}

func HandlePlayerConnect(room *Room, conn *websocket.Conn) *Player {
	player := Player{
		DeiviceID:  "device-id",
		Connection: conn,
		Room:       room,
	}
	player.InChannel = make(chan []byte, 1024)

	return &player
}
