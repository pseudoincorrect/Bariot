package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pseudoincorrect/bariot/internal/reader/service"
	e "github.com/pseudoincorrect/bariot/pkg/utils/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
)

const closedConn = "wsasend"

var upgrader = websocket.Upgrader{}

type WsServer interface {
	Close()
}

var _ WsServer = (*wsServer)(nil)

type wsServer struct {
	waitGroup *sync.WaitGroup // wait server conn closing
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
	Host    string
	Port    string
	Service service.Reader
}

// Start the configuration of the server
func Start(conf Config) wsServer {
	addr := ":" + conf.Port
	// addr := conf.Host + ":" + conf.Port
	httpServerExitDone := &sync.WaitGroup{}
	httpServerExitDone.Add(1)
	srv := StartServer(addr, httpServerExitDone, conf.Service)
	return wsServer{waitGroup: httpServerExitDone, server: srv}
}

// StartServer create endpoint and start the HTTP server
func StartServer(addr string, wg *sync.WaitGroup, s service.Reader) *http.Server {
	server := &http.Server{Addr: addr}
	logger.Debug("Reader, start server on", server.Addr)
	http.HandleFunc("/reader/thing", getSingleThingEndpoint(s))
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			e.HandleFatal(e.ErrConn, err, "listen and serve")
		}
	}()
	return server
}

// getSingleThingEndpoint return a HTTP/WS handler to get a continuous stream of thing data
func getSingleThingEndpoint(s service.Reader) http.HandlerFunc {
	singleThingHandler := func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			e.Handle(e.ErrConn, err, "upgrade")
			return
		}
		defer c.Close()
		msgJson, err := authorizeConn(c, s)
		if err != nil {
			return
		}
		thingId := msgJson.ThingId
		stop := make(chan bool)
		thingData := make(chan string)
		// function to be called when data is received
		handler := func(msg string) {
			thingData <- msg
		}
		// handler will be subscribed to thing id subject
		go s.ReceiveThingData(thingId, handler, stop)
		cnt := 1
		for data := range thingData {
			err := sendThingData(c, data)
			if err != nil {
				stop <- true
				e.Handle(e.ErrWrite, err, "send thing data")
			}
			cnt += 1
		}
	}
	return singleThingHandler
}

func authorizeConn(c *websocket.Conn, s service.Reader) (*ThingAuthMsg, error) {
	msgType, message, err := c.ReadMessage()
	logger.Debug("WS message: ", string(message), ", type: ", msgType)
	if err != nil {
		err = e.Handle(e.ErrConn, err, "read message")
		return nil, err
	}
	msgJson, err := decodeAuth(message)
	if err != nil {
		return nil, err
	}
	err = authorizeGetThingData(s, msgJson)
	if err != nil {
		prevErr := err
		err = closeServerConn(c, WsUnauthorized)
		if err != nil {
			return nil, err
		}
		return nil, prevErr
	}
	return msgJson, nil
}

// authorizeGetThingData a user to get a thing's data
func authorizeGetThingData(s service.Reader, msg *ThingAuthMsg) error {
	err := s.AuthorizeSingleThing(msg.Token, msg.ThingId)
	if err != nil {
		return e.Handle(e.ErrAuthz, err, "authorize single thing")
	}
	return nil
}

// decode an json auth message
func decodeAuth(msg []byte) (*ThingAuthMsg, error) {
	authMsg := ThingAuthMsg{}
	err := json.Unmarshal(msg, &authMsg)
	if err != nil {
		err = e.Handle(e.ErrParsing, err, "json unmarshal auth")
		return nil, err
	}
	msgJson := ThingAuthMsg{}
	err = json.Unmarshal(msg, &msgJson)
	if err != nil {
		e.Handle(e.ErrParsing, err, "json unmarshal data")
		return nil, err
	}
	return &msgJson, nil
}

// sendThingData sends mock data over a websocket connection
func sendThingData(c *websocket.Conn, data string) error {
	err := c.WriteMessage(websocket.TextMessage, []byte(data))
	if err != nil {
		if !strings.Contains(err.Error(), closedConn) {
			err = e.Handle(e.ErrConn, err, "Write message")
		}
		return err
	}
	return nil
}

// closeServerConn closes the websocket connection
func closeServerConn(c *websocket.Conn, msg string) error {
	err := c.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.ClosePolicyViolation, msg))
	return err
}
