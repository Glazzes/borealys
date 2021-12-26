package config

import (
	"context"
	"github.com/go-redis/redis"
)

var (
	Ctx = context.Background()
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "cache:6379",
		DB: 0,
		Password: "",
	})
)