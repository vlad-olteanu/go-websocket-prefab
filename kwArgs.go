package wsprefab

import "time"

type KwArgs struct {
	ReadBufferSize  int
	WriteBufferSize int
	WriteWait       time.Duration
	PongWait        time.Duration
	PingPeriod      time.Duration
	MaxMessageSize  int //bytes
}
