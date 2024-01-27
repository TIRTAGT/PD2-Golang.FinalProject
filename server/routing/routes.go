package routing

import (
	"fmt"
	"net/http"

	"github.com/TIRTAGT/PD2-Golang.FinalProject/controllers"
	_ "github.com/TIRTAGT/PD2-Golang.FinalProject/controllers"
)

func HandleRoute(w http.ResponseWriter, r *http.Request) {
	var REQUEST_URI = r.URL.Path

	// Hapus tanda / yang tidak dilanjutkan dengan apapun pada akhir URI, kecuali jika URI hanya / saja
	if REQUEST_URI != "/" && REQUEST_URI[len(REQUEST_URI)-1] == '/' {
		REQUEST_URI = REQUEST_URI[:len(REQUEST_URI)-1]
	}

	// Cek jika ada controller yang dapat handle request URI tersebut
	// RUMUS: /controller/<URI>/<AKHIRAN_URI>.go
	// Contoh: /test
	// Maka: /controller/test/test.go
	ControllerHandler, IsControllerExist := controllers.ControllerMap[REQUEST_URI]

	// Jika ada controller yang terdaftar untuk menghandle request URI tersebut
	if IsControllerExist {
		IsHandledByController := false
		IsMethodNotSupported := false

		var HasilController *string

		// Panggil request handler pada controller sesuai dengan metode request
		switch r.Method {
		case "GET":
			if ControllerHandler.GET != nil {
				HasilController = ControllerHandler.GET(w, r)
				IsHandledByController = true
			} else {
				IsMethodNotSupported = true
			}

		case "POST":
			if ControllerHandler.POST != nil {
				HasilController = ControllerHandler.POST(w, r)
				IsHandledByController = true
			} else {
				IsMethodNotSupported = true
			}

		case "PUT":
			if ControllerHandler.PUT != nil {
				HasilController = ControllerHandler.PUT(w, r)
				IsHandledByController = true
			} else {
				IsMethodNotSupported = true
			}

		case "DELETE":
			if ControllerHandler.DELETE != nil {
				HasilController = ControllerHandler.DELETE(w, r)
				IsHandledByController = true
			} else {
				IsMethodNotSupported = true
			}
		}

		if IsMethodNotSupported {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("405 Halaman yang dituju tidak mendukung metode request anda (" + r.Method + ")"))
			return
		}

		if !IsHandledByController {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Terjadi kesalahan pada server"))
			return
		}

		// Jika controller mengembalikan sesuatu, kirimkan ke browser
		if HasilController != nil && len(*HasilController) > 0 {
			var tmp = []byte(*HasilController)

			// Secara default, kirim konten sebagai HTML
			w.Header().Add("Content-Type", "text/html")
			w.Header().Add("Content-Length", fmt.Sprint(len(tmp)))
			w.WriteHeader(http.StatusOK)
			w.Write(tmp)
		}

		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Tidak ada controller yang dapat melayani request ke " + REQUEST_URI))
}
