package home

import (
	"net/http"

	"github.com/TIRTAGT/PD2-Golang.FinalProject/server/controller/utility"
)

func GET(w http.ResponseWriter, r *http.Request) *string {
	// Siapkan view
	konten_halaman := utility.FView("/index.html", []utility.VariablePair{
		{Key: "name", Value: "Katon"},
	})

	// Kirim konten halaman ke browser (sebagai response dari server)
	return &konten_halaman
}
