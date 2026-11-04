# AGENTS.md: AI Smart Aquarium

Tugas besar kelompok mata kuliah Teknologi IoT. Sistem memantau suhu dan kekeruhan akuarium, mengklasifikasikan kondisi dengan tiga model Random Forest, lalu mengendalikan pemanas, kipas, dan pompa filter. Sistem mendukung banyak akuarium (multi-node) dan tidak memakai platform IoT pihak ketiga.

## Cara memulai sesi

1. Baca file ini dan `docs/kontrak.md`.
2. Tanyakan peran pengguna kalau belum disebutkan, lalu baca instruksi perannya di `docs/agent/`.
3. Cek `git status` dan branch aktif.
4. Sampaikan rencana singkat sebelum mengubah banyak file.
5. Setelah selesai, jalankan pemeriksaan yang tertulis di instruksi peran dan laporkan hasilnya apa adanya.

## Urutan kebenaran

Kalau ada perbedaan, ikuti urutan ini:

1. `docs/kontrak.md`
2. Instruksi peran di `docs/agent/`
3. `docs/referensi/Laporan_Baru_Bab1-4.md`

Laporan dipakai sebagai latar belakang. Beberapa bagiannya sudah dikoreksi oleh kontrak.

## Teknologi final

Jangan mengganti teknologi berikut tanpa persetujuan kelompok:

| Lapisan | Teknologi |
| --- | --- |
| Perangkat | NodeMCU ESP8266, DS18B20 di D5, turbidity di A0, relay IN1 D6 pemanas, IN2 D7 kipas, IN3 D1 pompa |
| Komunikasi | MQTT dengan broker Mosquitto |
| Backend | Go, inferensi dari kode hasil m2cgen |
| Model | 3 Random Forest, 50 pohon, kedalaman 8, suhu acuan 27,48 °C |
| Database | PostgreSQL |
| Website | Vue 3, Vite, Tailwind CSS, ECharts, login JWT |
| Notifikasi | Bot Telegram |
| Hosting | Docker Compose di VM Azure for Students, Caddy, DuckDNS |

Blynk, ThingSpeak, Firebase, dan platform IoT sejenis dilarang. Jangan menambahkan framework besar atau layanan berbayar.

## Kepemilikan folder

Ubah folder milik peran sendiri saja. Kalau perlu mengubah folder lain, hentikan pekerjaan dan beri tahu pengguna.

| Peran | Folder |
| --- | --- |
| 1. Perangkat, fondasi, infrastruktur | `firmware/`, `deploy/`, `compose.yaml`, `backend/cmd/`, `backend/Dockerfile`, `frontend/Dockerfile`, `.env.example` |
| 2. Simulator, pengujian, dokumentasi | `tools/simulator/`, `docs/laporan/` |
| 3. ML dan logika backend | `ml/`, `backend/internal/model`, `ingest`, `decision`, `notify`, `mqtt` |
| 4. API backend | `backend/internal/api`, `auth`, `ws`, `db`, `db/migrations/` |
| 5. Frontend | `frontend/`, kecuali `frontend/Dockerfile` |

File bersama adalah `docs/kontrak.md` dan `backend/internal/core/ports.go`. Keduanya hanya boleh berubah lewat pull request terpisah yang disetujui peran terdampak.

## Aturan wajib

- Go tidak terpasang di laptop. Jalankan semua perintah Go lewat Docker. Jangan meminta pengguna memasang Go.
- Jangan membaca, menampilkan, atau meng-commit `.env`, `firmware/smart_aquarium/config.h`, `deploy/mosquitto/passwd`, dan file kunci.
- Jangan menulis password, token, atau alamat rahasia di dalam kode atau dokumen.
- Jangan push langsung ke `main`. Pakai branch `feat/firmware`, `feat/infra`, `feat/sim`, `feat/ml`, `feat/api`, atau `feat/web`.
- Mintalah izin sebelum menjalankan perintah yang menghapus data, misalnya `docker compose down -v`, `git reset --hard`, atau `rm -r`.
- Jangan menjalankan perintah di VM atau mengubah Azure kecuali pengguna meminta secara langsung.
- Jangan menandai pekerjaan selesai sebelum pemeriksaannya dijalankan. Kalau pemeriksaan tidak bisa dijalankan, katakan alasannya.
- Jangan menciptakan field, endpoint, atau topik MQTT di luar kontrak.
- Di Windows, jalankan perintah Docker dari PowerShell. Pada PowerShell 5, tanda kutip ganda di JSON pada argumen perintah perlu ditulis `\"`.
- Teks antarmuka memakai bahasa Indonesia. Nama variabel dan kode mengikuti kebiasaan bahasanya.

## Perintah harian

```bash
docker compose up -d
docker compose logs -f backend
docker compose restart backend
docker compose run --rm --no-deps backend go test ./...
docker compose run --rm --no-deps backend go mod tidy
docker compose down
cd frontend && npm run dev
```

Simulator dijalankan dengan `docker compose run --rm simulator` setelah Peran 2 membuatnya.
