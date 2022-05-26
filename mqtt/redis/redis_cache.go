package cache

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/pseudoincorrect/bariot/pkg/env"
)

func Connect() {

	redis_host := env.GetEnv("REDIS_HOST")

	client := redis.NewClient(&redis.Options{
		Addr:     redis_host,
		Password: "",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	log.Println(pong, err)
}
