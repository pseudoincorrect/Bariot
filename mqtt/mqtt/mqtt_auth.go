package mqtt

import (
	"encoding/json"
	"log"

	"github.com/mainflux/senml"
	authClient "github.com/pseudoincorrect/bariot/pkg/auth/client"
)

type Authorizer func(topic string, jwt string) error

func CreateAuthorizer(auth authClient.Auth) (Authorizer, error) {
	authorizer := func(topic string, payload string) error {
		log.Println("Authorizing MQTT msg")
		log.Println("topic = ", topic)
		log.Println("jwt = ", payload)
		return nil
	}
	return authorizer, nil
}

type AuthenticatedMsg struct {
	Token   string `json:"token"`
	Sensors senml.Pack
}

func ExtractJwt(payload []byte) (string, error) {
	msg := AuthenticatedMsg{}
	err := json.Unmarshal(payload, &msg)
	if err != nil {
		log.Println(err)
		return "", err
	}
	log.Println("JSON decoded jwt = ", msg.Token)
	log.Println("JSON decoded data = ", msg.Sensors)
	return msg.Token, nil
}
