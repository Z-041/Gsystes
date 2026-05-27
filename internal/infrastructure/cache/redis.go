package cache

import (
    "context"
    "time"

    "github.com/gsystes/backend/internal/infrastructure/config"
    "github.com/gsystes/backend/internal/infrastructure/logger"
    "github.com/redis/go-redis/v9"
)

var globalRedis *redis.Client

func InitRedis(cfg config.RedisConfig) error {
    client := redis.NewClient(&redis.Options{
        Addr:         cfg.Addr(),
        Password:     cfg.Password,
        DB:           cfg.DB,
        DialTimeout:  5 * time.Second,
        ReadTimeout:  3 * time.Second,
        WriteTimeout: 3 * time.Second,
    })

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.Ping(ctx).Err(); err != nil {
        logger.Warn("redis connection failed, cache disabled", logger.ErrorField(err))
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