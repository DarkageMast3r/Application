package cache

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

// ErrNotFound is een generieke fout voor wanneer een sleutel niet in de cache wordt gevonden.
var ErrNotFound = errors.New("key not found in cache")

// Cache defines the consistent interface for cache operations (e.g., Redis, Memcached)
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error                          // Variadic voor meerdere sleutels
	Exists(ctx context.Context, key string) (bool, error)                   // Nieuw: expliciete check
	Increment(ctx context.Context, key string) (int64, error)               // Nieuw: voor rate limiting
	Expire(ctx context.Context, key string, expiration time.Duration) error // Nieuw: om TTL in te stellen
	// Keys(context.Context, string) ([]string, error) // Optioneel: als je Keys echt nodig hebt, maar vaak vermeden voor grote caches
}

// RedisClient implements the Cache interface using go-redis client
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient initializes a new Redis client and returns the Cache interface.
func NewRedisClient() (Cache, error) { // Return Cache interface, not *redis.Client
	_ = godotenv.Load()

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")
	password := os.Getenv("REDIS_PASSWORD")
	databaseStr := os.Getenv("REDIS_DB")
	database, err := strconv.Atoi(databaseStr)
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB value: %w", err)
	}

	addr := host + ":" + port

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping om verbinding te testen
	_, err = client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connection successful.")
	return &RedisClient{client: client}, nil // Return de implementatie van de interface
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", ErrNotFound // Gebruik je eigen ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("redis GET failed for key %s: %w", key, err)
	}
	return val, nil
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis SET failed for key %s: %w", key, err)
	}
	return nil
}

func (r *RedisClient) Del(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("redis DEL failed for keys %v: %w", keys, err)
	}
	return nil
}

func (r *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("redis EXISTS failed for key %s: %w", key, err)
	}
	return count > 0, nil
}

func (r *RedisClient) Increment(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis INCR failed for key %s: %w", key, err)
	}
	return val, nil
}

func (r *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := r.client.Expire(ctx, key, expiration).Err()
	if err != nil {
		return fmt.Errorf("redis EXPIRE failed for key %s: %w", key, err)
	}
	return nil
}

// Keys is over het algemeen niet aan te raden in productieomgevingen vanwege prestaties
// func (r *RedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
// 	return r.client.Keys(ctx, pattern).Result()
// }
