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

type Movement struct {
	Force            bool
	DestinationStage int
	Target           string
	Broadcast        bool
}
type Room struct {
	RoomID           string
	PlayerBroadcast  chan []byte
	PlayerRegister   chan *Player
	PlayerUnregister chan *Player
	MoveControl      chan Movement
	Signals          chan ControlSignal
	Players          map[*Player]bool
	AssignedSequence map[string]int
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

type ControlSignal struct {
	Target *Player
	Type   string
	Args   []string
}

func NewRoom(roomID string) *Room {
	room := &Room{
		RoomID:           roomID,
		PlayerBroadcast:  make(chan []byte, 1024),
		PlayerRegister:   make(chan *Player),
		PlayerUnregister: make(chan *Player),
		Players:          make(map[*Player]bool),
		MoveControl:      make(chan Movement),
		Signals:          make(chan ControlSignal),
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
			log.Println("Player Registered: ", player.DeiviceID)
			update := r.PlayerSequenceUpdate()
			// Send the player sequence update to the newly registered player
			for _, seqUpdate := range update {
				if seqUpdate.Player == nil {
					continue
				}
				eventMessage := model.EventMessage{
					EventType: model.EventTypeAsignSequence,
					Sequence:  &seqUpdate.Sequence,
				}
				message, err := json.Marshal(eventMessage)
				if err != nil {
					log.Println("Error Marshalling Event Message: ", err)
					continue
				}
				select {
				case seqUpdate.Player.InChannel <- message:
				default:
					log.Println("Player Channel is full, disconnecting player")
					r.PlayerUnregister <- seqUpdate.Player
				}
			}
		case player := <-r.PlayerUnregister:
			if _, ok := r.Players[player]; ok {
				delete(r.Players, player)
				log.Println("Player Unregistered: ", player.DeiviceID)
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
			// Forward the message to all players except the sender

		case move := <-r.MoveControl:
			if move.Broadcast {
				for player := range r.Players {
					if player == nil {
						continue
					} else {
						eventMessage := model.EventMessage{
							EventType: model.EventMoveCommand,
							MoveCommand: &model.MoveCommandMessage{
								Force:            move.Force,
								DestinationStage: move.DestinationStage,
							},
						}
						message, err := json.Marshal(eventMessage)
						if err != nil {
							log.Println("Error Marshalling Event Message: ", err)
							continue
						}
						// Send the message to all players
						select {
						case player.InChannel <- message:
						default:
							log.Println("Player Channel is full, disconnecting player")
							r.PlayerUnregister <- player
						}
					}
				}
			} else {
				for player := range r.Players {
					if player == nil {
						continue
					} else if player.DeiviceID == move.Target {
						eventMessage := model.EventMessage{
							EventType: model.EventMoveCommand,
							MoveCommand: &model.MoveCommandMessage{
								Force:            move.Force,
								DestinationStage: move.DestinationStage,
							},
						}
						message, err := json.Marshal(eventMessage)
						if err != nil {
							log.Println("Error Marshalling Event Message: ", err)
							continue
						}
						select {
						case player.InChannel <- message:
						default:
							log.Println("Player Channel is full, disconnecting player")
							r.PlayerUnregister <- player
						}
					}
				}
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
				select {
				case player.InChannel <- messageBytes:
				default:
					log.Println("Player Channel is full, disconnecting player")
					r.PlayerUnregister <- player
				}
			}

		}
	}
}
func (r *Room) GetPlayerByDeviceID(deviceID string) *Player {
	for player := range r.Players {
		if player.DeiviceID == deviceID {
			return player
		}
	}
	return nil
}
