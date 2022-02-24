package client

import (
	"context"
	"log"

	pb "github.com/pseudoincorrect/bariot/pkg/grpc/auth"
	"google.golang.org/grpc"
)

const (
	Admin = "admin"
	User  = "user"
	Thing = "thing"
)

type Auth interface {
	StartAuthClient() error
	IsAdmin(context.Context, string) (bool, error)
	IsWhichUser(context.Context, string) (string, string, error)
}

var _ Auth = (*authClient)(nil)

func New(conf AuthClientConf) Auth {
	return &authClient{Conf: conf, Conn: nil, Client: nil}
}

type AuthClientConf struct {
	Port string
	Host string
}

type authClient struct {
	Conf   AuthClientConf
	Conn   *grpc.ClientConn
	Client pb.AuthClient
}

func (c *authClient) StartAuthClient() error {
	addr := c.Conf.Host + ":" + c.Conf.Port
	log.Println("init user service GRPC client to ", addr)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Println("did not connect:", err)
		return err
	}
	// defer conn.Close()
	c.Conn = conn
	c.Client = pb.NewAuthClient(conn)
	return nil
}

func (c *authClient) IsAdmin(ctx context.Context, jwt string) (bool, error) {
	claims, err := c.Client.GetClaimsToken(ctx, &pb.GetClaimsTokenRequest{Jwt: jwt})
	if err != nil {
		log.Println("IsWhichUser GetClaimsToken error:", err)
		return false, err
	}
	return claims.GetRole() == Admin, nil
}

func (c *authClient) IsWhichUser(ctx context.Context, jwt string) (string, string, error) {
	claims, err := c.Client.GetClaimsToken(ctx, &pb.GetClaimsTokenRequest{Jwt: jwt})
	if err != nil {
		log.Println("IsWhichUser GetClaimsToken error:", err)
		return "", "", err
	}
	return claims.GetRole(), claims.GetSubject(), nil
}
