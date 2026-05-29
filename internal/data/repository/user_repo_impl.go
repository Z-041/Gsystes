package repository

import (
	"github.com/gsystes/backend/internal/data/model"
	domainEntity "github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domainRepo.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) toDomain(m *model.User) *domainEntity.User {
	if m == nil {
		return nil
	}
	user := &domainEntity.User{
		ID:        m.ID,
		Username:  m.Username,
		Password:  m.Password,
		Nickname:  m.Nickname,
		Email:     m.Email,
		Phone:     m.Phone,
		Avatar:    m.Avatar,
		Status:    m.Status,
		RoleID:    m.RoleID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Role.ID != 0 {
		user.Role = &domainEntity.Role{
			ID:          m.Role.ID,
			Name:        m.Role.Name,
			Code:        m.Role.Code,
			Description: m.Role.Description,
			Status:      m.Role.Status,
			CreatedAt:   m.Role.CreatedAt,
			UpdatedAt:   m.Role.UpdatedAt,
		}
	}
	return user
}

func (r *userRepository) toModel(d *domainEntity.User) *model.User {
	return &model.User{
		ID:        d.ID,
		Username:  d.Username,
		Password:  d.Password,
		Nickname:  d.Nickname,
		Email:     d.Email,
		Phone:     d.Phone,
		Avatar:    d.Avatar,
		Status:    d.Status,
		RoleID:    d.RoleID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (r *userRepository) Create(user *domainEntity.User) error {
	m := r.toModel(user)
	return r.db.Create(m).Error
}

func (r *userRepository) Update(user *domainEntity.User) error {
	return r.db.Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"nickname": user.Nickname,
		"email":    user.Email,
		"phone":    user.Phone,
		"avatar":   user.Avatar,
		"status":   user.Status,
		"role_id":  user.RoleID,
	}).Error
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepository) FindByID(id uint) (*domainEntity.User, error) {
	var m model.User
	if err := r.db.Preload("Role").First(&m, id).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *userRepository) FindByUsername(username string) (*domainEntity.User, error) {
	var m model.User
	if err := r.db.Where("username = ?", username).Preload("Role").First(&m).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *userRepository) FindByPage(page, pageSize int, conditions map[string]interface{}) ([]domainEntity.User, int64, error) {
	var models []model.User
	var total int64

	query := r.db.Model(&model.User{})
	for key, value := range conditions {
		query = query.Where(key, value)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Role").Offset(offset).Limit(pageSize).Order("id DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]domainEntity.User, len(models))
	for i, m := range models {
		entities[i] = *r.toDomain(&m)
	}
	return entities, total, nil
}

func (r *userRepository) FindByRoleID(roleID uint) ([]domainEntity.User, error) {
	var models []model.User
	if err := r.db.Where("role_id = ?", roleID).Preload("Role").Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]domainEntity.User, len(models))
	for i, m := range models {
		entities[i] = *r.toDomain(&m)
	}
	return entities, nil
}

func (r *userRepository) BatchUpdateRole(userIDs []uint, roleID uint) error {
	return r.db.Model(&model.User{}).Where("id IN ?", userIDs).Update("role_id", roleID).Error
}

func (r *userRepository) BatchCreate(users []*domainEntity.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, u := range users {
			m := r.toModel(u)
			if err := tx.Create(m).Error; err != nil {
				return err
			}
			u.ID = m.ID
		}
		return nil
	})
}

func (r *userRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.User{}).Count(&count).Error
	return count, err
}
