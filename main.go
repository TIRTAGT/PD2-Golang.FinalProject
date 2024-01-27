package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/TIRTAGT/PD2-Golang.FinalProject/config"
	"github.com/TIRTAGT/PD2-Golang.FinalProject/server"
)

var ServerSuksesStart = false

func main() {
	HTTPWaitGroup := &sync.WaitGroup{}
	HTTPWaitGroup.Add(1)

	// Mulai server
	server.Start(config.HTTP_LISTEN_ADDRESS, config.HTTP_LISTEN_PORT, HTTPWaitGroup)

	// Berikan jeda 2 detik untuk starting server (menghindari Race Condition)
	time.Sleep(2 * time.Second)

	// Tahan CTRL+C / SIGINT untuk mematikan server
	// @source https://stackoverflow.com/a/44598343
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			fmt.Println("\nSedang mematikan Golang Web Server...")

			if sig == os.Interrupt {
				server.ServerInstance.Shutdown(context.TODO())
				break
			}
		}
	}()

	// Cek apakah server sudah berjalan
	if server.IsServerRunning {
		// Cek apakah listen IP adalah 0.0.0.0, ganti dengan anyhost sebagai display
		if config.HTTP_LISTEN_ADDRESS == "0.0.0.0" {
			fmt.Println("Golang Web Server dimulai di semua interface, pada port " + fmt.Sprintf("%d", config.HTTP_LISTEN_PORT) + "")
			fmt.Println("⚠️  Tergantung konfigurasi Firewall anda, Golang Web Server mungkin dapat diakses dari jaringan publik.")
		} else {
			fmt.Println("Golang Web Server dimulai pada: http://" + server.ServerInstance.Addr + "/")
		}

		fmt.Println("Tekan CTRL+C atau kirim SIGTERM/SIGINT untuk menghentikan server.")
		ServerSuksesStart = true

		// Izinkan server untuk menerima koneksi browser client
		server.IzinkanKoneksi = true
	}

	// Tunggu selama server masih berjalan
	HTTPWaitGroup.Wait()

	if ServerSuksesStart {
		fmt.Println("Golang Web Server telah dimatikan.")
	}
}
