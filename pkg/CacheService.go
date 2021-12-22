package pkg

import (
	"github.com/glazzes/borealys/pkg/config"
	"time"
)

var (
	duration = time.Second * 240
)

type CacheService interface {
	SetKey(key string, output string)
	GetKey(key string) []string
}

type SimpleCacheService struct {}

func (c *SimpleCodeRunnerService) SetKey(key, output string)  {
	config.RedisClient.Set(key, output, duration)
}

func (c *SimpleCacheService) GetKey(key string) []string {
	return make([]string, 0)
}
