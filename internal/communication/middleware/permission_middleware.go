package middleware

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/domain/entity"
	"github.com/gsystes/backend/internal/domain/repository"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

type PermissionMiddleware struct {
	roleRepo repository.RoleRepository
	cache    sync.Map
}

func NewPermissionMiddleware(roleRepo repository.RoleRepository) *PermissionMiddleware {
	return &PermissionMiddleware{roleRepo: roleRepo}
}

func (m *PermissionMiddleware) Require(permCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := infraMiddleware.GetClaims(c)
		if claims == nil {
			utils.Forbidden(c, "no permission")
			c.Abort()
			return
		}

		roleID := claims.RoleID
		if roleID == 0 {
			utils.Forbidden(c, "no permission")
			c.Abort()
			return
		}

		permissions := m.getCachedPermissions(roleID)
		for _, p := range permissions {
			if p.Code == permCode {
				c.Next()
				return
			}
		}

		utils.Forbidden(c, "no permission: "+permCode)
		c.Abort()
	}
}

func (m *PermissionMiddleware) getCachedPermissions(roleID uint) []entity.Permission {
	if cached, ok := m.cache.Load(roleID); ok {
		return cached.([]entity.Permission)
	}

	permissions, err := m.roleRepo.GetPermissions(roleID)
	if err != nil {
		return nil
	}

	m.cache.Store(roleID, permissions)
	return permissions
}

func (m *PermissionMiddleware) InvalidateCache(roleID uint) {
	m.cache.Delete(roleID)
}