package server

import (
	"context"
	"log"
	"os"
	"testing"

	pb "github.com/pseudoincorrect/bariot/pkg/things/grpc"
	"github.com/pseudoincorrect/bariot/pkg/utils/debug"
	"github.com/pseudoincorrect/bariot/tests/mocks/services"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const port string = "50051"
const host string = "localhost"

var mThings services.ThingsMock
var conn *grpc.ClientConn
var client pb.ThingsClient

func connect() pb.ThingsClient {
	addr := host + ":" + port
	debug.LogInfo("init user service GRPC client to ", addr)
	c, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("did not connect:", err)
	}
	conn = c
	return pb.NewThingsClient(conn)
}

// On windows, the firewall will issue a warning
// Defender -> create inbound rule -> port -> 50051
func TestMain(m *testing.M) {
	mThings = services.NewThingsMock()
	conf := ServerConf{
		Service: &mThings,
		Port:    port,
	}
	go func() {
		s, err := Start(conf)
		if err != nil {
			s.Stop()
			log.Fatal("Could not start grpc server")
		}
	}()
	client = connect()
	code := m.Run()
	conn.Close()
	os.Exit(code)
}

func TestGetUserOfThing(t *testing.T) {
	ctx := context.Background()
	thingId := "000.000.001"
	mThings.UserId = "000.000.002"
	res, err := client.GetUserOfThing(ctx, &pb.GetUserOfThingRequest{ThingId: thingId})
	assert.Nil(t, err, "should get Admin token without error")
	assert.Equal(t, mThings.UserId, res.UserId, "Should get the right user id")
}
