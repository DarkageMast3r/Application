package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/time/rate"
)

func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	// Maak de oorspronkelijke limiter aan
	localLimiter := rate.NewLimiter(r, b)

	// Probeer een Redis client uit de context te halen
	return func(c *gin.Context) {
		// Eerst proberen we Redis als die beschikbaar is
		if redisClient, exists := c.Get("redis_client"); exists {
			rdb := redisClient.(*redis.Client)

			window := time.Duration(float64(time.Second)*(1/float64(r))) * time.Duration(b)
			key := fmt.Sprintf("rate_limit:%s:%s", c.ClientIP(), c.Request.URL.Path)

			current, err := rdb.Incr(c.Request.Context(), key).Result()
			if err == nil {
				if current == 1 {
					rdb.Expire(c.Request.Context(), key, window)
				}

				if current > int64(b) {
					c.String(http.StatusTooManyRequests, "Rate limit exceeded")
					c.Abort()
					return
				}

				c.Next()
				return
			}
			// Als Redis faalt, vallen we terug op de lokale limiter
		}

		// Fallback naar de originele implementatie
		if !localLimiter.AllowN(time.Now(), 1) {
			c.String(http.StatusTooManyRequests, "Rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}

// Helper om Redis client aan Gin context toe te voegen
func WithRedisClient(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("redis_client", rdb)
		c.Next()
	}
}
