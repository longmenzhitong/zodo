package redish

import (
	"github.com/go-redis/redis"
	"zodo/internal/conf"
)

var client *redis.Client

func init() {
	if conf.IsRedisStorage() {
		initClient()
	}
}

func Client() *redis.Client {
	if client == nil {
		initClient()
	}
	return client
}

func initClient() {
	c := redis.NewClient(&redis.Options{
		Addr:     conf.Data.Storage.Redis.Address,
		Password: conf.Data.Storage.Redis.Password,
		DB:       conf.Data.Storage.Redis.Db,
	})
	_, err := c.Ping().Result()
	if err != nil {
		panic(err)
	}
	client = c
}
