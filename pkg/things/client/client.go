package client

import (
	"context"

	pb "github.com/pseudoincorrect/bariot/pkg/things/grpc"
	e "github.com/pseudoincorrect/bariot/pkg/utils/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Things interface {
	StartThingsClient() error
	GetUserOfThing(ctx context.Context, thingId string) (userId string, err error)
}

// Static type checking
var _ Things = (*thingsClient)(nil)

func New(conf Conf) thingsClient {
	return thingsClient{Conf: conf, Conn: nil, Client: nil}
}

type Conf struct {
	Port string
	Host string
}

type thingsClient struct {
	Conf   Conf
	Conn   *grpc.ClientConn
	Client pb.ThingsClient
}

// StartThingsClient starts the auth client GRPC server
func (c *thingsClient) StartThingsClient() error {
	addr := c.Conf.Host + ":" + c.Conf.Port
	logger.Info("init user service GRPC client to ", addr)
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Error("did not connect:", err)
		return err
	}
	c.Conn = conn
	c.Client = pb.NewThingsClient(conn)
	return nil
}

// IsAdmin checks if the user is an admin given a token
func (c *thingsClient) GetUserOfThing(ctx context.Context, thingId string) (userId string, err error) {
	res, err := c.Client.GetUserOfThing(ctx, &pb.GetUserOfThingRequest{ThingId: thingId})
	if err != nil {
		logger.Error("IsWhichUser GetClaimsUserToken error:", err)
		return "", e.Handle(e.ErrGrpc, err, "get user of thing grpc")
	}
	return res.UserId, nil
}
