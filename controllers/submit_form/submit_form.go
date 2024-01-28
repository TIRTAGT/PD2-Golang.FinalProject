package submit_form

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/smtp"

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

	// #region Koneksi ke Server SMTP
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

	// Ambil data dari form
	// FIXME: VALIDASI DATA UNTUK MENGHINDARI XSS ATTACK
	name := r.PostForm.Get("name")

	// Siapkan isi email berdasarkan template
	konten_halaman := utility.FView("/submit_form/mail_template.html", []utility.VariablePair{
		{Key: "name", Value: name},
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

	// Kirim konten halaman ke browser (sebagai response dari server)
	pesan_berhasil := "Email berhasil dikirim !"
	return &pesan_berhasil
}
