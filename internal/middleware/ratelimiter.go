package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v10"
	"github.com/rosariocannavo/go_auth/internal/nats"
	"github.com/rosariocannavo/go_auth/internal/ratelimiter"
)

type RateLimitMiddleware struct {
	RedisLimiter *ratelimiter.RedisRateLimiter
}

func NewRateLimitMiddleware() *RateLimitMiddleware {
	return &RateLimitMiddleware{ratelimiter.SetupRedisRateLimiter()}
}

const RateRequest = "rate_request_%s"

func (r *RateLimitMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := r.RedisLimiter.Allow(c, fmt.Sprintf(RateRequest, "userName"), redis_rate.Limit{
			Rate:   1, //max req per client
			Burst:  5,
			Period: time.Second,
		})

		if err != nil || res.Allowed <= 0 {

			message := fmt.Sprintf("Timestamp: %s | Handler: %s | Status: %d | Response: %s", time.Now().UTC().Format(time.RFC3339), "middleware/RateLimiter", http.StatusTooManyRequests, "error: Too many requests")
			nats.NatsConnection.PublishMessage(message)

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			return
		}

		c.Next()
	}
}
