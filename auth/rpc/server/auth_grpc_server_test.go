package server

import (
	"context"
	"log"
	"os"
	"testing"

	mockAuth "github.com/pseudoincorrect/bariot/auth/mock"
	pb "github.com/pseudoincorrect/bariot/pkg/auth/grpc"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port string = "50051"
const host string = "localhost"

var mAuth *mockAuth.Mock
var conn *grpc.ClientConn
var client pb.AuthClient

func connect() {
	addr := host + ":" + port
	log.Println("init user service GRPC client to ", addr)
	c, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("did not connect:", err)
	}
	conn = c
	client = pb.NewAuthClient(conn)
}

// On windows, the firewall will issue a warning
// Defender -> create inbound rule -> port -> 50051
func TestMain(m *testing.M) {
	mAuth = new(mockAuth.Mock)
	conf := ServerConf{
		AuthService: mAuth,
		Port:        port,
	}
	go func() {
		s, err := Start(conf)
		if err != nil {
			s.Stop()
			log.Fatal("Could not start grpc server")
		}
	}()
	connect()
	code := m.Run()
	conn.Close()
	os.Exit(code)
}

func TestGetAdminToken(t *testing.T) {
	ctx := context.Background()
	mockToken := "123.123.123"
	mAuth.On("GetAdminToken").Return(nil, mockToken)
	token, err := client.GetAdminToken(ctx, &pb.GetAdminTokenRequest{})
	assert.Nil(t, err, "should get Admin token without error")
	assert.Equal(t, token.Jwt, mockToken, "token should be", mockToken)
}

func TestGetUserToken(t *testing.T) {
	ctx := context.Background()
	mockToken := "123.123.123"
	userId := "000.000.001"
	mAuth.On("GetUserToken", userId).Return(nil, mockToken)
	req := new(pb.GetUserTokenRequest)
	req.UserId = userId
	token, err := client.GetUserToken(ctx, req)
	assert.Nil(t, err, "should get User token without error")
	assert.Equal(t, token.Jwt, mockToken, "token should be", mockToken)
}

func TestGetThingToken(t *testing.T) {
	ctx := context.Background()
	mockToken := "123.123.123"
	userId := "000.000.002"
	thingId := "000.000.001"
	mAuth.On("GetThingToken", thingId, userId).Return(nil, mockToken)
	req := new(pb.GetThingTokenRequest)
	req.UserId = userId
	req.ThingId = thingId
	token, err := client.GetThingToken(ctx, req)
	assert.Nil(t, err, "should get User token without error")
	assert.Equal(t, token.Jwt, mockToken, "token should be", mockToken)
}

func TestGetClaimsUserToken(t *testing.T) {
	ctx := context.Background()
	mockToken := "123.123.123"
	role := "admin"
	mAuth.On("GetClaimsUserToken", mockToken).Return(nil, role, "001", 24)
	req := new(pb.GetClaimsUserTokenRequest)
	req.Jwt = mockToken
	claims, err := client.GetClaimsUserToken(ctx, req)
	assert.Nil(t, err, "should get User token without error")
	assert.Equal(t, claims.Role, role, "token should be", mockToken)
}

func TestGetClaimsThingToken(t *testing.T) {
	ctx := context.Background()
	mockToken := "123.123.123"
	userId := "000.000.001"
	thingId := "000.000.002"
	mAuth.On("GetClaimsThingToken", mockToken).Return(nil, userId, thingId, 24)
	req := new(pb.GetClaimsThingTokenRequest)
	req.Jwt = mockToken
	claims, err := client.GetClaimsThingToken(ctx, req)
	assert.Nil(t, err, "should get User token without error")
	assert.Equal(t, claims.UserId, userId, "token should be", mockToken)
}
