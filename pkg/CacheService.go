package pkg

import (
	"encoding/json"
	"github.com/glazzes/borealys/pkg/config"
	"log"
	"time"
)

var (
	duration = time.Second * 240
	infoLogger = log.Logger{}
)

type CacheService interface {
	SetOutput(key string, output []string)
	GetOutputByKey(key string) []string
}

type SimpleCacheService struct {}

func (c *SimpleCodeRunnerService) SetOutput(key, output string) {
	if err := config.RedisClient.Set(key, output, duration).Err(); err != nil {
		log.Fatalln(err)
	}
}

func (c *SimpleCacheService) GetOutputByKey(key string) []string {
	result,err := config.RedisClient.Get(key).Result()
	checkNilErr(err)

	deserialized := []string{}
	if err = json.Unmarshal([]byte(result), &deserialized); err != nil{
		log.Fatal(err)
	}

	return deserialized
}

func checkNilErr(err error){
	log.Fatal(err)
}
