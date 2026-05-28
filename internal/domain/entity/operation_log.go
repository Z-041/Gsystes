package entity

import "time"

type OperationLog struct {
	ID         uint
	UserID     uint
	Username   string
	Method     string
	Path       string
	Query      string
	Body       string
	StatusCode int
	Latency    int64
	ClientIP   string
	UserAgent  string
	CreatedAt  time.Time
}
