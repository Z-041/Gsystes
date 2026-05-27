package repository

import (
	"github.com/gsystes/backend/internal/data/model"
	domainEntity "github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domainRepo.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) toDomain(m *model.Role) *domainEntity.Role {
	if m == nil {
		return nil
	}
	return &domainEntity.Role{
		ID:          m.ID,
		Name:        m.Name,
		Code:        m.Code,
		Description: m.Description,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *roleRepository) Create(role *domainEntity.Role) error {
	return r.db.Create(&model.Role{
		Name:        role.Name,
		Code:        role.Code,
		Description: role.Description,
		Status:      role.Status,
	}).Error
}

func (r *roleRepository) Update(role *domainEntity.Role) error {
	return r.db.Model(&model.Role{}).Where("id = ?", role.ID).Updates(map[string]interface{}{
		"name":        role.Name,
		"code":        role.Code,
		"description": role.Description,
		"status":      role.Status,
	}).Error
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}

func (r *roleRepository) FindByID(id uint) (*domainEntity.Role, error) {
	var m model.Role
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *roleRepository) FindAll() ([]domainEntity.Role, error) {
	var models []model.Role
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]domainEntity.Role, len(models))
	for i, m := range models {
		entities[i] = *r.toDomain(&m)
	}
	return entities, nil
}

func (r *roleRepository) FindByPage(page, pageSize int) ([]domainEntity.Role, int64, error) {
	var models []model.Role
	var total int64

	if err := r.db.Model(&model.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]domainEntity.Role, len(models))
	for i, m := range models {
		entities[i] = *r.toDomain(&m)
	}
	return entities, total, nil
}