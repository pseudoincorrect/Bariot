package mqtt

import (
	"context"
	"log"
	"strings"

	authClient "github.com/pseudoincorrect/bariot/pkg/auth/client"
	"github.com/pseudoincorrect/bariot/pkg/errors"
)

type Authorizer func(topic string, jwt string) error

func CreateAuthorizer(auth authClient.Auth) (Authorizer, error) {
	authorizer := func(topic string, token string) error {
		thingId, err := auth.IsWhichThing(context.Background(), token)
		if err != nil {
			log.Println("MQTT token AUTHENTICATION error")
			return errors.ErrAuthorization
		}
		topicThingId, _ := extractThingIdFromTopic(topic)
		if topicThingId != thingId {
			log.Println("Thing", thingId, "NOT AUTHORIZED to publish on topic", topic)
			return errors.ErrAuthorization
		}
		return nil
	}
	return authorizer, nil
}

func extractThingIdFromTopic(topic string) (string, error) {
	thingId := strings.Split(topic, "/")[1]
	return thingId, nil
}
