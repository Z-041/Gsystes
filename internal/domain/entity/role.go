package entity

import "time"

type Role struct {
	ID          uint
	Name        string
	Code        string
	Description string
	Status      int
	Permissions []*Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RoleStatus int

const (
	RoleStatusActive   RoleStatus = 1
	RoleStatusInactive RoleStatus = 2
)
