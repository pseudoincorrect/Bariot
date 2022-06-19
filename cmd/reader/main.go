package main

import (
	"log"

	websocket "github.com/pseudoincorrect/bariot/internal/reader/ws"
	"github.com/pseudoincorrect/bariot/pkg/env"
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
