package config

import (
	"context"
	"github.com/go-redis/redis"
)

var (
	Ctx = context.Background()
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB: 0,
		Password: "",
	})
)