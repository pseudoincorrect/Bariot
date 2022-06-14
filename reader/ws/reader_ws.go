package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pseudoincorrect/bariot/reader/service"
)

const closedConn = "wsasend"

var upgrader = websocket.Upgrader{}

type WsServer interface {
	Close()
}

var _ WsServer = (*wsServer)(nil)

type wsServer struct {
	waitGroup *sync.WaitGroup
	server    *http.Server
}

// Close the HTTPserver
func (s *wsServer) Close() {
	if err := s.server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
	s.waitGroup.Wait()
}

type Config struct {
	Host string
	Port string
	S    service.ReaderSvc
}

// Start the configuration of the server
func Start(conf Config) wsServer {
	addr := conf.Host + ":" + conf.Port
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv := StartServer(addr, httpServerExitDone, conf.S)
	return wsServer{waitGroup: httpServerExitDone, server: srv}
}

// StartServer create endpoint and start the HTTP server
func StartServer(addr string, wg *sync.WaitGroup, s service.ReaderSvc) *http.Server {
	server := &http.Server{Addr: addr}
	http.HandleFunc("/thing", getSingleThingHandler(s))
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe() error: %v", err)
		}
	}()
	return server
}

// getSingleThingHandler return a HTTP/WS handler to get a continuous stream of thing data
func getSingleThingHandler(s service.ReaderSvc) http.HandlerFunc {
	singleThingHandler := func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade error:", err)
			return
		}
		defer c.Close()
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			return
		}
		authMsg := ThingUpdateMsg{}
		err = json.Unmarshal(message, &authMsg)
		if err != nil {
			log.Println("err:", err)
		}
		msgJson := ThingUpdateMsg{}
		err = json.Unmarshal(message, &msgJson)
		if err != nil {
			log.Println("err:", err)
			return
		}
		authorized, err := s.AuthorizeSingleThing(msgJson.Token, msgJson.ThingId)
		if err != nil {
			log.Println("err:", err)
		}
		if !authorized {
			err = c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.ClosePolicyViolation, WsUnauthorized))
			if err != nil {
				log.Println("err:", err)
			}
		}
		cnt := 1
		for {
			thingData := "{\"thingId\": \"000.000.001\", \"data\": \"" + strconv.Itoa(cnt) + "\"}"
			err = c.WriteMessage(websocket.TextMessage, []byte(thingData))
			if err != nil {
				if !strings.Contains(err.Error(), closedConn) {
					log.Println("Server Write err:", err)
				}
				return
			}
			cnt += 1
		}
	}
	return singleThingHandler
}
