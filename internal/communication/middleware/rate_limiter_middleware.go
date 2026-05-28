package middleware

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gsystes/backend/internal/infrastructure/cache"
	"github.com/gsystes/backend/internal/infrastructure/config"
	infraMiddleware "github.com/gsystes/backend/internal/infrastructure/middleware"
	"github.com/gsystes/backend/internal/infrastructure/utils"
	"github.com/redis/go-redis/v9"
)

const (
	rateLimitKeyPrefix = "rate_limit:"
	cleanupInterval    = 2 * time.Minute
)

type rateLimiter struct {
	visitors map[string]*windowEntry
	mu       sync.RWMutex
	once     sync.Once
	stopCh   chan struct{}
}

type windowEntry struct {
	timestamps []time.Time
	lastAccess time.Time
}

var memoryLimiter = &rateLimiter{
	visitors: make(map[string]*windowEntry),
	stopCh:   make(chan struct{}),
}

func StopMemoryLimiter() {
	close(memoryLimiter.stopCh)
}

func (rl *rateLimiter) startCleanup() {
	rl.once.Do(func() {
		go func() {
			ticker := time.NewTicker(cleanupInterval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					rl.cleanup()
				case <-rl.stopCh:
					return
				}
			}
		}()
	})
}

func (rl *rateLimiter) allow(key string, rate int, window time.Duration) bool {
	rl.startCleanup()

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	entry, exists := rl.visitors[key]
	if !exists {
		entry = &windowEntry{}
		rl.visitors[key] = entry
	}
	entry.lastAccess = now

	cutoff := now.Add(-window)
	filtered := make([]time.Time, 0, len(entry.timestamps))
	for _, ts := range entry.timestamps {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}
	entry.timestamps = filtered

	if len(filtered) >= rate {
		return false
	}

	entry.timestamps = append(entry.timestamps, now)
	return true
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for key, entry := range rl.visitors {
		if now.Sub(entry.lastAccess) > cleanupInterval {
			delete(rl.visitors, key)
		}
	}
}

func allowWithRedis(key string, rate int, window time.Duration) bool {
	rdb := cache.GetRedis()
	if rdb == nil {
		return memoryLimiter.allow(key, rate, window)
	}

	ctx := context.Background()
	now := time.Now().UnixMicro()
	windowMicros := window.Microseconds()
	cutoff := now - windowMicros
	pipe := rdb.Pipeline()

	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(cutoff, 10))
	countCmd := pipe.ZCard(ctx, key)
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: now})
	pipe.Expire(ctx, key, window*2)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return memoryLimiter.allow(key, rate, window)
	}

	return countCmd.Val() < int64(rate)
}

func buildRateKey(c *gin.Context, prefix string) string {
	userID := infraMiddleware.GetUserID(c)
	if userID > 0 {
		return rateLimitKeyPrefix + prefix + ":user:" + c.ClientIP()
	}
	return rateLimitKeyPrefix + prefix + ":ip:" + c.ClientIP()
}

func RateLimiter(rate int, burst int, window time.Duration) gin.HandlerFunc {
	rlCfg := config.GetConfig().RateLimit
	if !rlCfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		key := buildRateKey(c, "api")
		if allowWithRedis(key, burst, window) {
			c.Next()
			return
		}
		utils.Error(c, http.StatusTooManyRequests, "too many requests, please try again later")
		c.Abort()
	}
}
