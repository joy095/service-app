package middleware

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	ginmiddleware "github.com/ulule/limiter/v3/drivers/middleware/gin"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

func RateLimiterMiddleware() gin.HandlerFunc {
	// Define the rate limit: 100 requests per minute
	rate, err := limiter.NewRateFromFormatted("100-M")
	if err != nil {
		panic(err)
	}

	// Setup go-redis v9 client
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Create store using go-redis v9
	store, err := redisstore.NewStoreWithOptions(rdb, limiter.StoreOptions{
		Prefix:   "rate_limiter",
		MaxRetry: 3,
	})
	if err != nil {
		panic(err)
	}

	// Create limiter instance and wrap with Gin middleware
	limiterInstance := limiter.New(store, rate)
	return ginmiddleware.NewMiddleware(limiterInstance)
}
