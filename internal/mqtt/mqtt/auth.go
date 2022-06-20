package mqtt

import (
	"context"
	"strings"

	authClient "github.com/pseudoincorrect/bariot/pkg/auth/client"
	cdb "github.com/pseudoincorrect/bariot/pkg/cache"
	e "github.com/pseudoincorrect/bariot/pkg/utils/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
)

type Authorizer func(topic string, jwt string) error

// CreateAuthorizer creates an authorizer function that can be used to authorize a mqtt topic
func CreateAuthorizer(auth authClient.Auth, cache cdb.ThingCache) (Authorizer, error) {
	authorizer := func(topic string, token string) error {
		thingId, err := authenticate(auth, cache, token)
		if err != nil {
			return err
		}
		topicThingId, _ := extractThingIdFromTopic(topic)
		if topicThingId != thingId {
			logger.Error("Thing", thingId, "NOT AUTHORIZED to publish on topic", topic)
			return e.ErrAuthz
		}
		return nil
	}
	return authorizer, nil
}

// extractThingIdFromTopic extracts the thing id from a topic
func extractThingIdFromTopic(topic string) (string, error) {
	thingId := strings.Split(topic, "/")[1]
	return thingId, nil
}

// authenticate authenticates a token and returns the thing id
func authenticate(auth authClient.Auth, cache cdb.ThingCache, token string) (string, error) {
	res, thingId, err := cache.GetThingIdByToken(token)
	if err != nil {
		logger.Error("Error Redis cache")
		return "", e.ErrCache
	}
	if res == cdb.CacheHit {
		return thingId, nil
	}

	thingId, err = auth.IsWhichThing(context.Background(), token)
	if err != nil {
		logger.Error("MQTT token AUTHENTICATION error")
		return "", e.ErrAuthz
	}

	err = cache.SetTokenWithThingId(token, thingId)
	if err != nil {
		logger.Error("Error Redis cache")
		return "", err
	}
	return thingId, nil
}
