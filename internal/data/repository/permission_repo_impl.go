package repository

import (
	"github.com/gsystes/backend/internal/data/model"
	domainEntity "github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) domainRepo.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) toDomain(m *model.Permission) *domainEntity.Permission {
	if m == nil {
		return nil
	}
	return &domainEntity.Permission{
		ID:        m.ID,
		Name:      m.Name,
		Code:      m.Code,
		Type:      m.Type,
		ParentID:  m.ParentID,
		Path:      m.Path,
		Method:    m.Method,
		Sort:      m.Sort,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *permissionRepository) Create(p *domainEntity.Permission) error {
	return r.db.Create(&model.Permission{
		Name:     p.Name,
		Code:     p.Code,
		Type:     p.Type,
		ParentID: p.ParentID,
		Path:     p.Path,
		Method:   p.Method,
		Sort:     p.Sort,
	}).Error
}

func (r *permissionRepository) Update(p *domainEntity.Permission) error {
	return r.db.Model(&model.Permission{}).Where("id = ?", p.ID).Updates(map[string]interface{}{
		"name":      p.Name,
		"code":      p.Code,
		"type":      p.Type,
		"parent_id": p.ParentID,
		"path":      p.Path,
		"method":    p.Method,
		"sort":      p.Sort,
	}).Error
}

func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Permission{}, id).Error
}

func (r *permissionRepository) FindByID(id uint) (*domainEntity.Permission, error) {
	var m model.Permission
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *permissionRepository) FindByCode(code string) (*domainEntity.Permission, error) {
	var m model.Permission
	if err := r.db.Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *permissionRepository) FindAll() ([]domainEntity.Permission, error) {
	var models []model.Permission
	if err := r.db.Order("sort ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]domainEntity.Permission, len(models))
	for i, m := range models {
		entities[i] = *r.toDomain(&m)
	}
	return entities, nil
}

func (r *permissionRepository) FindByPage(page, pageSize int) ([]domainEntity.Permission, int64, error) {
	var models []model.Permission
	var total int64

	if err := r.db.Model(&model.Permission{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := r.db.Offset(offset).Limit(pageSize).Order("sort ASC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]domainEntity.Permission, len(models))
	for i, m := range models {
		entities[i] = *r.toDomain(&m)
	}
	return entities, total, nil
}

func (r *permissionRepository) FindByRoleID(roleID uint) ([]domainEntity.Permission, error) {
	var models []model.Permission
	err := r.db.Raw(`
		SELECT p.* FROM sys_permissions p
		INNER JOIN sys_role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.sort ASC
	`, roleID).Scan(&models).Error
	if err != nil {
		return nil, err
	}
	entities := make([]domainEntity.Permission, len(models))
	for i, m := range models {
		entities[i] = *r.toDomain(&m)
	}
	return entities, nil
}
