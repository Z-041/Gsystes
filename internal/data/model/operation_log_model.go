package model

import "time"

type OperationLog struct {
	ID         uint      `gorm:"primarykey"`
	UserID     uint      `gorm:"column:user_id;index"`
	Username   string    `gorm:"column:username;type:varchar(64)"`
	Method     string    `gorm:"column:method;type:varchar(10)"`
	Path       string    `gorm:"column:path;type:varchar(256)"`
	Query      string    `gorm:"column:query;type:varchar(512)"`
	Body       string    `gorm:"column:body;type:text"`
	StatusCode int       `gorm:"column:status_code;type:int"`
	Latency    int64     `gorm:"column:latency;type:bigint"`
	ClientIP   string    `gorm:"column:client_ip;type:varchar(20)"`
	UserAgent  string    `gorm:"column:user_agent;type:varchar(256)"`
	CreatedAt  time.Time `gorm:"column:created_at"`
}

func (OperationLog) TableName() string {
	return "sys_operation_logs"
}
