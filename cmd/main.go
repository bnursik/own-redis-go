package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"own-redis/internal/server"
	"own-redis/internal/storage"
)

//TODO gofmt

const (
	defaultPort = 8080
	usageText   = `Own Redis

Usage:
  own-redis [--port <N>]
  own-redis --help

Options:
  --help       Show this screen.
  --port N     Port number.`
)

func main() {
	port := flag.Int("port", defaultPort, "Port number to listen on")
	help := flag.Bool("help", false, "Show usage information")
	flag.Parse()

	if *help {
		fmt.Println(usageText)
		os.Exit(0)
	}

	// Create key-value store
	store := storage.NewStore()

	// Start UDP server
	srv := server.NewUDPServer(*port, store)
	log.Printf("Starting UDP server on port %d", *port)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
