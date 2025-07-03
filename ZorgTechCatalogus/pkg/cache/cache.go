package cache

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Keys(context.Context, string) *redis.StringSliceCmd
	Del(context.Context, ...string) *redis.IntCmd
}

func NewRedisClient() *redis.Client {
	_ = godotenv.Load()

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	database := os.Getenv("REDIS_DB")

	addr := host + ":" + port

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	})
}
