package caches

import (
	"errors"

	"github.com/pseudoincorrect/bariot/pkg/cache"
)

// Static type checking
var _ cache.ThingCache = (*CacheMock)(nil)

type CacheMock struct {
	ThrowErr string
}

func NewCacheMock() CacheMock {
	return CacheMock{}
}

// Connect to redis
func (c *CacheMock) Connect() error {
	if c.ThrowErr != "" {
		return errors.New(c.ThrowErr)
	}
	return nil
}

// Delete token key from cache
func (c *CacheMock) DeleteToken(token string) error {
	if c.ThrowErr != "" {
		return errors.New(c.ThrowErr)
	}
	return nil
}

// Delete thingId key from cache
func (c *CacheMock) DeleteThingId(thingId string) error {
	if c.ThrowErr != "" {
		return errors.New(c.ThrowErr)
	}
	return nil
}

// Get thingId by token
func (c *CacheMock) GetThingIdByToken(token string) (
	_ cache.CacheRes, thingId string, err error) {
	if c.ThrowErr != "" {
		return cache.CacheError, "", errors.New(c.ThrowErr)
	}
	return cache.CacheHit, thingId, nil
}

// Get token by thingId
func (c *CacheMock) GetTokenByThingId(thingId string) (
	_ cache.CacheRes, token string, err error) {
	if c.ThrowErr != "" {
		return cache.CacheError, "", errors.New(c.ThrowErr)
	}
	return cache.CacheHit, token, nil
}

// Set token key with thingId value
func (c *CacheMock) SetTokenWithThingId(token string, thingId string) error {
	if c.ThrowErr != "" {
		return errors.New(c.ThrowErr)
	}
	return nil
}

// DeleteTokenAndTokenByThingId delete token and tokenByThingId keys
func (c *CacheMock) DeleteTokenAndThingByThingId(thingId string) error {
	if c.ThrowErr != "" {
		return errors.New(c.ThrowErr)
	}
	return nil
}
