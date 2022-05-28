package cache

import (
	"log"

	"github.com/go-redis/redis"
)

type Cache interface {
	Connect() error
}

var _ Cache = (*cache)(nil)

type cache struct {
	client *redis.Client
	conf   Conf
}

func New(conf Conf) Cache {
	return &cache{client: nil, conf: conf}
}

type Conf struct {
	RedisHost string
	RedisPort string
}

func (c *cache) Connect() error {
	addr := c.conf.RedisHost + ":" + c.conf.RedisPort

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Panic("Could not connect to Redis")
	}

	return nil
}
