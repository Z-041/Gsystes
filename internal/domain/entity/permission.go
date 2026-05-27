package entity

import "time"

type Permission struct {
	ID        uint
	Name      string
	Code      string
	Type      int
	ParentID  uint
	Path      string
	Method    string
	Sort      int
	CreatedAt time.Time
	UpdatedAt time.Time
}