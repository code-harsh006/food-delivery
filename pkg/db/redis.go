package db

import (
	"context"
	"fmt"
	"food-delivery-backend/pkg/config"
	"food-delivery-backend/pkg/logger"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(cfg *config.RedisConfig) error {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// Test connection
	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Connected to Redis")
	return nil
}

func GetRedis() *redis.Client {
	return RedisClient
}
