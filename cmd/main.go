package main

import (
	"gobank/internal/api"
	"log"
)

// TODO: move to configuration

const (
	ServerAddress = "localhost:8080"
)

func main() {
	s := api.NewServer()
	err := s.Run(ServerAddress)
	if err != nil {
		log.Fatal("cannot start api:", err)
	}
}
