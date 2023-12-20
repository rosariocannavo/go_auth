package redis

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/rosariocannavo/go_auth/config"
)

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
	})
}
