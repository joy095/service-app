package middleware

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	ginmiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

func createRedisStore() (limiter.Store, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Test the connection
	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	store, err := redisstore.NewStoreWithOptions(rdb, limiter.StoreOptions{
		Prefix:   "rate_limiter",
		MaxRetry: 3,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create redis store: %w", err)
	}
	return store, nil
}

// ParseCustomRate allows formats like "10-2m", "30-20m", "5-1h", etc.
func ParseCustomRate(rateStr string) (limiter.Rate, error) {
	parts := strings.Split(rateStr, "-")
	if len(parts) != 2 {
		return limiter.Rate{}, fmt.Errorf("invalid rate format: %s", rateStr)
	}
	limit, err := strconv.Atoi(parts[0])
	if err != nil {
		return limiter.Rate{}, fmt.Errorf("invalid limit: %s", parts[0])
	}

	durationStr := parts[1]
	var period time.Duration

	if strings.HasSuffix(durationStr, "m") {
		minutes, err := strconv.Atoi(strings.TrimSuffix(durationStr, "m"))
		if err != nil {
			return limiter.Rate{}, err
		}
		period = time.Duration(minutes) * time.Minute
	} else if strings.HasSuffix(durationStr, "h") {
		hours, err := strconv.Atoi(strings.TrimSuffix(durationStr, "h"))
		if err != nil {
			return limiter.Rate{}, err
		}
		period = time.Duration(hours) * time.Hour
	} else {
		return limiter.Rate{}, fmt.Errorf("unsupported period: %s", durationStr)
	}

	return limiter.Rate{
		Period: period,
		Limit:  int64(limit),
	}, nil
}

// NewRateLimiter creates middleware with custom periods like "10-2m"
func NewRateLimiter(rateStr string) gin.HandlerFunc {
	rate, err := ParseCustomRate(rateStr)
	if err != nil {
		log.Printf("Error parsing rate: %v", err)
		// Return a fallback middleware that just passes through
		return func(c *gin.Context) {
			c.Next()
		}
	}

	store, err := createRedisStore()
	if err != nil {
		log.Printf("Error creating Redis store: %v", err)
		// Return a fallback middleware that just passes through
		return func(c *gin.Context) {
			c.Next()
		}
	}

	limiterInstance := limiter.New(store, rate)
	return ginmiddleware.NewMiddleware(limiterInstance)
}

// CombinedRateLimiter accepts multiple custom rate strings
func CombinedRateLimiter(rateStrings ...string) gin.HandlerFunc {
	middlewares := make([]gin.HandlerFunc, len(rateStrings))
	for i, rateStr := range rateStrings {
		middlewares[i] = NewRateLimiter(rateStr)
	}
	return func(c *gin.Context) {
		for _, mw := range middlewares {
			mw(c)
			if c.IsAborted() {
				return
			}
		}
	}
}
