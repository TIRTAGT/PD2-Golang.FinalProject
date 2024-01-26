package config

/**
=========================== FIXME: TODO: IMPORTANT: PENTING: ===========================
	Untuk keamanan saat bekerja dengan public version control,
	harap pastikan untuk rename file example.config.go menjadi config.go sebelum melakukan
	perubahan konfigurasi apapun.

	Walau filenya akan tetap bekerja tanpa direname, sangat disarankan rename untuk
	menghindari commit file config beserta seluruh informasi privatenya.
=========================================================================================
**/

const (
	HTTP_LISTEN_ADDRESS = "localhost"
	HTTP_LISTEN_PORT    = 8080

	DATABASE_HOSTNAME = HTTP_LISTEN_ADDRESS
	DATABASE_PORT     = 3306
	DATABASE_USERNAME = "root"
	DATABASE_PASSWORD = "root"
	DATABASE_NAME     = "PD2_Golang_FinalProject"
)
