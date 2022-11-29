package redis

import (
	"github.com/go-redis/redis"
)

func NewRedis() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
	})

	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}

	return client, nil
}
