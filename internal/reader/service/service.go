package service

import (
	"context"

	natsGo "github.com/nats-io/nats.go"
	auth "github.com/pseudoincorrect/bariot/pkg/auth/client"
	nats "github.com/pseudoincorrect/bariot/pkg/nats/client"
	things "github.com/pseudoincorrect/bariot/pkg/things/client"
	e "github.com/pseudoincorrect/bariot/pkg/utils/errors"
)

const natsThingsSubject = "thingsMsg"

type Reader interface {
	AuthorizeSingleThing(userToken string, thingId string) error
	ReceiveThingData(thingId string, handler func(string), stop chan bool) error
}

type reader struct {
	auth   auth.Auth
	things things.Things
	nats   nats.Nats
}

var _ Reader = (*reader)(nil)

// New creates a new reader service
func New(a auth.Auth, t things.Things, n nats.Nats) reader {
	return reader{a, t, n}
}

// AuthorizeSingleThing Check whether a thingId belong to a user (ID) with a user token
func (s *reader) AuthorizeSingleThing(userToken string, thingId string) error {
	ctx := context.Background()
	_, userId, err := s.auth.IsWhichUser(ctx, userToken)
	if err != nil {
		return e.Handle(e.ErrAuthn, err, "unauthorized")
	}
	userId2, err := s.things.GetUserOfThing(ctx, thingId)
	if err != nil {
		return err
	}
	if userId != userId2 {
		return e.ErrAuthz
	}
	return nil
}

// ReceiveThingData will connect to nats topic to receive the corresponding data and channel them
func (s *reader) ReceiveThingData(
	thingId string, handler func(string), stop chan bool,
) error {
	subject := natsThingsSubject + "." + thingId
	sub, err := s.nats.Subscribe(
		subject,
		GetReceiveThingIdDataHandler(handler),
	)
	if err != nil {
		return err
	}
	<-stop
	err = sub.Unsubscribe()
	if err != nil {
		e.Handle(e.ErrNats, err, "nat unsubscribe")
	}
	return err
}

func GetReceiveThingIdDataHandler(handler func(string)) natsGo.MsgHandler {
	return func(msg *natsGo.Msg) {
		// logger.Debug("--- GetReceiveThingIdDataHandler ---")
		// logger.Debug(msg.Subject)
		// logger.Debug(string(msg.Data))
		// logger.Debug("-------- Got a msg from NATS, sending to WEBSOCKETS -----")
		handler(string(msg.Data))
	}
}
