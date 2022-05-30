package cache

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	appErr "github.com/pseudoincorrect/bariot/pkg/errors"
)

type CacheRes int64

const (
	CacheOk CacheRes = iota
	CacheHit
	CacheMiss
	CacheError
)

type ThingCache interface {
	deleteToken(token string) error
	deleteThingId(thingId string) error
	Connect() error
	GetThingIdByToken(token string) (_ CacheRes, thingId string, err error)
	GetTokenByThingId(thingId string) (_ CacheRes, token string, err error)
	SetTokenWithThingId(token string, thingId string) error
	DeleteTokenAndTokenByThingId(thingId string) error
}

var _ ThingCache = (*cache)(nil)

type cache struct {
	client *redis.Client
	conf   Conf
}

func New(conf Conf) ThingCache {
	return &cache{client: nil, conf: conf}
}

type Conf struct {
	RedisHost string
	RedisPort string
}

func (c *cache) Connect() error {
	var ctx = context.Background()
	addr := c.conf.RedisHost + ":" + c.conf.RedisPort

	c.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	res, err := c.client.Ping(ctx).Result()
	if err != nil {
		log.Panic("Could not connect to Redis")
	}
	log.Println("Ping res = ", res)
	return nil
}

func (c *cache) deleteToken(token string) error {
	var ctx = context.Background()
	key, err := c.client.Del(ctx, token).Result()
	if err != nil {
		log.Println("Error deleteToken")
		return appErr.ErrCache
	}
	log.Println("Deleted token key", key)
	return nil
}

func (c *cache) deleteThingId(thingId string) error {
	var ctx = context.Background()
	key, err := c.client.Del(ctx, thingId).Result()
	if err != nil {
		log.Println("Error deleteToken")
		return appErr.ErrCache
	}
	log.Println("Deleted thingId key", key)
	return nil
}

func (c *cache) GetThingIdByToken(token string) (
	_ CacheRes, thingId string, err error) {
	var ctx = context.Background()
	thingId, err = c.client.Get(ctx, token).Result()

	if err == redis.Nil {
		log.Println("ThingCache token MISS")
		return CacheMiss, "", nil
	} else if err != nil {
		log.Println("Error ThingCache")
		return CacheError, "", appErr.ErrCache
	}

	log.Println("GetThingIdByToken(): ", token[0:10], " = ", thingId)
	return CacheHit, thingId, nil
}

func (c *cache) GetTokenByThingId(thingId string) (
	_ CacheRes, token string, err error) {
	var ctx = context.Background()
	token, err = c.client.Get(ctx, thingId).Result()

	if err == redis.Nil {
		log.Println("ThingCache thingId MISS")
		return CacheMiss, "", nil
	} else if err != nil {
		log.Println("Error ThingCache")
		return CacheError, "", appErr.ErrCache
	}

	log.Println("GetTokenByThingId(): ", thingId, " = ", token)
	return CacheHit, token, nil
}

func (c *cache) SetTokenWithThingId(token string, thingId string) error {
	var ctx = context.Background()

	err := c.client.Set(ctx, token, thingId, 0).Err()
	if err != nil {
		log.Println("Error ThingCache, adding token (key) to cache")
		return appErr.ErrCache
	}

	err = c.client.Set(ctx, thingId, token, 0).Err()
	if err != nil {
		log.Println("Error ThingCache, adding thingId (key) to cache")
		return appErr.ErrCache
	}

	return nil
}

func (c *cache) DeleteTokenAndTokenByThingId(thingId string) error {
	res, token, err := c.GetTokenByThingId(thingId)
	if err != nil {
		return appErr.ErrCache
	}
	if res == CacheMiss {
		return nil
	}
	if res == CacheHit {
		err = c.deleteToken(token)
		if err != nil {
			return appErr.ErrCache
		}
		err = c.deleteThingId(thingId)
		if err != nil {
			return appErr.ErrCache
		}
	}
	log.Println("Deleted Token", token[0:10], "and thingId", thingId)
	return nil
}
