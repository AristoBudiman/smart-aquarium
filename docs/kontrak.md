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

- Semua endpoint kecuali `POST /api/auth/login` dan `GET /api/health` memakai `Authorization: Bearer JWT`.
- Waktu memakai ISO 8601 UTC, misalnya `2026-11-08T03:00:00Z`.
- Error berbentuk `{"error":"pesan"}`.

| Endpoint | Kirim | Terima |
| --- | --- | --- |
| `GET /api/health` | – | `{"status":"ok"}` |
| `POST /api/auth/login` | `{"username":"","password":""}` | `{"token":""}` |
| `GET /api/nodes` | – | Array `NodeSummary` |
| `GET /api/nodes/{id}` | – | `NodeDetail` |
| `PATCH /api/nodes/{id}` | `{"name":""}` | `NodeDetail` |
| `GET /api/nodes/{id}/telemetry?from=&to=` | waktu ISO 8601 | Array `{ts,temp,ntu,status}` |
| `PUT /api/nodes/{id}/mode` | `{"mode":"auto"}` | `NodeDetail` |
| `PUT /api/nodes/{id}/actuators` | `{"heater":0,"fan":1,"pump":0}` | `202`, atau `409` kalau mode Auto |
| `GET /api/nodes/{id}/logs` | – | `{"actuators":[...],"alerts":[...]}` |
| `PUT /api/me/telegram` | `{"telegram_chat_id":""}` | `204` |

Bentuk data:

```json
{
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
