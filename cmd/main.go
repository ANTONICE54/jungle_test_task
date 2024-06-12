package main

import (
	"app/internal/server"
	"log"
)

func main() {
	server, err := server.NewServer()
	if err != nil {
		log.Fatal("Cannot create the server: ", err)
	}

	err = server.Start()
	if err != nil {
		log.Fatal("Cannot start the server:", err)
	}
}
