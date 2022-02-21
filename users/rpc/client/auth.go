package client

import (
	"context"
	"log"

	pb "github.com/pseudoincorrect/bariot/users/rpc/auth"
	"google.golang.org/grpc"
)

const (
	admin = "admin"
	user  = "user"
	thing = "thing"
)

type Auth interface {
	StartAuthClient() error
	GetAdminToken(context.Context) (string, error)
	GetUserToken(context.Context, string) (string, error)
	IsAdmin(context.Context, string) (bool, error)
	IsWhichUser(context.Context, string) (bool, string, error)
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

func (c *authClient) GetAdminToken(ctx context.Context) (string, error) {
	resToken, err := c.Client.GetAdminToken(ctx, &pb.GetAdminTokenRequest{})
	if err != nil {
		log.Println("GRPC get admin token error:", err)
		return "", err
	}
	return resToken.GetJwt(), nil
}

func (c *authClient) GetUserToken(ctx context.Context, userId string) (string, error) {
	resToken, err := c.Client.GetUserToken(ctx, &pb.GetUserTokenRequest{UserId: userId})
	if err != nil {
		log.Println("GRPC get admin token error:", err)
		return "", err
	}
	return resToken.GetJwt(), nil
}

func (c *authClient) IsAdmin(ctx context.Context, jwt string) (bool, error) {
	claims, err := c.Client.GetClaimsToken(ctx, &pb.GetClaimsTokenRequest{Jwt: jwt})
	if err != nil {
		return false, err
	}
	return claims.GetRole() == admin, nil
}

func (c *authClient) IsWhichUser(ctx context.Context, jwt string) (bool, string, error) {
	claims, err := c.Client.GetClaimsToken(ctx, &pb.GetClaimsTokenRequest{Jwt: jwt})
	if err != nil {
		return false, "", err
	}
	return claims.GetRole() == user, claims.GetSubject(), nil
}
