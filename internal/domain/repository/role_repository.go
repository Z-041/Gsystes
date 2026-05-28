package repository

import "github.com/gsystes/backend/internal/domain/entity"

type RoleRepository interface {
	Create(role *entity.Role) error
	Update(role *entity.Role) error
	Delete(id uint) error
	FindByID(id uint) (*entity.Role, error)
	FindByCode(code string) (*entity.Role, error)
	FindAll() ([]entity.Role, error)
	FindByPage(page, pageSize int) ([]entity.Role, int64, error)
	AssignPermissions(roleID uint, permissionIDs []uint) error
	GetPermissions(roleID uint) ([]entity.Permission, error)
}
