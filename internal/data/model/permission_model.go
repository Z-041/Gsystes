package model

import "time"

type Permission struct {
	ID        uint      `gorm:"primarykey"`
	Name      string    `gorm:"column:name;type:varchar(64);not null"`
	Code      string    `gorm:"column:code;type:varchar(64);uniqueIndex;not null"`
	Type      int       `gorm:"column:type;type:tinyint"`
	ParentID  uint      `gorm:"column:parent_id;default:0"`
	Path      string    `gorm:"column:path;type:varchar(256)"`
	Method    string    `gorm:"column:method;type:varchar(32)"`
	Sort      int       `gorm:"column:sort;default:0"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (Permission) TableName() string {
	return "sys_permissions"
}
