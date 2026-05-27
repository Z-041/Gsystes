package entity

import "time"

type User struct {
	ID        uint
	Username  string
	Password  string
	Nickname  string
	Email     string
	Phone     string
	Avatar    string
	Status    int
	RoleID    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserStatus int

const (
	UserStatusActive   UserStatus = 1
	UserStatusInactive UserStatus = 2
)
