package cache

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
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
	databaseStr := os.Getenv("REDIS_DB")
	database, err := strconv.Atoi(databaseStr)
	if err != nil {
		log.Fatalf("Invalid REDIS_DB value: %v", err)
	}

	addr := host + ":" + port

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	})
}
