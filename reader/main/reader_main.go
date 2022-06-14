package main

import (
	"log"

	"github.com/pseudoincorrect/bariot/pkg/env"
	websocket "github.com/pseudoincorrect/bariot/reader/ws"
)

type config struct {
	wsPort string
}

// Load config from environment variables
func loadConfig() config {
	var conf = config{
		wsPort: env.GetEnv("WS_PORT"),
	}
	return conf
}

func main() {
	conf := loadConfig()
	log.Println(conf)
	log.Println("Reader service online")
	websocket.Start(websocket.Config{Port: "8080"})
}
