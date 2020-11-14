package db

import (
    "fmt"

    "github.com/spf13/viper"
    "github.com/go-redis/redis"
)

var RedisClient *redis.Client

func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     fmt.Sprintf("%s:%s", viper.GetString("redis.host"), viper.GetString("redis.port")),
        Password: viper.GetString("redis.auth"),
        DB:       0,
    })

    _, err := RedisClient.Ping().Result()
    if err != nil {
        panic("redis ping error")
    }
}