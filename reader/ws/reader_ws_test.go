package ws

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

const host = "localhost"
const port = "8080"
const address = host + ":" + port

func TestThingGetEndpoint(t *testing.T) {
	srv := Start(Config{Host: host, Port: port})
	// sendSomething()
	howMany := 6
	rec := receiveUpdateTest(6)
	// log.Println("rec: ", rec)
	assert.Equal(t, howMany, len(rec), "should return "+strconv.Itoa(howMany)+" messages")
	assert.True(t, strings.Contains(rec[0], "thingId"), "should contain 'thingId'")
	srv.Close()
}

func getJsonStartThingUpdate(token string, thingId string) string {
	m := ThingUpdateMsg{}
	m.Token = token
	m.ThingId = thingId
	jsonM, err := json.Marshal(m)
	if err != nil {
		log.Fatal("err json:", err)
	}
	return string(jsonM)
}

func receiveUpdateTest(howMany int) []string {
	var ret []string
	u := url.URL{Scheme: "ws", Host: address, Path: "/thing"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})
	cnt := 0

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}
			ret = append(ret, string(message))
			cnt += 1
			if cnt >= howMany {
				close(done)
				return
			}
		}
	}()

	msg := getJsonStartThingUpdate("123.123.123", "000.000.001")

	err = c.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Println("write message error:", err)
		return nil
	}
	<-done
	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close error:", err)
		return nil
	}
	return ret
}
