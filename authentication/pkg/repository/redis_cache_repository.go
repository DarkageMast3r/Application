package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"authentication/pkg/cache" // Zorg dat dit pad klopt (je Redis client interface)
)

type RedisCacheRepository struct {
	client cache.Cache
}

func NewRedisCacheRepository(client cache.Cache) CacheRepository {
	return &RedisCacheRepository{client: client}
}

func (r *RedisCacheRepository) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) { // Aangenomen dat je cache.Cache zo'n fout retourneert
			return "", nil // Niet gevonden is geen fout
		}
		return "", fmt.Errorf("failed to get from cache: %w", err)
	}
	return val, nil
}

func (r *RedisCacheRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration)
	if err != nil {
		return fmt.Errorf("failed to set in cache: %w", err)
	}
	return nil
}

func (r *RedisCacheRepository) Del(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}
	return nil
}

func (r *RedisCacheRepository) Exists(ctx context.Context, key string) (bool, error) {
	// Aangenomen dat je cache.Cache een Exists methode heeft, of dat Get + nil check werkt
	val, err := r.client.Get(ctx, key)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return false, nil // Niet gevonden
		}
		return false, fmt.Errorf("failed to check existence in cache: %w", err)
	}
	return val != "", nil // Als er een waarde is, bestaat de sleutel
}

func (r *RedisCacheRepository) Increment(ctx context.Context, key string) (int64, error) {
	// Aangenomen dat je cache.Cache een Incr methode heeft
	// Als niet, moet je Get, parse, increment, Set implementeren, wat race conditions kan hebben.
	// Voor rate limiting is een atomaire Increment essentieel.
	val, err := r.client.Increment(ctx, key) // Stel dat je cache.Cache een Increment methode heeft
	if err != nil {
		return 0, fmt.Errorf("failed to increment cache key: %w", err)
	}
	return val, nil
}

func (r *RedisCacheRepository) Expire(ctx context.Context, key string, expiration time.Duration) error {
	// Aangenomen dat je cache.Cache een Expire methode heeft
	err := r.client.Expire(ctx, key, expiration)
	if err != nil {
		return fmt.Errorf("failed to set expiration for cache key: %w", err)
	}
	return nil
}

// AddBlacklistedAccessToken voegt een token toe aan de blacklist met een TTL
func (r *RedisCacheRepository) AddBlacklistedAccessToken(ctx context.Context, token string, expiresAt time.Time) error {
	// De waarde kan willekeurig zijn, zolang de sleutel bestaat
	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil // Token is al verlopen
	}
	err := r.client.Set(ctx, "blacklist:"+token, "true", ttl)
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}
	return nil
}

// IsAccessTokenBlacklisted controleert of een token op de blacklist staat
func (r *RedisCacheRepository) IsAccessTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	val, err := r.client.Get(ctx, "blacklist:"+token)
	if err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return false, nil // Niet gevonden is niet geblacklist
		}
		return false, fmt.Errorf("failed to check blacklist: %w", err)
	}
	return val == "true", nil
}
