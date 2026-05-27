package repository

import "github.com/gsystes/backend/internal/domain/entity"

type PermissionRepository interface {
	Create(permission *entity.Permission) error
	Update(permission *entity.Permission) error
	Delete(id uint) error
	FindByID(id uint) (*entity.Permission, error)
	FindByCode(code string) (*entity.Permission, error)
	FindAll() ([]entity.Permission, error)
	FindByPage(page, pageSize int) ([]entity.Permission, int64, error)
	FindByRoleID(roleID uint) ([]entity.Permission, error)
}
