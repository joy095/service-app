// internal/middleware/rate_limiter.go
package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	ginmiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

// NewRateLimiter returns a middleware with a specific rate like "5-M" or "10-S"
func NewRateLimiter(rateString string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(rateString)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	store, err := redisstore.NewStoreWithOptions(rdb, limiter.StoreOptions{
		Prefix:   "rate_limiter",
		MaxRetry: 3,
	})
	if err != nil {
		panic(err)
	}

	limiterInstance := limiter.New(store, rate)
	return ginmiddleware.NewMiddleware(limiterInstance)
}
