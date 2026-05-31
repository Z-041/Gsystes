package middleware

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/domain/entity"
	"github.com/gsystes/backend/internal/domain/repository"
	"github.com/gsystes/backend/internal/infrastructure/cache"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
	"github.com/gsystes/backend/internal/infrastructure/utils"
)

const (
	permCacheTTL            = 10 * time.Minute
	permCacheCleanupInt     = 30 * time.Minute
	permCacheInvalidateChan = "perm_cache:invalidate"
)

type permCacheEntry struct {
	permissions []entity.Permission
	expiresAt   time.Time
}

type permInvalidateMsg struct {
	RoleID uint `json:"role_id"`
}

type PermissionMiddleware struct {
	roleRepo  repository.RoleRepository
	cache     sync.Map
	cleanupMu sync.Mutex
	lastClean time.Time
	stopCh    chan struct{}
}

func NewPermissionMiddleware(roleRepo repository.RoleRepository) *PermissionMiddleware {
	m := &PermissionMiddleware{
		roleRepo: roleRepo,
		stopCh:   make(chan struct{}),
	}
	m.startRedisListener()
	return m
}

func (m *PermissionMiddleware) startRedisListener() {
	rdb := cache.GetRedis()
	if rdb == nil {
		return
	}

	go func() {
		pubsub := rdb.Subscribe(context.Background(), permCacheInvalidateChan)
		defer pubsub.Close()

		ch := pubsub.Channel()
		for {
			select {
			case msg := <-ch:
				var inv permInvalidateMsg
				if err := json.Unmarshal([]byte(msg.Payload), &inv); err != nil {
					continue
				}
				m.cache.Delete(inv.RoleID)
			case <-m.stopCh:
				return
			}
		}
	}()
}

func (m *PermissionMiddleware) Stop() {
	close(m.stopCh)
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
	m.publishInvalidate(roleID)
}

func (m *PermissionMiddleware) publishInvalidate(roleID uint) {
	rdb := cache.GetRedis()
	if rdb == nil {
		return
	}

	data, _ := json.Marshal(permInvalidateMsg{RoleID: roleID})
	err := rdb.Publish(context.Background(), permCacheInvalidateChan, data).Err()
	if err != nil {
		logger.Warn("failed to publish perm cache invalidate",
			logger.UintField("role_id", roleID),
			logger.ErrorField(err),
		)
	}
}

func parseRoleID(key interface{}) uint {
	switch v := key.(type) {
	case uint:
		return v
	case int:
		return uint(v)
	case string:
		id, _ := strconv.ParseUint(v, 10, 64)
		return uint(id)
	default:
		return 0
	}
}
