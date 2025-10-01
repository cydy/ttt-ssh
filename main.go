package main

import (
	"fmt"
	"log"
	"os"

	"tictactoe-ssh/pkg/server"
)

func main() {
	port := os.Getenv("SSH_PORT")
	if port == "" {
		port = "2222"
	}

	fmt.Printf("Starting Tic-Tac-Toe SSH Server on port %s...\n", port)
	fmt.Println("Connect with: ssh -p", port, "localhost")

	if err := server.Start(port); err != nil {
		log.Fatal(err)
	}
}
