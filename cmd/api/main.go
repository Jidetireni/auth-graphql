package main

import (
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Application error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	// Initialize the server
	server, err := NewServer()
	if err != nil {
		return err
	}

	if err := server.Initialize(); err != nil {
		return err
	}

	if err := server.Start(); err != nil {
		return err
	}
	return nil
}
