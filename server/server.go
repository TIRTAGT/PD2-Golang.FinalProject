package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/TIRTAGT/PD2-Golang.FinalProject/server/routing"
)

type RequestHandler struct {
	mu sync.Mutex
}

func (h *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// fmt.Fprintf(w, "count is %d\n", n)
	// w.Header().Add("Content-Type", "text/html")
	w.Header().Add("Server", runtime.Version())

	var REQUEST_URI = r.URL.Path

	// Jika REQUEST_URI memiliki "..", jangan izinkan untuk mengakses.
	if strings.Contains(REQUEST_URI, "..") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 SuS Request"))

		// Print request log yang mirip kaya nginx
		fmt.Printf("%s - [%s] \"%s %s %s\" %d 400_SUS_REQUEST\n", r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.Path, r.Proto, http.StatusBadRequest)
		return
	}

	// Jika mulai dari /assets, berikan static file
	if strings.HasPrefix(REQUEST_URI, "/assets/") {
		var LOKASI_FILE = "." + REQUEST_URI

		// Cek informasi file tersebut di OS
		_, err := os.Stat(LOKASI_FILE);

		// Jika tidak ada, kirim 404
		if errors.Is(err, os.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 File Asset yang diminta tidak ditemukan"))

			// Print request log yang mirip kaya nginx
			fmt.Printf("%s - [%s] \"%s %s %s\" %d FILE_ASSET\n", r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.Path, r.Proto, http.StatusNotFound)
			return
		}

		// Jika ada, kirim file tersebut
		fmt.Printf("%s - [%s] \"%s %s %s\" %d\n", r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.Path, r.Proto, http.StatusOK)
		http.ServeFile(w, r, LOKASI_FILE)
		return
	}
	
	// Kasih requestnya ke server/routes.go
	routing.HandleRoute(w, r)

	// Print request log yang mirip kaya nginx
	fmt.Printf("%s - [%s] \"%s %s %s\" %d %s bytes\n", r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.Path, r.Proto, http.StatusOK, w.Header().Get("Content-Length"))
}

var ServerInstance http.Server

func Start(LISTEN_ADDRESS string, LISTEN_PORT int) {
	fmt.Println("Starting Golang Web Server...")

	var address = LISTEN_ADDRESS + ":" + fmt.Sprintf("%d", LISTEN_PORT)

	ServerInstance = http.Server{
		Addr: address,
		Handler: new(RequestHandler),
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("Golang Web Server akan dimulai pada: http://" + ServerInstance.Addr + "/")
	fmt.Println("Tekan CTRL+C atau kirim SIGTERM/SIGINT untuk menghentikan server.")
	
	// Sifatnya blocking (kodenya akan berhenti sampai servernya mati)
	// TODO: Cari alternatif yang non-blocking
	err := ServerInstance.ListenAndServe()

	if err != nil {
		fmt.Println("Ada kesalahan saat memulai server pada http://" + ServerInstance.Addr + "/")
		fmt.Println(err.Error())
	}
}