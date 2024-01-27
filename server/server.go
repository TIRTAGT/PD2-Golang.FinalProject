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

func (h *RequestHandler) ServeHTTP(w *Response, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	w.Header().Add("Server", runtime.Version())

	var REQUEST_URI = r.URL.Path

	// Pastikan REQUEST_URI tidak kosong
	if len(REQUEST_URI) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Null Request"))
		return
	}

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
		_, err := os.Stat(LOKASI_FILE)

		// Jika tidak ada, kirim 404
		if errors.Is(err, os.ErrNotExist) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 File Asset yang diminta tidak ditemukan"))

			// Print request log yang mirip kaya nginx
			fmt.Printf("%s - [%s] \"%s %s %s\" %d FILE_ASSET\n", r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.Path, r.Proto, http.StatusNotFound)
			return
		}

		// Print request log yang mirip kaya nginx
		fmt.Printf("%s - [%s] \"%s %s %s\" %d\n", r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.Path, r.Proto, http.StatusOK)

		// Jika ada, kirim file tersebut
		http.ServeFile(w, r, LOKASI_FILE)
		return
	}

	// Kasih requestnya ke server/routes.go
	routing.HandleRoute(w, r)

	// Jika ditemuakn header Content-Length, print request log beserta ukuran konten yang dikirim
	cl := w.Header().Get("Content-Length")

	if len(cl) > 0 {
		fmt.Printf("%s - [%s] \"%s %s %s\" %d %s bytes\n", r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.Path, r.Proto, w.statusCode, cl)
	} else {
		fmt.Printf("%s - [%s] \"%s %s %s\" %d\n", r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.Method, r.URL.Path, r.Proto, w.statusCode)
	}

}

type Response struct {
	http.ResponseWriter
	statusCode int
}

func (res *Response) WriteHeader(code int) {
	res.statusCode = code
	res.ResponseWriter.WriteHeader(code)
}

var ServerInstance http.Server
var IsServerRunning bool = false
var IzinkanKoneksi bool = false

func Start(LISTEN_ADDRESS string, LISTEN_PORT int, WaitGroup *sync.WaitGroup) {
	fmt.Println("Starting Golang Web Server...")

	var address = LISTEN_ADDRESS + ":" + fmt.Sprintf("%d", LISTEN_PORT)

	// Buat request handler baru
	var HandlerAsli = new(RequestHandler)

	// Tambahkan ResponseWriter custom agar bisa menggunakan fitur tambahan
	var HandlerCustom = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		// Jika tidak mengizinkan koneksi, kirim 503 Service Unavailable
		if !IzinkanKoneksi {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("503 Layanan Web Server tidak tersedia untuk saat ini"))
			return
		}

		HandlerAsli.ServeHTTP(&Response{w, http.StatusOK}, req)
	})

	ServerInstance = http.Server{
		Addr:           address,
		Handler:        HandlerCustom,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// Start server di background (latar belakang)
	// @source https://stackoverflow.com/a/42533360

	go func() {
		defer WaitGroup.Done()

		IsServerRunning = true
		err := ServerInstance.ListenAndServe()
		IsServerRunning = false

		if err != nil && err != http.ErrServerClosed {
			fmt.Println("Ada kesalahan saat memulai server pada " + ServerInstance.Addr + "/")
			fmt.Println(err.Error())
		}
	}()
}
