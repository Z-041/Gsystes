package model

type RolePermission struct {
	RoleID       uint `gorm:"column:role_id;primaryKey"`
	PermissionID uint `gorm:"column:permission_id;primaryKey"`
}

func (RolePermission) TableName() string {
	return "sys_role_permissions"
}
