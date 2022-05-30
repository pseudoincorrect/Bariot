package authClient

import (
	"context"
	"log"

	pb "github.com/pseudoincorrect/bariot/pkg/auth/grpc"
	"google.golang.org/grpc"
)

const (
	Admin = "admin"
	User  = "user"
	Thing = "thing"
)

type ctxt context.Context

type Auth interface {
	StartAuthClient() error
	IsAdmin(ctxt, string) (bool, error)
	IsWhichUser(ctxt, string) (string, string, error)
	IsWhichThing(ctxt, string) (string, error)
	GetThingToken(ctxt, string, string) (string, error)
	GetAdminToken(ctxt) (string, error)
	GetUserToken(ctxt, string) (string, error)
}

// Static type checking
var _ Auth = (*authClient)(nil)

func New(conf Conf) Auth {
	return &authClient{Conf: conf, Conn: nil, Client: nil}
}

type Conf struct {
	Port string
	Host string
}

type authClient struct {
	Conf   Conf
	Conn   *grpc.ClientConn
	Client pb.AuthClient
}

// StartAuthClient starts the auth client GRPC server
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

// IsAdmin checks if the user is an admin given a token
func (c *authClient) IsAdmin(ctx ctxt, jwt string) (bool, error) {
	claims, err := c.Client.GetClaimsUserToken(ctx, &pb.GetClaimsUserTokenRequest{Jwt: jwt})
	if err != nil {
		log.Println("IsWhichUser GetClaimsUserToken error:", err)
		return false, err
	}
	return claims.GetRole() == Admin, nil
}

// IsWhichUser checks if the user is a user given a token
func (c *authClient) IsWhichUser(ctx ctxt, jwt string) (string, string, error) {
	claims, err := c.Client.GetClaimsUserToken(ctx, &pb.GetClaimsUserTokenRequest{Jwt: jwt})
	if err != nil {
		log.Println("IsWhichUser GetClaimsUserToken error:", err)
		return "", "", err
	}
	return claims.GetRole(), claims.GetSubject(), nil
}

// IsWhichThing whom a thing belong to, given a token
func (c *authClient) IsWhichThing(ctx ctxt, jwt string) (string, error) {
	claims, err := c.Client.GetClaimsThingToken(ctx, &pb.GetClaimsThingTokenRequest{Jwt: jwt})
	if err != nil {
		log.Println("IsWhichThing GetClaimsThingToken error:", err)
		return "", err
	}
	return claims.GetSubject(), nil
}

// GetThingToken returns a token for a thing given a user id and a thing id
func (c *authClient) GetThingToken(ctx ctxt, thingId string, userId string) (string, error) {
	res, err := c.Client.GetThingToken(ctx, &pb.GetThingTokenRequest{ThingId: thingId, UserId: userId})
	if err != nil {
		log.Println("GetThingToken error:", err)
		return "", err
	}
	return res.Jwt, nil
}

// GetAdminToken returns a token for an admin
func (c *authClient) GetAdminToken(ctx ctxt) (string, error) {
	resToken, err := c.Client.GetAdminToken(ctx, &pb.GetAdminTokenRequest{})
	if err != nil {
		log.Println("GRPC get admin token error:", err)
		return "", err
	}
	return resToken.GetJwt(), nil
}

// GetUserToken returns a token for a user
func (c *authClient) GetUserToken(ctx ctxt, userId string) (string, error) {
	resToken, err := c.Client.GetUserToken(ctx, &pb.GetUserTokenRequest{UserId: userId})
	if err != nil {
		log.Println("GRPC get admin token error:", err)
		return "", err
	}
	return resToken.GetJwt(), nil
}
