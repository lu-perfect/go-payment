package main

import (
	"gobank/internal/api"
	"log"
)

func main() {
	s := api.NewServer()
	err := s.Run()
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
