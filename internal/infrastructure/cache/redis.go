package cache

import (
	"context"
	"net"
	"time"

	"github.com/gsystes/backend/internal/infrastructure/config"
	"github.com/gsystes/backend/internal/infrastructure/logger"
	"github.com/redis/go-redis/v9"
)

var globalRedis *redis.Client

func InitRedis(cfg config.RedisConfig) error {
	addr := cfg.Addr()

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		conn = nil
		logger.Warn("redis unavailable, cache disabled", logger.StringField("addr", addr), logger.ErrorField(err))
		return nil
	}
	conn.Close()

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		DialTimeout:  3 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		logger.Warn("redis ping failed, cache disabled", logger.ErrorField(err))
		return nil
	}

	globalRedis = client
	logger.Info("redis connected successfully")
	return nil
}

func GetRedis() *redis.Client {
	return globalRedis
}

func Close() error {
	if globalRedis != nil {
		return globalRedis.Close()
	}
	return nil
}
