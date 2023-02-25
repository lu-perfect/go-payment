package main

import (
	"gobank/internal/api"
	"gobank/internal/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	api.NewServer(config).Run()
}
