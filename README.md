# PD2-Golang.FinalProject
Kelompok 3 - Final Project Pemrograman Dasar 2

<b>Anggota Kelompok</b>
<ul>
  <li>Matthew Tirtawidjaja (<a href="https://github.com/TIRTAGT">@TIRTAGT</a>)</li>
  <li>Katon (<a href="https://github.com/tonkaton">@tonkaton</a>)</li>
</ul>

----

### Cara Instalasi / Penggunaan Program
<ol>
	<li>Clone/Download seluruh repository <a href="https://github.com/TIRTAGT/PD2-Golang.FinalProject">PD2-Golang.FinalProject</a> ke folder yang diinginkan.</li>
	<!-- insert foto contoh cara clone/download disini -->
	<li>Buka folder tersebut di terminal (CMD / PowerShell / Terminal / Git Bash).</li>
	<!-- insert foto contoh cara buka folder di terminal -->
	<li>Sesuaikan alamat IP dan/atau port dimulainya HTTP server proyek ini</li>
	<!-- insert foto contoh caranya -->
	<li>Jalankan program dengan perintah <kbd>go run .</kbd></li>
	<!-- insert foto contoh hasil menjalankan perintahnya -->
	<li>Kunjungi alamat IP dan port HTTP server yang tertera di terminal pada web browser, atau CTRL+klik (jika lingkungan perangkat mendukung)</li>
	<!-- insert foto tampilan awal web projectnya -->
</ol>

----

### Penjelasan Struktur Folder dan File Project

``assets/``

Folder ini berisi source code untuk frontend (CSS/JS) dan file-file framework frontend yang digunakan (contoh: Bootstrap, JQuery, AdminLTE, dll).

``controllers/``

Folder ini berisi source code untuk web backend project Golang Web Server ini.

``server/``

Folder ini berisi source code untuk implementasi Golang Web Server yang dibuat untuk melakukan simulasi teknik MVC (Model-View-Controller) pada Golang.

``views/``

Folder ini berisi source code untuk halaman frontend yang akan diambil dari controller pada rute tersebut (HTML).



``.gitignore``

File berisi daftar file/folder yang tidak perlu diupload ke git/version control

``go.mod``

File untuk menandakan project ini adalah sebuah module golang

``go.work``

File berisi module lokal apa saja yang digunakan oleh project ini

``README.md``

File ini

``server.go``

File perlu dijalankan untuk memulai server


----

### Framework/Tools yang digunakan

<ul>
	<li>Golang (Go) - go1.21.2</li>
</ul>
