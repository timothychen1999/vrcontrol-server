package model

type EventType string

const (
	EventMoveCommand       EventType = "move_command"
	EventTypeShotEvent     EventType = "shot_event"
	EventTypeLatern        EventType = "latern"
	EventTypeQA            EventType = "qa"
	EventTypeAsignSequence EventType = "assign_sequence"
)

type EventMessage struct {
	EventType   EventType           `json:"event_type"`
	MoveCommand *MoveCommandMessage `json:"move_command,omitempty"`
	ShotEvent   *ShotEventMessage   `json:"shot_event,omitempty"`
	Latern      *LaternEventMessage `json:"latern,omitempty"`
	QA          *QAEventMessage     `json:"qa,omitempty"`
	Sequence    *int                `json:"sequence,omitempty"`
}
type MoveCommandMessage struct {
	Force            bool `json:"force"`
	DestinationStage int  `json:"chapter"`
}
type LaternEventMessage struct {
	LaternID int        `json:"latern_id"`
	Postions []Vector3f `json:"postions"`
}
type ShotEventMessage struct {
	Position  Vector3f `json:"position"`
	Direction Vector3f `json:"direction"`
}
type QAEventMessage struct {
	QuestionID int  `json:"question_id"`
	StateID    int  `json:"state_id"`
	State      bool `json:"state"`
}
