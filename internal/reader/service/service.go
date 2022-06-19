package service

import (
	"context"
	"log"

	"github.com/nats-io/nats.go"
	auth "github.com/pseudoincorrect/bariot/pkg/auth/client"
	e "github.com/pseudoincorrect/bariot/pkg/errors"
	things "github.com/pseudoincorrect/bariot/pkg/things/client"
)

type Reader interface {
	AuthorizeSingleThing(userToken string, thingId string) error
	ReceiveThingData(thingId string, thingData chan string, stop chan bool)
}

type reader struct {
	auth     auth.Auth
	things   things.Things
	natsConn *nats.Conn
}

var _ Reader = (*reader)(nil)

// New creates a new reader service
func New(a auth.Auth, t things.Things, n *nats.Conn) Reader {
	return &reader{a, t, n}
}

func (s *reader) AuthorizeSingleThing(userToken string, thingId string) error {
	ctx := context.Background()
	_, userId, err := s.auth.IsWhichUser(ctx, userToken)
	if err != nil {
		return e.Handle(e.ErrAuthn, err, "unauthorized")
	}
	userId2, err := s.things.GetUserOfThing(ctx, userId)
	if err != nil {
		return err
	}
	log.Println("userId ", userId)
	log.Println("userId2", userId2)
	if userId != userId2 {
		return e.ErrAuthz
	}
	return nil
}

func (s *reader) ReceiveThingData(thingId string, thingData chan string, stop chan bool) {
}
