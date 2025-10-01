package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tictactoe-ssh/pkg/server"
)

func main() {
	port := os.Getenv("SSH_PORT")
	if port == "" {
		port = "2222"
	}

	fmt.Printf("Starting Tic-Tac-Toe SSH Server on port %s...\n", port)
	fmt.Println("Connect with: ssh -p", port, "-o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null localhost")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := server.Start(port); err != nil {
			log.Fatal(err)
		}
	}()

	// Wait for interrupt signal
	<-sigChan
	fmt.Println("\nReceived interrupt signal. Shutting down gracefully...")
}
