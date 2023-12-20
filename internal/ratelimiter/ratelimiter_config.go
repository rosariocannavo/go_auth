package ratelimiter

import (
	"context"
	"log"

	redisrate "github.com/go-redis/redis_rate/v10"
	"github.com/rosariocannavo/go_auth/internal/redis"
)

type RedisRateLimiter struct {
	*redisrate.Limiter
}

func SetupRedisRateLimiter() *RedisRateLimiter {

	_, err := redis.Client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	return &RedisRateLimiter{redisrate.NewLimiter(redis.Client)}
}
