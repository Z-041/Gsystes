package seed

import (
	"fmt"

	"github.com/gsystes/backend/internal/data/model"
	"github.com/gsystes/backend/internal/domain/entity"
	domainRepo "github.com/gsystes/backend/internal/domain/repository"
	domainService "github.com/gsystes/backend/internal/domain/service"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	"gorm.io/gorm"
)

var defaultRoles = []*entity.Role{
	{Name: "普通用户", Code: "user", Description: "基础用户角色，拥有基本查看权限", Status: 1},
	{Name: "编辑者", Code: "editor", Description: "编辑者角色，拥有内容管理权限", Status: 1},
	{Name: "审计员", Code: "auditor", Description: "审计员角色，拥有日志查看和审计权限", Status: 1},
}

var defaultPermissions = []*entity.Permission{
	{Name: "用户管理", Code: "user:manage", Type: 1, Path: "/users", Sort: 1},
	{Name: "创建用户", Code: "user:create", Type: 2, Path: "/users", Method: "POST", Sort: 2},
	{Name: "编辑用户", Code: "user:update", Type: 2, Path: "/users", Method: "PUT", Sort: 3},
	{Name: "删除用户", Code: "user:delete", Type: 2, Path: "/users", Method: "DELETE", Sort: 4},
	{Name: "查询用户", Code: "user:read", Type: 2, Path: "/users", Method: "GET", Sort: 5},
	{Name: "角色管理", Code: "role:manage", Type: 1, Path: "/roles", Sort: 6},
	{Name: "创建角色", Code: "role:create", Type: 2, Path: "/roles", Method: "POST", Sort: 7},
	{Name: "编辑角色", Code: "role:update", Type: 2, Path: "/roles", Method: "PUT", Sort: 8},
	{Name: "删除角色", Code: "role:delete", Type: 2, Path: "/roles", Method: "DELETE", Sort: 9},
	{Name: "查询角色", Code: "role:read", Type: 2, Path: "/roles", Method: "GET", Sort: 10},
	{Name: "权限管理", Code: "perm:manage", Type: 1, Path: "/permissions", Sort: 11},
	{Name: "分配权限", Code: "perm:assign", Type: 2, Path: "/permissions", Method: "POST", Sort: 12},
	{Name: "操作日志", Code: "log:manage", Type: 1, Path: "/logs", Sort: 13},
	{Name: "查询日志", Code: "log:read", Type: 2, Path: "/logs", Method: "GET", Sort: 14},
}

func ensureRoles(db *gorm.DB, roleRepo domainRepo.RoleRepository) error {
	for _, r := range defaultRoles {
		var count int64
		db.Model(&model.Role{}).Where("code = ?", r.Code).Count(&count)
		if count > 0 {
			continue
		}
		if err := roleRepo.Create(r); err != nil {
			return fmt.Errorf("failed to create role %s: %w", r.Code, err)
		}
		logger.Info("default role created", logger.StringField("code", r.Code), logger.StringField("name", r.Name))
	}
	return nil
}

func ensurePermissions(db *gorm.DB, permRepo domainRepo.PermissionRepository) error {
	for _, p := range defaultPermissions {
		var count int64
		db.Model(&model.Permission{}).Where("code = ?", p.Code).Count(&count)
		if count > 0 {
			continue
		}
		if err := permRepo.Create(p); err != nil {
			return fmt.Errorf("failed to create permission %s: %w", p.Code, err)
		}
		logger.Info("permission added", logger.StringField("code", p.Code))
	}
	return nil
}

func ensureSuperAdminRole(db *gorm.DB, roleRepo domainRepo.RoleRepository) (*model.Role, error) {
	var roleModel model.Role
	err := db.Where("code = ?", "super_admin").First(&roleModel).Error
	if err == nil {
		return &roleModel, nil
	}

	adminRole := &entity.Role{
		Name:        "超级管理员",
		Code:        "super_admin",
		Description: "系统超级管理员，拥有所有权限",
		Status:      1,
	}
	if err := roleRepo.Create(adminRole); err != nil {
		return nil, fmt.Errorf("failed to create admin role: %w", err)
	}
	logger.Info("admin role created", logger.StringField("code", "super_admin"))

	if err := db.Where("code = ?", "super_admin").First(&roleModel).Error; err != nil {
		return nil, fmt.Errorf("failed to get admin role: %w", err)
	}
	return &roleModel, nil
}

func assignAllPermissionsToRole(db *gorm.DB, roleID uint, roleRepo domainRepo.RoleRepository) error {
	var permModels []model.Permission
	if err := db.Order("sort ASC").Find(&permModels).Error; err != nil {
		return fmt.Errorf("failed to query all permissions: %w", err)
	}
	permIDs := make([]uint, len(permModels))
	for i, m := range permModels {
		permIDs[i] = m.ID
	}
	if err := roleRepo.AssignPermissions(roleID, permIDs); err != nil {
		return fmt.Errorf("failed to assign permissions to role %d: %w", roleID, err)
	}
	logger.Info("all permissions assigned to role", logger.UintField("role_id", roleID), logger.IntField("count", len(permIDs)))
	return nil
}

func InitSeedData(
	db *gorm.DB,
	userSvc *domainService.UserDomainService,
	userRepo domainRepo.UserRepository,
	roleRepo domainRepo.RoleRepository,
	permRepo domainRepo.PermissionRepository,
) error {
	// 1. Ensure super_admin role exists
	roleModel, err := ensureSuperAdminRole(db, roleRepo)
	if err != nil {
		return err
	}

	// 2. Ensure default roles exist (普通用户, 编辑者, 审计员)
	if err := ensureRoles(db, roleRepo); err != nil {
		return err
	}

	// 3. Ensure all default permissions exist (idempotent)
	if err := ensurePermissions(db, permRepo); err != nil {
		return err
	}

	// 4. Assign all permissions to super_admin role (idempotent)
	if err := assignAllPermissionsToRole(db, roleModel.ID, roleRepo); err != nil {
		return err
	}

	// 5. Create admin user if not exists
	existing, _ := userRepo.FindByUsername("admin")
	if existing != nil {
		logger.Info("seed data check completed, admin user already exists")
		return nil
	}

	adminUser := &entity.User{
		Username: "admin",
		Nickname: "超级管理员",
		Email:    "admin@gsystes.com",
		RoleID:   roleModel.ID,
	}
	if err := userSvc.Create(adminUser, "admin123"); err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}
	logger.Info("admin user created", logger.StringField("username", "admin"))

	logger.Info("seed data initialization completed")
	return nil
}
