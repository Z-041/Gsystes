package repository

import "github.com/gsystes/backend/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	Update(user *entity.User) error
	Delete(id uint) error
	FindByID(id uint) (*entity.User, error)
	FindByUsername(username string) (*entity.User, error)
	FindByPage(page, pageSize int, conditions map[string]interface{}) ([]entity.User, int64, error)
	FindByRoleID(roleID uint) ([]entity.User, error)
	BatchUpdateRole(userIDs []uint, roleID uint) error
}
