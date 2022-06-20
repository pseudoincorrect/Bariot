package client

import (
	"context"

	pb "github.com/pseudoincorrect/bariot/pkg/auth/grpc"
	"github.com/pseudoincorrect/bariot/pkg/utils/debug"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	Admin = "admin"
	User  = "user"
	Thing = "thing"
)

type Auth interface {
	StartAuthClient() error
	IsAdmin(ctx context.Context, token string) (iAdmin bool, err error)
	IsWhichUser(ctx context.Context, token string) (role string, userId string, err error)
	IsWhichThing(ctx context.Context, token string) (thingId string, err error)
	GetThingToken(ctx context.Context, thingId string, userId string) (token string, err error)
	GetAdminToken(ctx context.Context) (token string, err error)
	GetUserToken(ctx context.Context, userId string) (token string, err error)
}

// Static type checking
var _ Auth = (*authClient)(nil)

func New(conf Conf) authClient {
	return authClient{Conf: conf, Conn: nil, Client: nil}
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
	debug.LogInfo("init user service GRPC client to ", addr)
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		debug.LogError("did not connect:", err)
		return err
	}
	// defer conn.Close()
	c.Conn = conn
	c.Client = pb.NewAuthClient(conn)
	return nil
}

// IsAdmin checks if the user is an admin given a token
func (c *authClient) IsAdmin(ctx context.Context, token string) (isAdmin bool, err error) {
	claims, err := c.Client.GetClaimsUserToken(ctx, &pb.GetClaimsUserTokenRequest{Jwt: token})
	if err != nil {
		debug.LogError("IsWhichUser GetClaimsUserToken error:", err)
		return false, err
	}
	return claims.GetRole() == Admin, nil
}

// IsWhichUser checks if the user is a user given a token, return role, user id
func (c *authClient) IsWhichUser(ctx context.Context, token string) (role string, userId string, err error) {
	claims, err := c.Client.GetClaimsUserToken(ctx, &pb.GetClaimsUserTokenRequest{Jwt: token})
	if err != nil {
		debug.LogError("IsWhichUser GetClaimsUserToken error:", err)
		return "", "", err
	}
	return claims.GetRole(), claims.GetSubject(), nil
}

// IsWhichThing whom a thing belong to, given a Thing token, return thing ID
func (c *authClient) IsWhichThing(ctx context.Context, token string) (thingId string, err error) {
	claims, err := c.Client.GetClaimsThingToken(ctx, &pb.GetClaimsThingTokenRequest{Jwt: token})
	if err != nil {
		debug.LogError("IsWhichThing GetClaimsThingToken error:", err)
		return "", err
	}
	return claims.GetSubject(), nil
}

// GetThingToken returns a token for a thing given a user id and a thing id
func (c *authClient) GetThingToken(ctx context.Context, thingId string, userId string) (token string, err error) {
	res, err := c.Client.GetThingToken(ctx, &pb.GetThingTokenRequest{ThingId: thingId, UserId: userId})
	if err != nil {
		debug.LogError("GetThingToken error:", err)
		return "", err
	}
	return res.Jwt, nil
}

// GetAdminToken returns a token for an admin
func (c *authClient) GetAdminToken(ctx context.Context) (token string, err error) {
	resToken, err := c.Client.GetAdminToken(ctx, &pb.GetAdminTokenRequest{})
	if err != nil {
		debug.LogError("GRPC get admin token error:", err)
		return "", err
	}
	return resToken.GetJwt(), nil
}

// GetUserToken returns a token for a user
func (c *authClient) GetUserToken(ctx context.Context, userId string) (token string, err error) {
	resToken, err := c.Client.GetUserToken(ctx, &pb.GetUserTokenRequest{UserId: userId})
	if err != nil {
		debug.LogError("GRPC get admin token error:", err)
		return "", err
	}
	return resToken.GetJwt(), nil
}
