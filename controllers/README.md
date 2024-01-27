# PD2-Golang.FinalProject
<b>Kelompok 3 - Final Project Pemrograman Dasar 2</b><br>
<b>Anggota Kelompok :</b>
<ul>
  <li>Matthew Tirtawidjaja (<a href="https://github.com/TIRTAGT">@TIRTAGT</a>)</li>
  <li>Katon Kurnia Wijaya (<a href="https://github.com/tonkaton">@tonkaton</a>)</li>
</ul>

----

Folder ini berisi source code untuk web backend project Golang Web Server ini.

----

## Cara membuat controller baru
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
6. Restart server
7. Done, seharusnya sudah bisa akses

----

## Cara hapus controller yang sudah ada
1. Hapus folder controller yang rutenya ingin dihapus
2. Hapus lokasi module tersebut dari file go.work di folder root (folder utama)
```bash
use (
	// ...
	./controllers
	./controllers/<nama_rute>    <- hapus ini
	
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
4. Restart server
5. Done, seharusnya sudah tidak bisa akses