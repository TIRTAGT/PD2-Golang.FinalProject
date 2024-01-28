package submit_form

import (
	"crypto/tls"
	"fmt"
	"html"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/TIRTAGT/PD2-Golang.FinalProject/config"
	"github.com/TIRTAGT/PD2-Golang.FinalProject/server/controller/utility"
)

func GET(w http.ResponseWriter, r *http.Request) *string {
	// Redirect ke halaman utama
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return nil
}

func POST(w http.ResponseWriter, r *http.Request) *string {
	// Pastikan ada data yang dikirim
	if err := r.ParseForm(); err != nil {
		// Jika tidak ada, redirect ke halaman utama
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return nil
	}

	// Ambil data dari form
	Form_name := r.Form.Get("name")
	Form_email := r.PostForm.Get("email")
	Form_subject := r.PostForm.Get("subject")
	Form_message := r.PostForm.Get("message")

	// Pastikan panjang nama lebih dari 0
	if len(Form_name) == 0 {
		pesan_error := "Nama tidak boleh kosong"
		return &pesan_error
	}

	// Pastikan panjang email lebih dari 0, dan memiliki format email yang benar
	if !IsEmailValid(Form_email) {
		pesan_error := "Email tidak valid"
		return &pesan_error
	}

	// Pastikan panjang subject lebih dari 0
	if len(Form_subject) == 0 {
		pesan_error := "Subject tidak boleh kosong"
		return &pesan_error
	}

	// Pastikan panjang message lebih dari 0
	if len(Form_message) == 0 {
		pesan_error := "Message tidak boleh kosong"
		return &pesan_error
	}

	// Escape semua karakter untuk menghindari XSS attack
	Form_name = html.EscapeString(Form_name)
	Form_email = html.EscapeString(Form_email)
	Form_subject = html.EscapeString(Form_subject)
	Form_message = html.EscapeString(Form_message)

	// #region Koneksi ke Server SMTP
	// @docs: https://golang.org/pkg/net/smtp/
	// @reference: https://gist.github.com/jpillora/cb46d183eca0710d909a
	// @reference: https://gist.github.com/muskankhedia/084243aceb55fdedf5c834691cefa3c1
	// @reference: https://stackoverflow.com/a/9951508/12718814

	TLSConn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", config.MAIL_SMTP_HOSTNAME, config.MAIL_SMTP_PORT), &tls.Config{
		ServerName: config.MAIL_SMTP_HOSTNAME,
	})

	if err != nil {
		fmt.Println("Error: Tidak dapat melakukan koneksi ke server SMTP, error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	SMTPSession, err := smtp.NewClient(TLSConn, config.MAIL_SMTP_HOSTNAME)

	if err != nil {
		SMTPSession.Close()

		fmt.Println("Error: Tidak dapat memulai sesi SMTP pada server, error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	StatusTLS, IsTLSOk := SMTPSession.TLSConnectionState()
	if !IsTLSOk {
		SMTPSession.Close()

		fmt.Println("Error: Koneksi dengan enkripsi TLS tidak dapat dilakukan ke server SMTP")
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	// Batalkan koneksi jika menggunakan enkripsi yang tidak aman sesuai standar sekarang (SSLv3, TLSv1.0, TLSv1.1)
	if StatusTLS.Version == tls.VersionSSL30 || StatusTLS.Version == tls.VersionTLS10 || StatusTLS.Version == tls.VersionTLS11 {
		SMTPSession.Close()

		fmt.Printf("Error: Versi TLS yang digunakan tidak sesuai (%x)\n", StatusTLS.Version)
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	// Login tanpa enkripsi tambahan / PLAIN (koneksi sudah dienkripsi oleh TLS)
	MailServerLogin := smtp.PlainAuth("", config.MAIL_SMTP_USERNAME, config.MAIL_SMTP_PASSWORD, config.MAIL_SMTP_HOSTNAME)

	if err = SMTPSession.Auth(MailServerLogin); err != nil {
		SMTPSession.Close()

		fmt.Println("Error: Tidak dapat melakukan login ke server SMTP", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	// Set pengirim
	if err = SMTPSession.Mail(config.MAIL_SMTP_SENDER_ADDRESS); err != nil {
		SMTPSession.Close()

		fmt.Println("Error: Tidak dapat mengatur pengirim email", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	// Set penerima
	if err = SMTPSession.Rcpt(config.MAIL_SMTP_RECEIVER_ADDRESS); err != nil {
		SMTPSession.Close()

		fmt.Println("Error: Tidak dapat mengatur penerima email", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}
	// #endregion

	// Siapkan isi email berdasarkan template
	konten_halaman := utility.FView("/submit_form/mail_template.html", []utility.VariablePair{
		{Key: "name", Value: Form_name},
		{Key: "email", Value: Form_email},
		{Key: "subject", Value: Form_subject},
		{Key: "message", Value: Form_message},
		{Key: "page_title", Value: "Web MK"},
	})

	// Kirim email
	EmailWriter, err := SMTPSession.Data()
	if err != nil {
		SMTPSession.Close()

		fmt.Println("Error: Tidak dapat menulis email", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	EmailWriter.Write([]byte(fmt.Sprintf("From: \"%s\" <%s>\n", config.MAIL_SMTP_SENDER_NAME, config.MAIL_SMTP_SENDER_ADDRESS)))
	EmailWriter.Write([]byte(fmt.Sprintf("To: \"%s\" <%s>\n", config.MAIL_SMTP_RECEIVER_NAME, config.MAIL_SMTP_RECEIVER_ADDRESS)))
	EmailWriter.Write([]byte(fmt.Sprintf("Subject: %s\n", config.MAIL_SUBJECT)))
	EmailWriter.Write([]byte("MIME-version: 1.0;\n"))
	EmailWriter.Write([]byte("Content-Type: text/html; charset=\"UTF-8\";"))
	EmailWriter.Write([]byte("\n\n"))
	EmailWriter.Write([]byte(konten_halaman))

	// Tandai email sudah selesai ditulis
	if err = EmailWriter.Close(); err != nil {
		SMTPSession.Close()

		fmt.Println("Error: Tidak dapat menyelesaikan penulisan email", err)
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	// Tutup koneksi
	SMTPSession.Quit()

	fmt.Println("Email berhasil dikirim ke", config.MAIL_SMTP_RECEIVER_ADDRESS)

	// Kirim OK ke client, validate.js pada sisi client akan mengubahnya menjadi alert success
	pesan_berhasil := "OK"
	return &pesan_berhasil
}

func IsEmailValid(email string) bool {
	// Pastikan email tidak kosong
	if len(email) == 0 {
		return false
	}

	// Email tidak boleh ada spasi
	if strings.Contains(email, " ") {
		return false
	}

	// Email wajib ada 1 "@" saja
	if strings.Count(email, "@") > 1 {
		return false
	}

	// Email wajib ada 1 atau lebih "."
	if strings.Count(email, ".") < 1 {
		return false
	}

	return true
}
