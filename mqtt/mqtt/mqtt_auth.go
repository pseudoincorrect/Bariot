package mqtt

import (
	"context"
	"log"
	"strings"

	authClient "github.com/pseudoincorrect/bariot/pkg/auth/client"
	cdb "github.com/pseudoincorrect/bariot/pkg/cache"
	appErrors "github.com/pseudoincorrect/bariot/pkg/errors"
)

type Authorizer func(topic string, jwt string) error

func CreateAuthorizer(auth authClient.Auth, cache cdb.ThingCache) (Authorizer, error) {
	authorizer := func(topic string, token string) error {
		thingId, err := authenticate(auth, cache, token)
		if err != nil {
			return err
		}
		topicThingId, _ := extractThingIdFromTopic(topic)
		if topicThingId != thingId {
			log.Println("Thing", thingId, "NOT AUTHORIZED to publish on topic", topic)
			return appErrors.ErrAuthorization
		}
		return nil
	}
	return authorizer, nil
}

func extractThingIdFromTopic(topic string) (string, error) {
	thingId := strings.Split(topic, "/")[1]
	return thingId, nil
}

func authenticate(auth authClient.Auth, cache cdb.ThingCache, token string) (string, error) {
	log.Println("authenticate()")
	res, thingId, err := cache.GetThingIdByToken(token)
	if err != nil {
		log.Println("Error Redis cache")
		return "", appErrors.ErrCache
	}
	if res == cdb.CacheHit {
		return thingId, nil
	}

	thingId, err = auth.IsWhichThing(context.Background(), token)
	if err != nil {
		log.Println("MQTT token AUTHENTICATION error")
		return "", appErrors.ErrAuthorization
	}

	err = cache.SetTokenWithThingId(token, thingId)
	if err != nil {
		log.Println("Error Redis cache")
		return "", err
	}
	log.Println("Token gotten and cached, token: ", token[0:10])
	return thingId, nil
}
