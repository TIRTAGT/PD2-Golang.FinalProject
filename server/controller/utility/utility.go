package utility

import (
	"errors"
	"fmt"
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
		fmt.Println("Tidak dapat memuat view '" + view_name + "'", err)
		return ""
	}

	return hasil
}

// Gabungkan view dengan template
func render_template(view_name string, variabel_tambahan []VariablePair, konten_tambahan []VariablePair) (string, error) {
	// Coba load view
	hasil_view, err := View(view_name, variabel_tambahan)

	// Jika gagal, berikan error ke pemanggil fungsi
	if err != nil {
		return "view", err
	}

	// Jika ada <body> di view, hapus saja karena content view harusnya berisi elemen <body> saja
	hasil_view = strings.ReplaceAll(hasil_view, "<body>", "")
	hasil_view = strings.ReplaceAll(hasil_view, "</body>", "")

	// Jika tidak ada konten tambahan, buat variabel kosong
	if konten_tambahan == nil || len(konten_tambahan) == 0 {
		konten_tambahan = []VariablePair{}
	}
	
	IsViewInjected := false
	// Cari apakah sudah ada body_content di konten tambahan, hapus(timpah) dengan content view jika ada
	for i, v := range konten_tambahan {
		if v.Key == "body_content" {
			konten_tambahan[i].Value = hasil_view
			IsViewInjected = true
			break
		}
	}

	// Jika tidak ada body_content di konten tambahan, tambahkan content view ke konten tambahan
	if !IsViewInjected {
		konten_tambahan = append(konten_tambahan, VariablePair{"content", hasil_view})
	}

	// Coba load template
	hasil_template, err := View("/template.html", konten_tambahan)

	// Jika gagal, berikan error ke pemanggil fungsi
	if err != nil {
		return "template", err
	}

	// Jika berhasil, kembalikan hasilnya
	return hasil_template, nil
}

// Gabungkan view dengan template, return string kosong jika gagal (Force View)
func Frender_template(view_name string, variabel_tambahan []VariablePair, konten_tambahan []VariablePair) string {
	hasil, err := render_template(view_name, variabel_tambahan, konten_tambahan)

	if err != nil {
		if (hasil == "view") {
			fmt.Println("Tidak dapat memuat view '" + view_name + "'", err)
		} else if (hasil == "template") {
			fmt.Println("Tidak dapat memuat template.html untuk '" + view_name + "'", err)
		}

		return ""
	}

	return hasil
}