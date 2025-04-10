package sockets

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/timothychen1999/vrcontrol-server/model"
)

type Player struct {
	DeiviceID        string
	Connection       *websocket.Conn
	Room             *Room
	Stage            int
	InChannel        chan []byte
	Sequence         int
	LastUpdate       time.Time
	HeadPotion       model.Vector3f
	HeadForward      model.Vector3f
	LeftHandPostion  model.Vector3f
	LeftHandForward  model.Vector3f
	RightHandPostion model.Vector3f
	RightHandForward model.Vector3f
	LeftHandAvail    bool
	RightHandAvail   bool
}

func HandlePlayerConnect(room *Room, conn *websocket.Conn) *Player {
	player := Player{
		DeiviceID:  "device-id",
		Connection: conn,
		Room:       room,
	}
	player.InChannel = make(chan []byte, BufferSize)
	go player.read()
	go player.write()
	room.PlayerRegister <- &player
	return &player
}
func (p *Player) read() {
	defer func() {
		p.Room.PlayerUnregister <- p
		p.Connection.Close()
	}()
	p.Connection.SetReadLimit(MaxMessageSize)
	p.Connection.SetReadDeadline(time.Now().Add(PongWait))
	p.Connection.SetPongHandler(func(string) error { p.Connection.SetReadDeadline(time.Now().Add(PongWait)); return nil })
	for {
		_, message, err := p.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, Newline, Space, -1))
		var playerMessage model.PlayerMessage
		err = json.Unmarshal(message, &playerMessage)
		if err != nil {
			log.Println("Error Unmarshalling Message: ", err)
			continue
		}
		switch playerMessage.MessageType {
		case string(model.MessageTypeHeartbeat):
			heartbeat := playerMessage.Heartbeat
			p.HeadPotion = heartbeat.HeadPotion
			p.HeadForward = heartbeat.HeadForward
			p.LeftHandPostion = heartbeat.LeftHandPostion
			p.LeftHandForward = heartbeat.LeftHandForward
			p.RightHandPostion = heartbeat.RightHandPostion
			p.RightHandForward = heartbeat.RightHandForward
			p.LeftHandAvail = heartbeat.LeftHandAvail
			p.RightHandAvail = heartbeat.RightHandAvail
			p.Stage = heartbeat.Stage
			p.DeiviceID = heartbeat.DeviceID
			p.LastUpdate = time.Unix(heartbeat.Timestamp, 0)

		default:
			// Other is broadcast message
			// Send to the room
			p.Room.PlayerBroadcast <- message
		}

	}
}
func (p *Player) write() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		p.Connection.Close()
	}()

	for {
		select {
		case message, ok := <-p.InChannel:
			p.Connection.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				// The hub closed the channel.
				p.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				err := p.Connection.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("Error: ", err)
					break
				}
			}
		case <-ticker.C:
			p.Connection.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := p.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}
}
