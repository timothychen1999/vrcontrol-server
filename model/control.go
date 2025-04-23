package model

type TargetType string

const (
	TargetBroadcast TargetType = "broadcast"
	TargetDevice    TargetType = "device"
	TargetServer    TargetType = "server"
)

type Control struct {
	TargetID   string     `json:"target_id"`
	TargetType TargetType `json:"target_type"`
	Command    string     `json:"command"`
	Args       []string   `json:"args"`
	Timestamp  int64      `json:"timestamp"`
	SourceID   string     `json:"source_id"`
}
