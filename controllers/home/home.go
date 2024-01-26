package home

import (
	"net/http"

	"github.com/TIRTAGT/PD2-Golang.FinalProject/server/controller/utility"
)

func GET(w http.ResponseWriter, r *http.Request) {
	// Beritahu browser kita akan mengirimkan konten HTML
	w.Header().Set("Content-Type", "text/html")

	// Tampilkan view test/index.html
	konten_halaman := utility.FView("/example.html", []utility.VariablePair{
		{ Key: "name", Value: "Katon", },
	})

	// Tulis konten halaman ke browser (sebagai response dari server)
	w.Write([]byte(konten_halaman))
}