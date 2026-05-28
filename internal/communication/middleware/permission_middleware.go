package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/domain/repository"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

type PermissionMiddleware struct {
	roleRepo repository.RoleRepository
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

		permissions, err := m.roleRepo.GetPermissions(roleID)
		if err != nil {
			utils.Forbidden(c, "failed to check permissions")
			c.Abort()
			return
		}

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
