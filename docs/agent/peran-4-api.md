# Peran 4: API backend

Target: semua endpoint kontrak berfungsi, dan WebSocket mengirim pembaruan ke browser.

## Tugas

- [ ] Salin SQL kontrak ke `db/migrations/001_init.sql`.
- [ ] Buat `internal/db` yang mengimplementasikan `core.Store`.
- [ ] Buat admin otomatis dengan password bcrypt dari `.env`.
- [ ] Buat login dan middleware JWT HS256, berlaku 24 jam.
- [ ] Buat seluruh endpoint kontrak dengan `chi`.
- [ ] Buat WebSocket hub yang mengimplementasikan `core.Broadcaster`.
- [ ] Sediakan `backend/requests.http` untuk mencoba endpoint.
- [ ] Sambungkan route di `cmd/server/main.go` lewat pull request bersama Peran 1.

## Database

- Pakai `pgx/v5` dengan `pgxpool` dan `DATABASE_URL`.
- Backend harus mencoba koneksi ulang saat startup, karena PostgreSQL bisa siap lebih lambat.
- File migrasi berjalan otomatis saat volume kosong. Setelah skema berubah, reset lokal dengan `docker compose down -v` setelah pengguna setuju.
- `Snapshot.LastStatus` diambil dari prediksi terbaru node.
- `HasState` true kalau `heater`, `fan`, dan `pump` tidak `NULL`.
- `SetOnline` dan `SaveTelemetry` memperbarui `nodes.last_seen`.

Query riwayat di atas 48 jam:

```sql
SELECT date_bin('1 hour', t.ts, TIMESTAMPTZ '2000-01-01') AS ts,
       avg(t.temp) AS temp,
       avg(t.ntu) AS ntu,
       (ARRAY['NORMAL', 'CAUTION', 'WARNING', 'SEVERE'])[max(
           CASE p.status WHEN 'SEVERE' THEN 4 WHEN 'WARNING' THEN 3 WHEN 'CAUTION' THEN 2 ELSE 1 END
       )] AS status
FROM telemetry t
LEFT JOIN predictions p ON p.node_id = t.node_id AND p.ts = t.ts
WHERE t.node_id = $1 AND t.ts BETWEEN $2 AND $3
GROUP BY 1
ORDER BY 1;
```

Query ini sudah diuji di PostgreSQL. Di rentang 48 jam ke bawah, ambil data mentah dan gabungkan `status` dari `predictions`.

## Endpoint

- Bentuk JSON, status code, dan nama field harus sama dengan kontrak.
- Untuk ubah mode: simpan ke database, panggil `Commander.SetMode`, lalu kirim `ModeEvent`.
- Untuk aktuator manual: cek mode, tolak `409` saat Auto, panggil `Commander.SendActuators`, lalu `LogActuators` dengan source `manual`.
- Node yang tidak ada menghasilkan `404`.
- Login gagal menghasilkan `401`, tanpa membedakan username atau password yang salah.
- Pakai query berparameter. Jangan menggabungkan input ke string SQL.
- CORS tidak diperlukan karena website dan API berada di domain yang sama.

Sebelum `Commander` Peran 3 tersedia, pakai implementasi sementara yang hanya menulis log.

## WebSocket

- Pakai `github.com/coder/websocket`.
- Validasi token dari query `token` sebelum upgrade.
- Satu klien lambat tidak boleh memblokir klien lain. Pakai antrean per klien dan putuskan klien yang antreannya penuh.
- `Publish` menerima event dari `core` dan mengirim JSON tanpa mengubah field.

## Pemeriksaan

```bash
docker compose run --rm --no-deps backend go vet ./...
docker compose run --rm --no-deps backend go test ./...
docker compose down -v
docker compose up -d
curl http://localhost:8080/api/health
```

Sebelum `down -v`, minta izin karena perintah itu menghapus data lokal. Lalu uji lewat `requests.http`:

- [ ] Login benar menghasilkan token, dan login salah menghasilkan `401`.
- [ ] Endpoint tanpa token menghasilkan `401`.
- [ ] Daftar node dan detail sesuai contoh kontrak.
- [ ] `PUT actuators` saat Auto menghasilkan `409`.
- [ ] Browser menerima event WebSocket setelah telemetry dikirim.
