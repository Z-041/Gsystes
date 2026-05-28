package model

import "time"

type User struct {
	ID        uint      `gorm:"primarykey"`
	Username  string    `gorm:"column:username;type:varchar(64);uniqueIndex;not null"`
	Password  string    `gorm:"column:password;type:varchar(256);not null"`
	Nickname  string    `gorm:"column:nickname;type:varchar(64)"`
	Email     string    `gorm:"column:email;type:varchar(128)"`
	Phone     string    `gorm:"column:phone;type:varchar(20)"`
	Avatar    string    `gorm:"column:avatar;type:varchar(256)"`
	Status    int       `gorm:"column:status;type:tinyint;default:1"`
	RoleID    uint      `gorm:"column:role_id"`
	Role      Role      `gorm:"foreignKey:RoleID"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (User) TableName() string {
	return "sys_users"
}
