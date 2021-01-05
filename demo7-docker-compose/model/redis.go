package model

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func ConnRedis(host, port, auth string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: auth,
		DB:       0,
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Fatalf("redis conn error", err)
	}
}
