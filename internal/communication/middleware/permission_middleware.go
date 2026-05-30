package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/domain/entity"
	"github.com/gsystes/backend/internal/domain/repository"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

const (
	permCacheTTL        = 10 * time.Minute
	permCacheCleanupInt = 30 * time.Minute
)

type permCacheEntry struct {
	permissions []entity.Permission
	expiresAt   time.Time
}

type PermissionMiddleware struct {
	roleRepo   repository.RoleRepository
	cache      sync.Map
	cleanupMu  sync.Mutex
	lastClean  time.Time
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
		entry := cached.(*permCacheEntry)
		if time.Now().Before(entry.expiresAt) {
			return entry.permissions
		}
	}

	permissions, err := m.roleRepo.GetPermissions(roleID)
	if err != nil {
		return nil
	}

	m.cache.Store(roleID, &permCacheEntry{
		permissions: permissions,
		expiresAt:   time.Now().Add(permCacheTTL),
	})

	m.maybeCleanup()

	return permissions
}

func (m *PermissionMiddleware) maybeCleanup() {
	m.cleanupMu.Lock()
	defer m.cleanupMu.Unlock()

	if time.Since(m.lastClean) < permCacheCleanupInt {
		return
	}
	m.lastClean = time.Now()

	now := time.Now()
	m.cache.Range(func(key, value interface{}) bool {
		entry := value.(*permCacheEntry)
		if now.After(entry.expiresAt) {
			m.cache.Delete(key)
		}
		return true
	})
}

func (m *PermissionMiddleware) InvalidateCache(roleID uint) {
	m.cache.Delete(roleID)
}
