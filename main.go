package main

import (
	"github.com/TIRTAGT/PD2-Golang.FinalProject/config"
	"github.com/TIRTAGT/PD2-Golang.FinalProject/server"
)

func main() {
	server.Start(config.HTTP_LISTEN_ADDRESS, config.HTTP_LISTEN_PORT)
}
