package model

import "time"

type Role struct {
	ID          uint      `gorm:"primarykey"`
	Name        string    `gorm:"column:name;type:varchar(64);not null"`
	Code        string    `gorm:"column:code;type:varchar(64);uniqueIndex;not null"`
	Description string    `gorm:"column:description;type:varchar(256)"`
	Status      int       `gorm:"column:status;type:tinyint;default:1"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (Role) TableName() string {
	return "sys_roles"
}