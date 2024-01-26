package utility

import (
	"errors"
	"os"
	"strings"
)

const (
	IS_VIEW_DIR = "404 Halaman yang diminta tidak ditemukan"
	IS_VIEW_NOT_FOUND = "404 Halaman yang diminta tidak ditemukan"
	IS_VIEW_NOT_READABLE = "500 Halaman yang diminta tidak dapat dibaca"
	ILLEGAL_VIEW_PATH = "500 Path view tidak valid"
)

type VariablePair struct {
	Key string
	Value string
}

// Siapkan view, return error jika gagal
func View(view_name string, variabel []VariablePair) (string, error) {
	// Jika view_name diawali dengan "/", hapus saja
	if view_name[0] == '/' {
		view_name = view_name[1:]
	}

	// Jika ada .. di view name, jangan izinkan
	if strings.Contains(view_name, "..") {
		return "", errors.New(ILLEGAL_VIEW_PATH)
	}

	// Buat tuple untuk hasil
	var hasil string = ""


	// Periksa apakah ada file view yang diminta pada folder views
	InfoView, err := os.Stat("./views/" + view_name)

	if errors.Is(err, os.ErrNotExist) {
		return hasil, errors.New(IS_VIEW_NOT_FOUND)
	}

	if InfoView.IsDir() {
		return hasil, errors.New(IS_VIEW_DIR)
	}

	// Baca filenya pakai ioutil
	ViewContent, err := os.ReadFile("./views/" + view_name)

	if err != nil {
		return hasil, errors.New(IS_VIEW_NOT_READABLE)
	}

	// Ubah isi file menjadi string
	hasil = string(ViewContent)

	// Jika ada variabel yang diberikan, ganti semua template variabel yang ada di view
	if len(variabel) > 0 {
		for _, v := range variabel {
			hasil = strings.ReplaceAll(hasil, "{{" + v.Key + "}}", v.Value)
			hasil = strings.ReplaceAll(hasil, "{{ " + v.Key + " }}", v.Value)
		}
	}

	// Jika berhasil, kembalikan hasilnya
	return hasil, nil
}

// Siapkan view, return string kosong jika gagal (Force View)
func FView(view_name string, variabel []VariablePair) string {
	hasil, err := View(view_name, variabel)

	if err != nil {
		return ""
	}

	return hasil
}