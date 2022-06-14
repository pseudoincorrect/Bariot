package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	mock "github.com/pseudoincorrect/bariot/reader/mock/service"
	"github.com/stretchr/testify/assert"
)

var registerOnce sync.Once

const host = "localhost"
const port = "8080"
const address = host + ":" + port

func getMock() *mock.MockReader {
	return new(mock.MockReader)
}

func TestThingGetEndpointAuthSuccess(t *testing.T) {
	resetHttpHandler()
	theMock := getMock()
	theMock.On("AuthorizeSingleThing").Return(nil, true).Once()
	conf := Config{Host: host, Port: port, S: theMock}
	// srv := startOnce(conf)
	srv := Start(conf)
	howMany := 6
	rec, _ := receiveUpdate(6)
	assert.Equal(t, howMany, len(rec), "should return "+strconv.Itoa(howMany)+" messages")
	assert.True(t, strings.Contains(rec[0], "thingId"), "should contain 'thingId'")
	srv.Close()
}

func TestThingGetEndpointAuthFail(t *testing.T) {
	resetHttpHandler()
	theMock := getMock()
	theMock.On("AuthorizeSingleThing").Return(nil, false).Once()
	conf := Config{Host: host, Port: port, S: theMock}
	// srv := startOnce(conf)
	srv := Start(conf)
	rec, err := receiveUpdate(6)
	assert.Nil(t, rec, "should get no message")
	assert.True(t, strings.Contains(err.Error(), WsUnauthorized), "should contain 'unauthorized'")
	srv.Close()
}

// func startOnce(c Config) wsServer {
// 	var srv wsServer
// 	registerOnce.Do(func() {
// 		srv = Start(c)
// 	})
// 	return srv
// }

func receiveUpdate(howManyMsg int) ([]string, error) {
	var retMsg []string
	var retErr error
	u := url.URL{Scheme: "ws", Host: address, Path: "/thing"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})
	msgCnt := 0
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				retErr = err
				return
			}
			retMsg = append(retMsg, string(message))
			msgCnt += 1
			if msgCnt >= howManyMsg {
				retErr = nil
				return
			}
		}
	}()
	msg := getJsonAuthMsg("123.123.123", "000.000.001")
	err = c.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write message error:", err)
		return nil, nil
	}

	<-done
	if retErr != nil {
		return nil, retErr
	}
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close error:", err)
		return nil, nil
	}
	return retMsg, nil
}

func getJsonAuthMsg(token string, thingId string) string {
	m := ThingUpdateMsg{}
	m.Token = token
	m.ThingId = thingId
	jsonM, err := json.Marshal(m)
	if err != nil {
		log.Fatal("err json:", err)
	}
	return string(jsonM)
}

// resetHttpHandler  reset http routes
// http is global registering the same route twice (for unit test)
// cause an error, we need to unregister the routes first
func resetHttpHandler() {
	http.DefaultServeMux = new(http.ServeMux)
}
