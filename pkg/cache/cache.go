package cache

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
	e "github.com/pseudoincorrect/bariot/pkg/utils/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
)

type CacheRes int64

const (
	CacheOk CacheRes = iota
	CacheHit
	CacheMiss
	CacheError
)

type ThingCache interface {
	Connect() error
	DeleteToken(token string) error
	DeleteThingId(thingId string) error
	GetThingIdByToken(token string) (_ CacheRes, thingId string, err error)
	GetTokenByThingId(thingId string) (_ CacheRes, token string, err error)
	SetTokenWithThingId(token string, thingId string) error
	DeleteTokenAndThingByThingId(thingId string) error
}

// Static type checking
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

// Connect to redis
func (c *cache) Connect() error {
	var ctx = context.Background()
	addr := c.conf.RedisHost + ":" + c.conf.RedisPort

	c.client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := c.client.Ping(ctx).Result()
	if err != nil {
		log.Panic("Could not connect to Redis")
	}
	return nil
}

// Delete token key from cache
func (c *cache) DeleteToken(token string) error {
	var ctx = context.Background()
	_, err := c.client.Del(ctx, token).Result()
	if err != nil {
		logger.Error("Error DeleteToken")
		return e.ErrCache
	}
	return nil
}

// Delete thingId key from cache
func (c *cache) DeleteThingId(thingId string) error {
	var ctx = context.Background()
	_, err := c.client.Del(ctx, thingId).Result()
	if err != nil {
		logger.Error("Error DeleteToken")
		return e.ErrCache
	}
	return nil
}

// Get thingId by token
func (c *cache) GetThingIdByToken(token string) (
	_ CacheRes, thingId string, err error) {
	var ctx = context.Background()
	thingId, err = c.client.Get(ctx, token).Result()

	if err == redis.Nil {
		logger.Error("ThingCache token MISS")
		return CacheMiss, "", nil
	} else if err != nil {
		logger.Error("Error ThingCache")
		return CacheError, "", e.ErrCache
	}
	return CacheHit, thingId, nil
}

// Get token by thingId
func (c *cache) GetTokenByThingId(thingId string) (
	_ CacheRes, token string, err error) {
	var ctx = context.Background()
	token, err = c.client.Get(ctx, thingId).Result()

	if err == redis.Nil {
		logger.Error("ThingCache thingId MISS")
		return CacheMiss, "", nil
	} else if err != nil {
		logger.Error("Error ThingCache")
		return CacheError, "", e.ErrCache
	}

	return CacheHit, token, nil
}

// Set token key with thingId value
func (c *cache) SetTokenWithThingId(token string, thingId string) error {
	var ctx = context.Background()

	err := c.client.Set(ctx, token, thingId, 0).Err()
	if err != nil {
		logger.Error("Error ThingCache, adding token (key) to cache")
		return e.ErrCache
	}

	err = c.client.Set(ctx, thingId, token, 0).Err()
	if err != nil {
		logger.Error("Error ThingCache, adding thingId (key) to cache")
		return e.ErrCache
	}

	return nil
}

// DeleteTokenAndThingByThingId delete token and tokenByThingId keys
func (c *cache) DeleteTokenAndThingByThingId(thingId string) error {
	res, token, err := c.GetTokenByThingId(thingId)
	if err != nil {
		return e.ErrCache
	}
	if res == CacheMiss {
		return nil
	}
	if res == CacheHit {
		err = c.DeleteToken(token)
		if err != nil {
			return e.ErrCache
		}
		err = c.DeleteThingId(thingId)
		if err != nil {
			return e.ErrCache
		}
	}
	return nil
}
