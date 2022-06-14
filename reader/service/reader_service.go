package service

import (
	"github.com/nats-io/nats.go"
	auth "github.com/pseudoincorrect/bariot/pkg/auth/client"
)

type ReaderSvc interface {
	AuthorizeSingleThing(userToken string, thingId string) (bool, error)
}

type readerSvc struct {
	auth     auth.Auth
	natsConn *nats.Conn
}

var _ ReaderSvc = (*readerSvc)(nil)

func (s *readerSvc) AuthorizeSingleThing(userToken string, thingId string) (bool, error) {
	return true, nil
}
