# PD2-Golang.FinalProject
Kelompok 3 - Final Project Pemrograman Dasar 2

----

Folder ini berisi source code untuk web backend project Golang Web Server ini.

----

## Cara buat controller baru
1. Buat folder baru sesuai nama rute yang diinginkan
2. Jalankan perintah untuk membuat module golang baru:
```bash
go mod init github.com/TIRTAGT/PD2-Golang.FinalProject/controllers/<nama_rute>
```
3. Buat file golang baru dengan nama rute, package harus sesuai dengan nama rute
```go
// file bernama: <nama_rute>.go
package <nama rute>

import (
	"net/http"
	"github.com/TIRTAGT/PD2-Golang.FinalProject/server/controller/utility"
)
```
4. Tambahkan lokasi module golang untuk rute yang baru tersebut ke file go.work di folder root (folder utama)
```bash
use (
	// ...
	./controllers
	./controllers/<nama_rute>
	./server
	// ...
)
```
5. Tambahkan rute baru ke file controllers.go di folder controllers
```go
var ControllerMap = map[string]handlerstruct.ControllerStruct {
	"/<nama_rute>": {
		GET: <nama_rute>.GET,
		// isi sesuai jika butuh metode lain
	},
}
```
6. Done, silahkan coba akses.

----

## Cara hapus controller yang sudah ada
1. Hapus folder controller yang rutenya ingin dihapus
2. Hapus lokasi module tersebut dari file go.work di folder root (folder utama)
```bash
use (
	// ...
	./controllers
	./controllers/<nama_rute>    <- hapus ini
	./server
	// ...
)
```
3. Hapus rute controlller tersebut dari file controllers.go di folder controllers
```go
var ControllerMap = map[string]handlerstruct.ControllerStruct {
	"/<nama_rute>": {									-
		GET: <nama_rute>.GET,							|	Hapus semua bagian ini
		// metode lain jika ada							|
	},													-
}
```
4. Done, seharusnya sudah tidak bisa akses