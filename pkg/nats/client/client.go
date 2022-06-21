package client

import (
	"fmt"
	"log"
	"time"

	natsGo "github.com/nats-io/nats.go"
	e "github.com/pseudoincorrect/bariot/pkg/utils/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
)

const natsThingsSubject = "thingsMsg"

type Nats interface {
	Connect(opts []natsGo.Option) error
	Disconnect()
	Subscribe(subject string, queue string, handler natsGo.MsgHandler) (*natsGo.Subscription, error)
	CreatePublisher() NatsPubType
	Publish(subject string, payload string) error
}

type nats struct {
	host string
	port string
	conn *natsGo.Conn
}

var _ Nats = (*nats)(nil)

type Conf struct {
	Host string
	Port string
}

func New(conf Conf) nats {
	return nats{
		host: conf.Host,
		port: conf.Port,
		conn: nil,
	}
}

// connect setup a connection to a Nats server
func (nps *nats) Connect(opts []natsGo.Option) error {
	natsUrl := "nats://" + nps.host + ":" + nps.port
	logger.Info("Connecting to NATS Server:", natsUrl)
	nc, err := natsGo.Connect(natsUrl, opts...)
	if err != nil {
		return e.Handle(e.ErrConn, err, "nats connect")
	}
	nps.conn = nc
	return nil
}

//Disconnect disconnects from the NATS server.
func (nps *nats) Disconnect() {
	nps.conn.Close()
}

// natsSubscribe subscribe to a topic/subject with a custom handler and a queue
func (nps *nats) Subscribe(subject string, queue string, handler natsGo.MsgHandler) (*natsGo.Subscription, error) {
	sub, err := nps.conn.QueueSubscribe(subject, queue, handler)
	if err != nil {
		err = e.Handle(e.ErrNats, err, "nats subscribe thing id")
		return nil, err
	}
	nps.conn.Flush()
	err = nps.conn.LastError()
	if err != nil {
		err = e.Handle(e.ErrNats, err, "nats flush")
		return nil, err
	}
	return sub, nil
}

// Publish publishes a message to the NATS server.
func (nps *nats) Publish(subject string, payload string) error {
	if err := nps.conn.Publish(subject, []byte(payload)); err != nil {
		logger.Error("Error publishing message: ", err)
		return err
	}
	return nil
}

// natsSetupConnOptions setup the Nats connection option such as reconnect
func NatsSetupConnOptions(name string) []natsGo.Option {
	opts := []natsGo.Option{natsGo.Name(name)}
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second
	opts = append(opts, natsGo.ReconnectWait(reconnectDelay))
	opts = append(opts, natsGo.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, natsGo.DisconnectErrHandler(func(nc *natsGo.Conn, err error) {
		str := fmt.Sprintf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
		logger.Info(str)
	}))
	opts = append(opts, natsGo.ReconnectHandler(func(nc *natsGo.Conn) {
		logger.Info("Reconnected [", nc.ConnectedUrl(), "]")
	}))
	opts = append(opts, natsGo.ClosedHandler(func(nc *natsGo.Conn) {
		log.Panic("Exiting:", nc.LastError())
	}))
	return opts
}

type NatsPubType func(thingId string, payload string) error

// CreatePublisher return a function to publish on a given subject.
func (nps *nats) CreatePublisher() NatsPubType {
	return func(thingId string, payload string) error {
		err := nps.Publish(formatThingIdSubject(thingId), payload)
		return err
	}
}

// formatThingIdSubject return the subject for a thingId
func formatThingIdSubject(thingId string) string {
	return natsThingsSubject + "." + thingId
}
