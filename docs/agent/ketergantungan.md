# Ketergantungan Antarperan

Baca file ini sebelum mengerjakan tugas. Tujuannya mencegah agent membuat bagian yang masih bergantung pada pekerjaan peran lain.

## Aturan utama

- Sumber status adalah branch `origin/main`, bukan branch teman, chat, atau file lokal.
- Sebelum memulai tugas, cek semua prasyarat tugas itu.
- Kalau prasyarat belum ada, jangan membuat implementasi asli untuk bagian peran lain.
- Jangan mengubah folder peran lain untuk mengatasi hambatan.
- Jangan mengubah `docs/kontrak.md` atau `ports.go` tanpa persetujuan peran terdampak.
- Kalau hanya butuh antarmuka, pakai interface dari `ports.go`. Implementasi palsu hanya boleh berada di file `_test.go` atau di folder sendiri, dan harus diberi komentar `SEMENTARA`.
- Kalau tugas benar-benar tertahan, hentikan tugas itu dan kerjakan tugas lain yang prasyaratnya sudah ada.

## Cara mengecek prasyarat

Jalankan dari root repo:

```bash
git fetch origin
git cat-file -e origin/main:PATH
```

Exit code 0 berarti path ada di `main`. Kalau agent tidak bisa menjalankan `git fetch`, minta pengguna menjalankannya.

| Kode | Prasyarat | Path yang dicek | Pemilik |
| --- | --- | --- | --- |
| P1 | Fondasi dan interface bersama | `backend/internal/core/ports.go` | Peran 1 |
| P2 | Simulator | `tools/simulator/main.py` | Peran 2 |
| P3 | Model Go | `backend/internal/model/model.go` | Peran 3 |
| P4 | Migrasi database | `db/migrations/001_init.sql` | Peran 4 |
| P5 | Implementasi `Store` | `backend/internal/db` | Peran 4 |
| P6 | Implementasi `Commander` | `backend/internal/mqtt` | Peran 3 |
| P7 | Implementasi `Broadcaster` | `backend/internal/ws` | Peran 4 |
| P8 | REST API | `backend/internal/api` | Peran 4 |
| P9 | Firmware MQTT | `firmware/smart_aquarium/smart_aquarium.ino` | Peran 1 |

## Batas tiap peran

| Peran | Boleh sekarang | Wajib menunggu | Tidak boleh dilakukan |
| --- | --- | --- | --- |
| 1 | Firmware, tes alat, deploy, migrasi VM setelah P4 | Uji Auto dengan backend menunggu P5 dan P6 | Mengubah folder backend, frontend, atau simulator milik peran lain |
| 2 | Simulator, skrip uji, gambar, laporan | Uji alur `cmd` asli menunggu P5 dan P6; uji website menunggu P8 | Mengubah backend atau frontend untuk membuat tes lolos |
| 3 | Model, `decision`, MQTT, `ingest`, `notify`, unit test dengan `Store` palsu | Menjalankan alur penuh menunggu P4 dan P5; menyambungkan `main.go` menunggu P5, P6, dan P7 | Menulis SQL, tabel, query database, API, atau WebSocket |
| 4 | Migrasi, `db`, auth, REST, WebSocket, `Commander` sementara di folder sendiri | Perintah aktuator dan mode asli menunggu P6; menyambungkan `main.go` menunggu P5, P6, dan P7 | Publish MQTT langsung, membuat logika model, atau mengubah aturan aktuator |
| 5 | Halaman, store, grafik, tampilan ponsel, data palsu | Mematikan data palsu menunggu P8; WebSocket asli menunggu P7 | Mengubah API, menebak field baru, atau membuat backend tiruan di luar `src/api/` |

## Urutan merge

1. Peran 4: P4 dan P5.
2. Peran 3: P3 dan P6.
3. Peran 4: P7 dan P8.
4. Peran 3 dan 4: penyambungan `backend/cmd/server/main.go`.
5. Peran 5: pindah dari data palsu ke API asli.
6. Peran 1: deploy ke VM dan jalankan migrasi.

Peran 2 dan firmware Peran 1 boleh merge kapan saja, selama hanya mengubah folder sendiri.

## Jika tertahan

Laporkan dengan format ini:

```text
TERTAHAN
Tugas: ...
Menunggu: P5 Implementasi Store
Pemilik: Peran 4
Sudah dicek: git cat-file -e origin/main:backend/internal/db
Tugas lain yang bisa dikerjakan: ...
```
