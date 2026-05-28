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

func InitSeedData(
	db *gorm.DB,
	userSvc *domainService.UserDomainService,
	userRepo domainRepo.UserRepository,
	roleRepo domainRepo.RoleRepository,
	permRepo domainRepo.PermissionRepository,
) error {
	existing, _ := userRepo.FindByUsername("admin")
	if existing != nil {
		logger.Info("seed data already exists, skipping")
		return nil
	}

	logger.Info("initializing seed data...")

	// 1. Create admin role
	adminRole := &entity.Role{
		Name:        "超级管理员",
		Code:        "super_admin",
		Description: "系统超级管理员，拥有所有权限",
		Status:      1,
	}
	if err := roleRepo.Create(adminRole); err != nil {
		return fmt.Errorf("failed to create admin role: %w", err)
	}

	var roleModel model.Role
	if err := db.Where("code = ?", "super_admin").First(&roleModel).Error; err != nil {
		return fmt.Errorf("failed to get admin role: %w", err)
	}
	logger.Info("admin role created", logger.StringField("code", "super_admin"))

	// 2. Create default permissions
	permissions := []*entity.Permission{
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
	permIDs := make([]uint, 0, len(permissions))
	for _, p := range permissions {
		if err := permRepo.Create(p); err != nil {
			return fmt.Errorf("failed to create permission %s: %w", p.Code, err)
		}
		var permModel model.Permission
		if err := db.Where("code = ?", p.Code).First(&permModel).Error; err != nil {
			return fmt.Errorf("failed to get permission %s: %w", p.Code, err)
		}
		permIDs = append(permIDs, permModel.ID)
	}
	logger.Info("default permissions created", logger.IntField("count", len(permissions)))

	// 3. Assign all permissions to admin role
	if err := roleRepo.AssignPermissions(roleModel.ID, permIDs); err != nil {
		return fmt.Errorf("failed to assign permissions to admin role: %w", err)
	}
	logger.Info("permissions assigned to admin role", logger.IntField("count", len(permIDs)))

	// 4. Create admin user
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
