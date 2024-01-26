package config

/**
=========================== FIXME: TODO: IMPORTANT: PENTING: ===========================
	Untuk mengaktifkan konfigurasi, file config.go pada folder config/example/ perlu
	di copy ke folder config/

	Untuk keamanan saat bekerja dengan public version control,
	harap pastikan untuk memindahkan file config.go terlebih dahulu sebelum mengubah
	konfigurasi apapun.

	Hal ini dilakukan untuk menghindari terjadinya commit file config beserta
	seluruh informasi privatenya secara tidak sengaja.
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
