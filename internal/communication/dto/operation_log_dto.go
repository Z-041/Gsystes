package dto

import "time"

type OperationLogResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Username    string    `json:"username"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	IP          string    `json:"ip"`
	Duration    int64     `json:"duration"`
	RequestBody string    `json:"request_body"`
	StatusCode  int       `json:"status_code"`
	CreatedAt   time.Time `json:"created_at"`
}