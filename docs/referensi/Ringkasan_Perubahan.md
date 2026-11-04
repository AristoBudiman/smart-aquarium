# Ringkasan Perubahan Laporan AI Smart Aquarium

Dokumen ini menjelaskan apa yang berubah dari laporan terakhir (versi Blynk) ke laporan baru (`Laporan_Baru_Bab1-4.md`), dan alasannya.

## 1. Kenapa ada perubahan

Ada tiga permintaan baru:
1. **Tidak boleh memakai Blynk, ThingSpeak, atau platform IoT sejenis.** Server, basis data, dan website harus dibuat sendiri.
2. **Sistem harus bisa menampung penambahan node (akuarium) baru.**
3. **Studi literatur ditambah menjadi 10 penelitian terkait** dari ScienceDirect atau jurnal yang bisa diakses gratis.

Yang **tidak berubah**: hardware dan wiring, tiga model Random Forest (suhu, kekeruhan, overall), logika aksi aktuator di Tabel 4.1, suhu acuan 27,48 °C, dan mekanisme fail-safe.

## 2. Arsitektur baru yang disepakati

| Bagian | Sebelumnya | Sekarang |
|---|---|---|
| Protokol node ↔ server | Library Blynk | **MQTT** |
| Broker / pusat data | Blynk.Cloud | **Mosquitto** (self-host, open source) |
| Server klasifikasi | Python di Render, polling Blynk tiap 60 detik | **Backend Go**, berbasis event (data masuk langsung diproses) |
| Inferensi model | File `.joblib` dimuat Python | Model dikonversi menjadi **kode Go dengan m2cgen** |
| Basis data | Riwayat Blynk (1 minggu) | **PostgreSQL** (tanpa batas penyimpanan dari platform) |
| Website | Dashboard Blynk | **Vue 3 + Vite + Tailwind + ECharts**, login **JWT** |
| Data real-time ke website | Otomatis oleh Blynk | **WebSocket** dari backend Go |
| Notifikasi | Blynk Events | **Bot Telegram** + riwayat di website |
| Hosting | Render (gratis, bisa tidur) | **VM Azure for Students**, semua layanan dalam **Docker Compose** |
| HTTPS & domain | Disediakan Blynk | **Caddy** (HTTPS otomatis) + **DuckDNS** (subdomain gratis) |
| Library ESP | Blynk | **PubSubClient + ArduinoJson + WiFiManager** |

Isi Docker Compose: `mosquitto`, `backend` (Go), `db` (PostgreSQL), `caddy` (reverse proxy + menyajikan website Vue).

Catatan: website Vue juga bisa di-deploy ke Vercel kalau tim lebih suka. Broker, backend, dan PostgreSQL tetap harus di VM karena harus menyala terus.

## 3. Cara kerja multi-node (inti perubahan desain)

- Setiap ESP memakai **chip ID** sebagai `node_id`.
- Semua topik MQTT memakai `node_id`: `aquarium/{node_id}/telemetry`, `/state`, `/status`, `/cmd`, `/config`, ditambah `server/heartbeat`.
- Backend subscribe ke `aquarium/+/telemetry`. Node yang belum dikenal **terdaftar otomatis** dengan mode bawaan Manual.
- Status online/offline memakai **Last Will (LWT)** MQTT.
- Mode Auto/Manual disimpan per node (topik `config`, retained).
- Menambah akuarium cukup dengan *flash* firmware yang sama, lalu mengatur WiFi lewat portal WiFiManager. Tidak ada kode server yang diubah.
- Untuk pengujian: 1 node fisik + beberapa node simulasi (program pengirim data palsu).

## 4. Padanan istilah lama → baru

| Lama (Blynk) | Baru |
|---|---|
| V0, V1 (suhu, NTU) | `aquarium/{node_id}/telemetry` |
| V2, V8, V9 (kelas) | Tabel `predictions` di PostgreSQL, dikirim ke website lewat WebSocket |
| V3 (mode) | `aquarium/{node_id}/config` + kolom `nodes.mode` |
| V4–V6 (aktuator) | `aquarium/{node_id}/cmd` (perintah) dan `/state` (status aktual) |
| V7 (heartbeat) | `server/heartbeat` |
| `getAll` / `update` HTTPS API Blynk | Subscribe / publish MQTT |
| Blynk Events | Bot Telegram |

## 5. Perubahan per bab

| Bagian | Status | Isi perubahan |
|---|---|---|
| 1.1 | Diubah sedikit | "platform cloud" → server buatan sendiri; antarmuka website; kalimat multi-node |
| 1.2 | Diubah | Tambah rumusan masalah **d** (penambahan node & tanpa platform pihak ketiga) |
| 1.3 b, c | Diubah | Blynk/HTTPS → MQTT; notifikasi lewat website + Telegram; website responsif |
| 1.3 f | Baru | Skalabilitas dan kemandirian sistem |
| 2.1 | Diubah | Blynk Cloud → Mosquitto + Go + PostgreSQL + Vue di VM; node_id & pendaftaran otomatis |
| 2.2 | Diubah | Batasan: satu VM; uji multi-node dengan simulator; target pengguna multi-akuarium |
| 2.3.1 | Diubah | Tambah fitur pemantauan multi-node, riwayat data, notifikasi |
| 2.3.2 | Diubah sedikit | Mode diatur per node; "Web/Mobile Dashboard" → website |
| 3.1.1 c | Diganti | Teori Blynk → **teori MQTT** (publish/subscribe, QoS, retained, LWT) |
| 3.1.2 | Diperluas | Tabel penelitian terkait 3 → **10 penelitian** |
| 3.1.3 | Ditulis ulang | Tiga kesenjangan: ML tidak dipakai untuk aksi, ketergantungan platform, desain satu lokasi |
| 3.1.4 | Diubah | Model dibatasi 50 pohon, kedalaman 8 (agar bisa dikonversi m2cgen); tabel hasil disesuaikan |
| 3.2.1 | Diubah sedikit | "ke Blynk.Cloud" → "ke broker MQTT"; chip ID sebagai node_id |
| 3.2.2, 3.2.3 | Tetap | – |
| 3.2.4 | Ditulis ulang | Mosquitto, backend Go, m2cgen, PostgreSQL, Vue, Telegram, Docker Compose di VM; 2 tabel baru |
| 4.1.1 | Diubah | Alur WiFiManager, MQTT, LWT; notifikasi dipindah ke server |
| 4.1.2 | Diubah | Server berbasis event, pendaftaran node otomatis, heartbeat paralel |
| Tabel 4.1 | Tetap | Kolom V4/V5/V6 diganti menjadi nama aktuator saja |
| 4.2 | Diubah | Lapisan Cloud = VM berisi 4 kontainer; UI = website + Telegram; multi-node |
| 4.3.1, 4.3.2 | Diubah | Topologi multi-node ke VM; tabel protokol baru |
| 4.3.3 | Diganti | Datastream Blynk → **tabel topik MQTT** |
| 4.3.4 – 4.3.6 | Baru | REST API & WebSocket, rancangan basis data (+ ERD), mekanisme penambahan node |
| 4.4 | Tetap | Nomor gambar wiring bergeser jadi 4.6 |
| 4.5 | Ditulis ulang | 6 halaman website; mockup baru menggantikan Blynk |
| Daftar Pustaka | Diperluas | [1]–[6] tetap, tambah [7]–[13] (penelitian baru) dan [14] (standar MQTT) |

## 6. Gambar yang perlu dibuat atau diubah

| Gambar | Status | Keterangan singkat |
|---|---|---|
| 4.1 Flowchart perangkat | Gambar ulang | Tambah WiFiManager & koneksi MQTT; hapus keputusan notifikasi |
| 4.2 Flowchart server | Gambar ulang | Berbasis event, pendaftaran node, heartbeat paralel |
| 4.3 Arsitektur | Gambar ulang | 2–3 node, kotak VM berisi 4 kontainer, website + Telegram |
| 4.4 Topologi | Gambar ulang | Multi-node, VM dengan port 1883/80/443 |
| 4.5 ERD | **Baru** | 6 tabel (dbdiagram.io / draw.io) |
| 4.6 Wiring | Tetap | Hanya nomor bergeser |
| 4.7 – 4.10 Mockup | **Baru** | Daftar node, detail node (desktop & ponsel), notifikasi Telegram (Figma) |

Detail isi setiap gambar ada di kotak "PERLU DIGAMBAR" pada laporan baru.

## 7. Sepuluh penelitian terkait

| No | Penelitian | Akses |
|---|---|---|
| 1–3 | Nasir & Natasya 2020; Qolbu dkk. 2026; Andini dkk. 2024 | Sudah ada sebelumnya |
| 4 | Islam dkk. 2023, *Microprocessors and Microsystems* | ScienceDirect |
| 5 | Adeleke dkk. 2023, *Internet of Things* | ScienceDirect |
| 6 | Olanubi dkk. 2024, *J. Electrical Systems and Information Technology* | Open access (SpringerOpen) |
| 7 | Kanwal dkk. 2024, *Sensors* | Open access (MDPI) |
| 8 | Arepalli & Naik 2024, *Ecological Informatics* | ScienceDirect, open access |
| 9 | Shete dkk. 2024, *MethodsX* | ScienceDirect, open access |
| 10 | Hridoy dkk. 2025, *Results in Chemistry* | ScienceDirect, open access |

Benang merah dari tabel: model ML yang akurat belum dipakai untuk menggerakkan aktuator, sistem bergantung pada Blynk/ThingSpeak/Firebase, dan desainnya untuk satu lokasi. Ketiga celah ini yang dijawab proyek kita.

## 8. Yang masih perlu dikerjakan

- [ ] Latih ulang model di Colab dengan `n_estimators=50, max_depth=8`, lalu perbarui Tabel 3.2 dengan angka hasil sendiri.
- [ ] Isi `[isi ...]` yang tersisa: tahun skripsi Revantianto [5], data dataset [6], dan nama subdomain DuckDNS.
- [ ] Buat atau perbarui gambar sesuai Bagian 6.
- [ ] Perbarui PPT: slide 3.1 (Blynk → MQTT, tabel 10 penelitian), 3.2 (Blynk.Cloud & software), 4.1–4.3, dan 4.5.
- [ ] Daftar akun Azure for Students dan DuckDNS sebelum tahap implementasi.
