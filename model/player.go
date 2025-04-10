package model

type MessageType string

const (
	MessageTypeHeartbeat   MessageType = "heartbeat"
	MessageTypeReadyToMove MessageType = "ready_to_move"
	MessageTypeShotEvent   MessageType = "shot_event"
	MessageTypeLatern      MessageType = "latern"
)

type Vector3f struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	Z float32 `json:"z"`
}

type PlayerMessage struct {
	MessageType string       `json:"message_type"`
	Heartbeat   *Heartbeat   `json:"heartbeat,omitempty"`
	ShotEvent   *ShotEvent   `json:"shot_event,omitempty"`
	Latern      *Latern      `json:"latern,omitempty"`
	ReadyToMove *ReadyToMove `json:"ready_to_move,omitempty"`
}

type Heartbeat struct {
	Timestamp        int64    `json:"timestamp"`
	DeviceID         string   `json:"device_id"`
	Stage            int      `json:"chapter"`
	Message          string   `json:"message"`
	HeadPotion       Vector3f `json:"head_position"`
	HeadForward      Vector3f `json:"head_forward,omitempty"`
	LeftHandPostion  Vector3f `json:"left_hand_position"`
	LeftHandForward  Vector3f `json:"left_hand_forward,omitempty"`
	RightHandPostion Vector3f `json:"right_hand_position"`
	RightHandForward Vector3f `json:"right_hand_forward,omitempty"`
	LeftHandAvail    bool     `json:"left_hand_available"`
	RightHandAvail   bool     `json:"right_hand_available"`
}
type ShotEvent struct {
	Timestamp int64    `json:"timestamp"`
	DeviceID  string   `json:"device_id"`
	Position  Vector3f `json:"position"`
	Direction Vector3f `json:"direction"`
}
type Latern struct {
	Timestamp int64      `json:"timestamp"`
	DeviceID  string     `json:"device_id"`
	LaternID  int        `json:"latern_id"`
	Postions  []Vector3f `json:"postions"`
}
type ReadyToMove struct {
	Timestamp int64  `json:"timestamp"`
	DeviceID  string `json:"device_id"`
	Stage     int    `json:"chapter"`
}
