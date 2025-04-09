package model

type Heartbeat struct {
	Timestamp int64  `json:"timestamp"`
	DeviceID  string `json:"device_id"`
	Stage     int    `json:"stage"`
}
