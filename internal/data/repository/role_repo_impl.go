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
	perms := make([]*domainEntity.Permission, len(m.Permissions))
	for i, p := range m.Permissions {
		perms[i] = &domainEntity.Permission{
			ID:        p.ID,
			Name:      p.Name,
			Code:      p.Code,
			Type:      p.Type,
			ParentID:  p.ParentID,
			Path:      p.Path,
			Method:    p.Method,
			Sort:      p.Sort,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}
	return &domainEntity.Role{
		ID:          m.ID,
		Name:        m.Name,
		Code:        m.Code,
		Description: m.Description,
		Status:      m.Status,
		Permissions: perms,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (r *roleRepository) toModel(d *domainEntity.Role) *model.Role {
	return &model.Role{
		ID:          d.ID,
		Name:        d.Name,
		Code:        d.Code,
		Description: d.Description,
		Status:      d.Status,
	}
}

func (r *roleRepository) Create(role *domainEntity.Role) error {
	return r.db.Create(r.toModel(role)).Error
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
	if err := r.db.Preload("Permissions").First(&m, id).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *roleRepository) FindByCode(code string) (*domainEntity.Role, error) {
	var m model.Role
	if err := r.db.Where("code = ?", code).Preload("Permissions").First(&m).Error; err != nil {
		return nil, err
	}
	return r.toDomain(&m), nil
}

func (r *roleRepository) FindAll() ([]domainEntity.Role, error) {
	var models []model.Role
	if err := r.db.Order("id ASC").Find(&models).Error; err != nil {
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

func (r *roleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
			return err
		}
		if len(permissionIDs) == 0 {
			return nil
		}
		assignments := make([]model.RolePermission, len(permissionIDs))
		for i, pid := range permissionIDs {
			assignments[i] = model.RolePermission{RoleID: roleID, PermissionID: pid}
		}
		return tx.Create(&assignments).Error
	})
}

func (r *roleRepository) GetPermissions(roleID uint) ([]domainEntity.Permission, error) {
	var perms []model.Permission
	err := r.db.Raw(`
		SELECT p.* FROM sys_permissions p
		INNER JOIN sys_role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.sort ASC
	`, roleID).Scan(&perms).Error
	if err != nil {
		return nil, err
	}
	entities := make([]domainEntity.Permission, len(perms))
	for i, p := range perms {
		entities[i] = domainEntity.Permission{
			ID:        p.ID,
			Name:      p.Name,
			Code:      p.Code,
			Type:      p.Type,
			ParentID:  p.ParentID,
			Path:      p.Path,
			Method:    p.Method,
			Sort:      p.Sort,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}
	return entities, nil
}
