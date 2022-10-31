package redish

import (
	"github.com/go-redis/redis"
	"zodo/internal/conf"
)

var Client *redis.Client

func init() {
	if conf.IsRedisStorage() {
		client := redis.NewClient(&redis.Options{
			Addr:     conf.Data.Storage.Redis.Address,
			Password: conf.Data.Storage.Redis.Password,
			DB:       conf.Data.Storage.Redis.Db,
		})
		_, err := client.Ping().Result()
		if err != nil {
			panic(err)
		}
		Client = client
	}
}
