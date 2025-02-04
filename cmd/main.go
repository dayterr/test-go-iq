package main

import (
	"log"

	"github.com/dayterr/test-go-iq/internal/config"
	"github.com/dayterr/test-go-iq/internal/server"
	"github.com/dayterr/test-go-iq/internal/storage"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal("no config, can't start the program")
	}

	h := server.NewHandler()
	s := storage.NewStorage(config.DatabaseURI)
	h.Storage = s
	r := server.CreateRouter(h)
	r.Run(":8080")
}
