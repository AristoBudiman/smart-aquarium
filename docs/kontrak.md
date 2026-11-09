# Kontrak Sistem

Dokumen ini adalah acuan bersama semua peran. Perubahan harus lewat pull request terpisah dan disetujui peran terdampak.

## Status dan keputusan

- Kelas model: `NORMAL`, `CAUTION`, `WARNING`, `SEVERE`.
- **Status utama** adalah kelas paling parah dari model suhu, model kekeruhan, dan model overall.
- Website dan notifikasi memakai status utama. Nilai `class_overall` tetap disimpan dan ditampilkan.
- Suhu acuan: 27,48 °C.

Aturan aktuator pada mode Auto:

| Kondisi | Pemanas | Kipas | Pompa |
| --- | --- | --- | --- |
| Suhu Normal | OFF | OFF | – |
| Suhu Caution | tetap | tetap | – |
| Suhu Warning/Severe, suhu < 27,48 °C | ON | OFF | – |
| Suhu Warning/Severe, suhu ≥ 27,48 °C | OFF | ON | – |
| Kekeruhan Normal | – | – | OFF |
| Kekeruhan Caution | – | – | tetap |
| Kekeruhan Warning/Severe | – | – | ON |

Backend membandingkan keputusan dengan status aktual terakhir dari topik `state`, bukan dengan perintah terakhir. Kalau status aktual belum diketahui, backend mengirim `cmd`.

## Identitas node

- `node_id` berformat `aq-` diikuti chip ID heksadesimal huruf kecil, misalnya `aq-1a2b3c`.
- Simulator memakai format `aq-sim001`.
- MQTT client ID sama dengan `node_id`, dan setiap node aktif harus unik.
- Node baru terdaftar otomatis dengan nama `Akuarium {node_id}` dan mode Manual.
- Setiap node punya kode klaim 8 karakter. Node baru belum punya pemilik sampai diklaim klien.
- Backend mengirim `config` retained saat mendaftarkan node.
- Firmware memakai mode Manual sampai `config` diterima.

## Validasi telemetry

- `temp` valid pada 0–50 °C. Di luar rentang itu, data ditolak dan backend membuat alert level `WARNING`.
- `ntu` harus 0 atau lebih.
- Firmware tidak mengirim telemetry kalau DS18B20 gagal dibaca.

## MQTT

| Topik | Arah | Isi | QoS | Retained |
| --- | --- | --- | --- | --- |
| `aquarium/{node_id}/telemetry` | Node → Backend | `{"temp":27.31,"ntu":12.4}` | 0 | Tidak |
| `aquarium/{node_id}/state` | Node → Backend | `{"heater":0,"fan":1,"pump":0,"failsafe":false}` | 0 | Tidak |
| `aquarium/{node_id}/status` | Node → Backend | `online` | 0 | Ya |
| `aquarium/{node_id}/status` | LWT broker | `offline` | 1 | Ya |
| `aquarium/{node_id}/cmd` | Backend → Node | `{"heater":0,"fan":1,"pump":0}` | 1 | Tidak |
| `aquarium/{node_id}/config` | Backend → Node | `{"mode":"auto"}` atau `{"mode":"manual"}` | 1 | Ya |
| `server/heartbeat` | Backend → semua node | `{"ts":1762171200}` | 0 | Tidak |

Node mengirim telemetry setiap 60 detik. Node mengirim `state` saat tersambung, saat relay berubah, dan setiap 60 detik. Backend mengirim heartbeat setiap 30 detik.

Backend memakai username MQTT `backend`. Semua node memakai username `node`. ACL produksi membatasi node ke topik miliknya melalui client ID.

## REST API

- Semua endpoint kecuali `POST /api/auth/login`, `POST /api/auth/register`, dan `GET /api/health` memakai `Authorization: Bearer JWT`.
- Klaim JWT memuat `uid` (ID user), `sub` (username), dan `adm` (`true` untuk admin). Token lama tanpa `adm` dianggap bukan admin.
- Semua endpoint `/api/nodes` dan `/api/nodes/{id}/...` hanya berlaku untuk node dengan `owner_id` milik user yang sedang login. `GET /api/nodes` hanya mengembalikan node milik user tersebut (admin mendapat array kosong). Node milik user lain atau belum diklaim menghasilkan `404 {"error":"Node tidak ditemukan"}`.
- Semua endpoint `/api/admin/...` wajib diakses oleh admin; user non-admin ditolak dengan `403 {"error":"Khusus admin"}`. Admin tidak dapat melihat data sensor, grafik, log, maupun mengendalikan atau mengklaim node klien.
- Waktu memakai ISO 8601 UTC, misalnya `2026-11-08T03:00:00Z`.
- Error berbentuk `{"error":"pesan"}`.

| Endpoint | Kirim | Terima |
| --- | --- | --- |
| `GET /api/health` | – | `{"status":"ok"}` |
| `POST /api/auth/login` | `{"username":"","password":""}` | `{"token":""}` |
| `POST /api/auth/register` | `{"username":"","password":""}` | `201 {"token":""}`, `400` jika tidak valid, `409` jika username sudah dipakai |
| `GET /api/me` | – | `Me` |
| `PUT /api/me/telegram` | `{"telegram_chat_id":""}` | `204` |
| `POST /api/nodes/claim` | `{"node_id":"","claim_code":""}` | `200 NodeDetail`, `403` jika admin, `404` jika node/kode salah, `409` jika milik akun lain, `429` jika >10 gagal/10 menit |
| `GET /api/nodes` | – | Array `NodeSummary` milik user |
| `GET /api/nodes/{id}` | – | `NodeDetail` |
| `PATCH /api/nodes/{id}` | `{"name":""}` | `NodeDetail` |
| `DELETE /api/nodes/{id}` | – | `204` (lepas kepemilikan, ubah ke mode Manual, `claim_code` tetap) |
| `GET /api/nodes/{id}/telemetry?from=&to=` | waktu ISO 8601 | Array `{ts,temp,ntu,status}` |
| `PUT /api/nodes/{id}/mode` | `{"mode":"auto"}` | `NodeDetail` |
| `PUT /api/nodes/{id}/actuators` | `{"heater":0,"fan":1,"pump":0}` | `202`, atau `409` kalau mode Auto |
| `GET /api/nodes/{id}/logs` | – | `{"actuators":[...],"alerts":[...]}` |
| `GET /api/admin/nodes` | – | Array `AdminNode` (urut `created_at DESC`) |
| `POST /api/admin/nodes` | `{"node_id":""}` | `201 AdminNode`, `400` jika format salah, `409` jika sudah ada |
| `POST /api/admin/nodes/{id}/claim-code` | – | `200 AdminNode` dengan `claim_code` baru |
| `DELETE /api/admin/nodes/{id}/owner` | – | `204` (lepas kepemilikan dan ubah ke mode Manual) |
| `DELETE /api/admin/nodes/{id}` | – | `204` (hapus node beserta datanya secara cascade) |

Bentuk data:

```json
{
  "Me": {
    "id": 1,
    "username": "budi",
    "is_admin": false,
    "telegram_linked": true
  },
  "AdminNode": {
    "node_id": "aq-1a2b3c",
    "name": "Akuarium aq-1a2b3c",
    "online": true,
    "last_seen": "2026-11-08T03:00:00Z",
    "created_at": "2026-11-08T01:00:00Z",
    "claim_code": "A1B2C3D4",
    "owner_username": "budi"
  },
  "NodeSummary": {
    "node_id": "aq-1a2b3c",
    "name": "Akuarium 1",
    "online": true,
    "mode": "auto",
    "last_seen": "2026-11-08T03:00:00Z",
    "temp": 27.3,
    "ntu": 12.4,
    "status": "NORMAL"
  },
  "NodeDetail": {
    "node_id": "aq-1a2b3c",
    "name": "Akuarium 1",
    "online": true,
    "mode": "auto",
    "last_seen": "2026-11-08T03:00:00Z",
    "temp": 27.3,
    "ntu": 12.4,
    "status": "NORMAL",
    "class_temp": "NORMAL",
    "class_turbidity": "NORMAL",
    "class_overall": "NORMAL",
    "heater": 0,
    "fan": 0,
    "pump": 0,
    "failsafe": false
  },
  "ActuatorLog": {
    "id": 1,
    "ts": "2026-11-08T03:00:00Z",
    "heater": 0,
    "fan": 1,
    "pump": 0,
    "source": "auto"
  },
  "Alert": {
    "id": 1,
    "ts": "2026-11-08T03:00:00Z",
    "level": "WARNING",
    "message": "Suhu 33,1 °C (WARNING)"
  }
}
```

Nilai sensor, kelas, dan aktuator boleh `null` kalau belum ada data. Rentang telemetry sampai 48 jam mengembalikan data mentah. Rentang di atas itu dirata-ratakan per jam.

## WebSocket

Browser tersambung ke `GET /ws?token=JWT`. Server mengirim JSON:

```json
{"type":"telemetry","node_id":"aq-1a2b3c","ts":"2026-11-08T03:00:00Z","temp":27.3,"ntu":12.4,"class_temp":"NORMAL","class_turbidity":"NORMAL","class_overall":"NORMAL","status":"NORMAL"}
{"type":"state","node_id":"aq-1a2b3c","heater":0,"fan":1,"pump":0,"failsafe":false}
{"type":"node_status","node_id":"aq-1a2b3c","online":true}
{"type":"mode","node_id":"aq-1a2b3c","mode":"auto"}
{"type":"alert","node_id":"aq-1a2b3c","level":"WARNING","message":"Suhu 33,1 °C (WARNING)"}
```

Aturan penyaringan event WebSocket:
- Klien biasa hanya menerima event dari node yang dimilikinya (`owner_id = user_id`).
- Admin hanya menerima event dengan `type: "node_status"` (baik untuk node yang sudah diklaim maupun node tanpa pemilik), dan tidak menerima event `telemetry`, `state`, `mode`, atau `alert`.
- Node tanpa pemilik hanya menyiarkan `node_status` ke admin.

## Notifikasi Telegram

Notifikasi Telegram untuk sebuah alert hanya dikirim ke `telegram_chat_id` milik pemilik node (`nodes.owner_id`). Jika node belum diklaim atau pemilik belum mengisi `telegram_chat_id`, alert tetap disimpan di database dan dicatat di log backend.

## Database

File `db/migrations/001_init.sql`:

```sql
CREATE TABLE users (
    id               BIGSERIAL PRIMARY KEY,
    username         TEXT NOT NULL UNIQUE,
    password_hash    TEXT NOT NULL,
    telegram_chat_id TEXT
);

CREATE TABLE nodes (
    node_id    TEXT PRIMARY KEY,
    name       TEXT NOT NULL,
    mode       TEXT NOT NULL DEFAULT 'manual' CHECK (mode IN ('auto', 'manual')),
    online     BOOLEAN NOT NULL DEFAULT FALSE,
    last_seen  TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    heater     SMALLINT CHECK (heater IN (0, 1)),
    fan        SMALLINT CHECK (fan IN (0, 1)),
    pump       SMALLINT CHECK (pump IN (0, 1)),
    failsafe   BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE telemetry (
    id      BIGSERIAL PRIMARY KEY,
    node_id TEXT NOT NULL REFERENCES nodes (node_id) ON DELETE CASCADE,
    ts      TIMESTAMPTZ NOT NULL,
    temp    DOUBLE PRECISION NOT NULL,
    ntu     DOUBLE PRECISION NOT NULL
);
CREATE INDEX telemetry_node_ts ON telemetry (node_id, ts DESC);

CREATE TABLE predictions (
    id              BIGSERIAL PRIMARY KEY,
    node_id         TEXT NOT NULL REFERENCES nodes (node_id) ON DELETE CASCADE,
    ts              TIMESTAMPTZ NOT NULL,
    class_temp      TEXT NOT NULL,
    class_turbidity TEXT NOT NULL,
    class_overall   TEXT NOT NULL,
    status          TEXT NOT NULL
);
CREATE INDEX predictions_node_ts ON predictions (node_id, ts DESC);

CREATE TABLE actuator_logs (
    id      BIGSERIAL PRIMARY KEY,
    node_id TEXT NOT NULL REFERENCES nodes (node_id) ON DELETE CASCADE,
    ts      TIMESTAMPTZ NOT NULL DEFAULT now(),
    heater  SMALLINT NOT NULL,
    fan     SMALLINT NOT NULL,
    pump    SMALLINT NOT NULL,
    source  TEXT NOT NULL CHECK (source IN ('auto', 'manual', 'failsafe'))
);
CREATE INDEX actuator_logs_node_ts ON actuator_logs (node_id, ts DESC);

CREATE TABLE alerts (
    id      BIGSERIAL PRIMARY KEY,
    node_id TEXT NOT NULL REFERENCES nodes (node_id) ON DELETE CASCADE,
    ts      TIMESTAMPTZ NOT NULL DEFAULT now(),
    level   TEXT NOT NULL,
    message TEXT NOT NULL,
    sent    BOOLEAN NOT NULL DEFAULT FALSE
);
CREATE INDEX alerts_node_ts ON alerts (node_id, ts DESC);
```

File `db/migrations/002_akun.sql` menambahkan dukungan peran admin (`users.is_admin`), kepemilikan node (`nodes.owner_id`), dan kode klaim 8 karakter (`nodes.claim_code`):

```sql
ALTER TABLE users ADD COLUMN is_admin BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE nodes ADD COLUMN owner_id BIGINT REFERENCES users (id) ON DELETE SET NULL;
ALTER TABLE nodes ADD COLUMN claim_code TEXT NOT NULL DEFAULT upper(substr(md5(random()::text || clock_timestamp()::text), 1, 8));
CREATE INDEX nodes_owner ON nodes (owner_id);
```

## Variabel lingkungan

| Nama | Dipakai | Keterangan |
| --- | --- | --- |
| `DOMAIN` | Caddy | Domain DuckDNS |
| `POSTGRES_PASSWORD` | PostgreSQL, backend | Password database |
| `DATABASE_URL` | Backend | Dibentuk di Compose |
| `MQTT_URL` | Backend | Dibentuk di Compose |
| `MQTT_USERNAME` | Backend | Kosong di lokal, `backend` di produksi |
| `MQTT_PASSWORD` | Backend | Kosong di lokal |
| `JWT_SECRET` | Backend | Minimal 32 karakter acak |
| `ADMIN_USERNAME` | Backend | Akun awal |
| `ADMIN_PASSWORD` | Backend | Akun awal |
| `TELEGRAM_BOT_TOKEN` | Backend | Kosong berarti notifikasi hanya dicatat di log |
