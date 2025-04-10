package sockets

import (
	"encoding/json"
	"log"
	"time"

	"github.com/timothychen1999/vrcontrol-server/model"
)

type MessageType string

const (
	MessageTypeUpdate MessageType = "update"
)

type Room struct {
	RoomID           string
	PlayerBroadcast  chan []byte
	PlayerRegister   chan *Player
	PlayerUnregister chan *Player
	Players          map[*Player]bool
}
type RoomMessage struct {
	MessageType        MessageType         `json:"message_type"`
	PlayerPostionInfos []PlayerPostionInfo `json:"player_position_info"`
	PlayerCount        int                 `json:"player_count"`
}
type PlayerPostionInfo struct {
	DeviceID         string         `json:"device_id"`
	HeadPotion       model.Vector3f `json:"head_position"`
	HeadForward      model.Vector3f `json:"head_forward,omitempty"`
	LeftHandPostion  model.Vector3f `json:"left_hand_position"`
	LeftHandForward  model.Vector3f `json:"left_hand_forward,omitempty"`
	RightHandPostion model.Vector3f `json:"right_hand_position"`
	RightHandForward model.Vector3f `json:"right_hand_forward,omitempty"`
	LeftHandAvail    bool           `json:"left_hand_available"`
	RightHandAvail   bool           `json:"right_hand_available"`
}

func NewRoom(roomID string) *Room {
	room := &Room{
		RoomID:           roomID,
		PlayerBroadcast:  make(chan []byte, 1024),
		PlayerRegister:   make(chan *Player),
		PlayerUnregister: make(chan *Player),
		Players:          make(map[*Player]bool),
	}
	return room
}
func (r *Room) Run() {
	updater := false
	updaterChannel := make(chan struct{})
	defer close(updaterChannel)
	for {
		select {
		case player := <-r.PlayerRegister:
			if !updater {
				updater = true
				go r.UpdateInfo(updaterChannel)
				log.Println("Updater Started")
			}
			r.Players[player] = true
		case player := <-r.PlayerUnregister:
			if _, ok := r.Players[player]; ok {
				delete(r.Players, player)
				close(player.InChannel)
				if len(r.Players) == 0 {
					updater = false
					updaterChannel <- struct{}{}
					log.Println("Updater Stopped")
				}
			}
		case message := <-r.PlayerBroadcast:
			//Handle Messages from Players
			log.Println("Message Received: ", string(message))
			var playerMessage model.PlayerMessage
			err := json.Unmarshal(message, &playerMessage)
			if err != nil {
				log.Println("Error Unmarshalling Player Message: ", err)
				continue
			}

		}
	}
}
func (r *Room) UpdateInfo(stop chan struct{}) {
	//Routine Update to all players runing at 30 tps
	ticker := time.NewTicker(time.Second / TickRate)
	defer ticker.Stop()
	for range ticker.C {
		select {
		case <-stop:
			return
		default:
			if len(r.Players) == 0 {
				continue
			}
			//Send Player Position Info to all players
			playerPostionInfos := make([]PlayerPostionInfo, 0, len(r.Players))
			for player := range r.Players {
				playerPostionInfos = append(playerPostionInfos, PlayerPostionInfo{
					DeviceID:         player.DeiviceID,
					HeadPotion:       player.HeadPotion,
					HeadForward:      player.HeadForward,
					LeftHandPostion:  player.LeftHandPostion,
					LeftHandForward:  player.LeftHandForward,
					RightHandPostion: player.RightHandPostion,
					RightHandForward: player.RightHandForward,
					LeftHandAvail:    player.LeftHandAvail,
					RightHandAvail:   player.RightHandAvail,
				})
			}
			roomMessage := RoomMessage{
				MessageType:        MessageTypeUpdate,
				PlayerPostionInfos: playerPostionInfos,
				PlayerCount:        len(r.Players),
			}
			messageBytes, err := json.Marshal(roomMessage)
			if err != nil {
				log.Println("Error Marshalling Room Message: ", err)
				continue
			}
			for player := range r.Players {
				player.InChannel <- messageBytes
			}

		}
	}
}
func ConnectToRoom() {

}
