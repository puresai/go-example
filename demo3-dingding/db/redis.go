package db

import (
    "fmt"
	"context"
	"time"

    "github.com/spf13/viper"
    "github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
		Password: viper.GetString("redis.auth"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.MaxActive"),DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolTimeout:  30 * time.Second,
	})


    _, err := RedisClient.Ping(Ctx).Result()
    if err != nil {
        panic("redis ping error")
    }
}