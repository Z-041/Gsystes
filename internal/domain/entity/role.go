package entity

import "time"

type Role struct {
	ID          uint
	Name        string
	Code        string
	Description string
	Status      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}