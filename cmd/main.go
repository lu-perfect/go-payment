package main

import (
	"gobank/internal/server"
	"log"
)

// TODO: move to configuration

const (
	ServerAddress = "localhost:8080"
)

func main() {
	s := server.New()
	err := s.Run(ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
