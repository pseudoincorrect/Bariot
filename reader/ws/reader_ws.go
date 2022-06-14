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

func (s *wsServer) Close() {
	if err := s.server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
	s.waitGroup.Wait()
}

type Config struct {
	Host string
	Port string
}

func Start(conf Config) wsServer {
	addr := conf.Host + ":" + conf.Port
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv := StartServer(addr, httpServerExitDone)
	return wsServer{waitGroup: httpServerExitDone, server: srv}
}

func StartServer(addr string, wg *sync.WaitGroup) *http.Server {
	srv := &http.Server{Addr: addr}
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/thing", singleThingHandler)

	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe() error: %v", err)
		}
	}()

	return srv
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}
		log.Printf("server recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}

func singleThingHandler(w http.ResponseWriter, r *http.Request) {
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
	// log.Printf("server recv: %s", message)
	msgJson := ThingUpdateMsg{}
	err = json.Unmarshal(message, &msgJson)
	if err != nil {
		log.Println("err:", err)
		return
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
