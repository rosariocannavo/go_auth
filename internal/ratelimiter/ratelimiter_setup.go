package ratelimiter

import (
	"context"
	"log"

	redisrate "github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

type RedisRateLimiter struct {
	*redisrate.Limiter
}

func SetupRedisRateLimiter(connString string) *RedisRateLimiter {

	//TODO move to init
	client := redis.NewClient(&redis.Options{
		Addr: connString,
	})
	//
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Error connecting to Redis:", err)
	}
	return &RedisRateLimiter{redisrate.NewLimiter(client)}
}
