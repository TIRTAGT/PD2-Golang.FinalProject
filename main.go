package main

import (
	"fmt"

	"github.com/TIRTAGT/PD2-Golang.FinalProject/server"
)

const (
	LISTEN_ADDRESS = "localhost"
	LISTEN_PORT = 8080
)

func main() {
	fmt.Println("Starting Golang Web Server...")

	server.Start(LISTEN_ADDRESS, LISTEN_PORT)
}