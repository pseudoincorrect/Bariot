package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	e "github.com/pseudoincorrect/bariot/pkg/utils/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
	"github.com/pseudoincorrect/bariot/tests/mocks/services"
	"github.com/stretchr/testify/assert"
)

const host = "localhost"
const port = "80"
const address = host + ":" + port

func TestThingGetEndpointAuthFail(t *testing.T) {
	resetHttpHandler()
	theMock := services.NewMockReader()
	theMock.On("AuthorizeSingleThing").Return(e.ErrAuthz).Once()
	conf := Config{Host: host, Port: port, Service: &theMock}
	srv := Start(conf)
	rec, err := receiveUpdate(6)
	assert.Nil(t, rec, "should get no message")
	assert.True(t, strings.Contains(err.Error(), WsUnauthorized), "should contain 'unauthorized'")
	srv.Close()
}

func TestThingGetEndpointAuthSuccess(t *testing.T) {
	resetHttpHandler()
	theMock := services.NewMockReader()
	theMock.On("AuthorizeSingleThing").Return(nil).Once()
	conf := Config{Host: host, Port: port, Service: &theMock}
	srv := Start(conf)
	howMany := 6
	rec, _ := receiveUpdate(6)
	assert.Equal(t, howMany, len(rec), "should return "+strconv.Itoa(howMany)+" messages")
	assert.True(t, strings.Contains(rec[0], "thingId"), "should contain 'thingId'")
	srv.Close()
}

func receiveUpdate(howManyMsg int) ([]string, error) {
	var msgStore []string
	var sendMsgErr error
	u := url.URL{Scheme: "ws", Host: address, Path: "reader/thing"}
	logger.Debug("WS URL:", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		e.HandleFatal(e.ErrConn, err, "websocket dial")
	}
	defer c.Close()
	done := make(chan struct{})
	go readMsgs(c, howManyMsg, &msgStore, &sendMsgErr, done)
	err = authorize(c)
	if err != nil {
		return nil, err
	}
	<-done
	if sendMsgErr != nil {
		return nil, sendMsgErr
	}
	err = closeClientConn(c)
	if err != nil {
		return nil, err
	}
	return msgStore, nil
}

func readMsgs(c *websocket.Conn, howManyMsg int, msgStore *[]string, ret *error, done chan struct{}) {
	defer close(done)
	msgCnt := 0
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			e.Handle(e.ErrRead, err, "client read message")
			*ret = err
			return
		}
		*msgStore = append(*msgStore, string(message))
		msgCnt += 1
		if msgCnt >= howManyMsg {
			*ret = nil
			return
		}
	}
}

func getJsonAuthMsg(token string, thingId string) string {
	m := ThingAuthMsg{}
	m.Token = token
	m.ThingId = thingId
	jsonM, err := json.Marshal(m)
	if err != nil {
		log.Fatal("err json:", err)
	}
	return string(jsonM)
}

func authorize(c *websocket.Conn) error {
	msg := getJsonAuthMsg("123.123.123", "000.000.001")
	err := writeMessage(c, msg)
	if err != nil {
		return err
	}
	return nil
}

func writeMessage(c *websocket.Conn, msg string) error {
	err := c.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		err = e.Handle(e.ErrWrite, err, "client write message")
		return err
	}
	return nil
}

func closeClientConn(c *websocket.Conn) error {
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		err = e.Handle(e.ErrWrite, err, "client write close message")
		return err
	}
	return nil
}

// resetHttpHandler  reset http routes
// http is global registering the same route twice (for unit test)
// cause an error, we need to unregister the routes first
func resetHttpHandler() {
	http.DefaultServeMux = new(http.ServeMux)
}
