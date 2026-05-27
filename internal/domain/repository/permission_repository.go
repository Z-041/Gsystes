package repository

import "github.com/gsystes/backend/internal/domain/entity"

type PermissionRepository interface {
	Create(permission *entity.Permission) error
	Update(permission *entity.Permission) error
	Delete(id uint) error
	FindByID(id uint) (*entity.Permission, error)
	FindAll() ([]entity.Permission, error)
	FindByRoleID(roleID uint) ([]entity.Permission, error)
}