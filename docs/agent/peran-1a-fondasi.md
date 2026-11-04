# Peran 1A: fondasi repo

Target: anggota lain bisa `git clone`, menyalin `.env.example`, menjalankan `docker compose up -d`, dan membuka `http://localhost:8080/api/health`. Go tidak terpasang di laptop, jadi semua proses Go berjalan lewat Docker.

## Tugas

- [ ] Buat struktur folder.
- [ ] Buat kerangka frontend Vue.
- [ ] Buat kerangka backend Go lewat Docker.
- [ ] Buat `ports.go` dan server health.
- [ ] Buat `compose.yaml`, `dev.conf`, dan `.env.example`.
- [ ] Tulis README singkat.
- [ ] Jalankan semua pemeriksaan.

## Struktur

```text
smart-aquarium/
├─ firmware/smart_aquarium/
├─ ml/dataset/
├─ backend/
│  ├─ cmd/server/main.go
│  └─ internal/core/ports.go
├─ db/migrations/
├─ frontend/
├─ tools/simulator/
├─ deploy/mosquitto/dev.conf
├─ docs/
├─ compose.yaml
├─ .env.example
└─ README.md
```

Tambahkan `.gitkeep` pada folder kosong. Jangan membuat file `001_init.sql` pada tahap ini, karena Peran 4 yang menambahkannya.

## Frontend

Jalankan dari root repo:

```bash
npm create vue@latest frontend -- --ts --router --pinia
cd frontend
npm install
npm install -D tailwindcss @tailwindcss/vite
npm install echarts vue-echarts axios
```

Di `frontend/vite.config.ts`, tambahkan plugin `@tailwindcss/vite`, lalu tambahkan:

```ts
server: {
  proxy: {
    '/api': 'http://localhost:8080',
    '/ws': { target: 'ws://localhost:8080', ws: true },
  },
},
```

Ganti isi CSS utama menjadi `@import "tailwindcss";`, lalu kosongkan halaman contoh bawaan.

## Backend

Jalankan dari root repo. Di Windows, pakai PowerShell, karena Git Bash mengubah path pada argumen Docker. Ganti `NAMA-AKUN` dengan akun GitHub pemilik repo.

```bash
mkdir backend
docker run --rm -v "${PWD}/backend:/app" -w /app golang:1.26-alpine sh -c "go mod init github.com/NAMA-AKUN/smart-aquarium/backend && go get github.com/eclipse/paho.mqtt.golang github.com/jackc/pgx/v5 github.com/go-chi/chi/v5 github.com/golang-jwt/jwt/v5 github.com/coder/websocket golang.org/x/crypto"
```

Buat `backend/internal/core/ports.go` persis seperti berikut:

```go
// Package core berisi tipe dan interface bersama Peran 3 dan Peran 4.
// Perubahan file ini sama dengan perubahan kontrak.
package core

import (
	"context"
	"time"
)

type Class string

const (
	ClassNormal  Class = "NORMAL"
	ClassCaution Class = "CAUTION"
	ClassWarning Class = "WARNING"
	ClassSevere  Class = "SEVERE"
)

func (c Class) Rank() int {
	switch c {
	case ClassSevere:
		return 3
	case ClassWarning:
		return 2
	case ClassCaution:
		return 1
	default:
		return 0
	}
}

// Worst mengembalikan kelas paling parah.
func Worst(classes ...Class) Class {
	hasil := ClassNormal
	for _, c := range classes {
		if c.Rank() > hasil.Rank() {
			hasil = c
		}
	}
	return hasil
}

type Mode string

const (
	ModeAuto   Mode = "auto"
	ModeManual Mode = "manual"
)

type Source string

const (
	SourceAuto     Source = "auto"
	SourceManual   Source = "manual"
	SourceFailsafe Source = "failsafe"
)

type Actuators struct {
	Heater int `json:"heater"`
	Fan    int `json:"fan"`
	Pump   int `json:"pump"`
}

type ActuatorState struct {
	Actuators
	Failsafe bool `json:"failsafe"`
}

type Telemetry struct {
	NodeID string    `json:"node_id"`
	TS     time.Time `json:"ts"`
	Temp   float64   `json:"temp"`
	NTU    float64   `json:"ntu"`
}

type Classes struct {
	Temp      Class `json:"class_temp"`
	Turbidity Class `json:"class_turbidity"`
	Overall   Class `json:"class_overall"`
	Status    Class `json:"status"`
}

// NodeSnapshot adalah kondisi terakhir yang dibutuhkan logika keputusan.
type NodeSnapshot struct {
	NodeID     string
	Mode       Mode
	Online     bool
	HasState   bool
	State      ActuatorState
	LastStatus Class
}

type Alert struct {
	ID      int64     `json:"id"`
	NodeID  string    `json:"node_id"`
	TS      time.Time `json:"ts"`
	Level   Class     `json:"level"`
	Message string    `json:"message"`
}

// Store diimplementasikan Peran 4 di internal/db dan dipakai Peran 3.
type Store interface {
	EnsureNode(ctx context.Context, nodeID string) (created bool, err error)
	Snapshot(ctx context.Context, nodeID string) (NodeSnapshot, error)
	SaveTelemetry(ctx context.Context, t Telemetry) error
	SavePrediction(ctx context.Context, nodeID string, ts time.Time, c Classes) error
	SaveState(ctx context.Context, nodeID string, s ActuatorState) error
	SetOnline(ctx context.Context, nodeID string, online bool) error
	LogActuators(ctx context.Context, nodeID string, a Actuators, source Source) error
	SaveAlert(ctx context.Context, a Alert) (id int64, err error)
	MarkAlertSent(ctx context.Context, id int64) error
	TelegramChatIDs(ctx context.Context) ([]string, error)
}

// Commander diimplementasikan Peran 3 di internal/mqtt dan dipakai Peran 4.
type Commander interface {
	SetMode(ctx context.Context, nodeID string, mode Mode) error
	SendActuators(ctx context.Context, nodeID string, a Actuators) error
}

// Broadcaster diimplementasikan Peran 4 di internal/ws dan dipakai Peran 3.
type Broadcaster interface {
	Publish(event any)
}

// Event WebSocket. Bentuk JSON-nya ditetapkan di docs/kontrak.md.
type TelemetryEvent struct {
	Type   string    `json:"type"` // "telemetry"
	NodeID string    `json:"node_id"`
	TS     time.Time `json:"ts"`
	Temp   float64   `json:"temp"`
	NTU    float64   `json:"ntu"`
	Classes
}

type StateEvent struct {
	Type   string `json:"type"` // "state"
	NodeID string `json:"node_id"`
	ActuatorState
}

type NodeStatusEvent struct {
	Type   string `json:"type"` // "node_status"
	NodeID string `json:"node_id"`
	Online bool   `json:"online"`
}

type ModeEvent struct {
	Type   string `json:"type"` // "mode"
	NodeID string `json:"node_id"`
	Mode   Mode   `json:"mode"`
}

type AlertEvent struct {
	Type    string `json:"type"` // "alert"
	NodeID  string `json:"node_id"`
	Level   Class  `json:"level"`
	Message string `json:"message"`
}
```

Buat `backend/cmd/server/main.go`:

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/api/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Peran 3 dan 4 menyambungkan Store, Commander, Broadcaster, dan route di sini.
	log.Println("backend berjalan di :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
```

Setelah file dibuat, jalankan:

```bash
docker run --rm -v "${PWD}/backend:/app" -w /app golang:1.26-alpine go mod tidy
```

## Docker Compose

`deploy/mosquitto/dev.conf`:

```text
listener 1883
allow_anonymous true
```

`.env.example`:

```text
DOMAIN=localhost
POSTGRES_PASSWORD=dev_password
JWT_SECRET=ganti_dengan_minimal_32_karakter_acak
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123
MQTT_USERNAME=
MQTT_PASSWORD=
TELEGRAM_BOT_TOKEN=
```

`compose.yaml`:

```yaml
services:
  db:
    image: postgres:17-alpine
    environment:
      POSTGRES_DB: aquarium
      POSTGRES_USER: aquarium
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    ports:
      - "5432:5432"
    volumes:
      - db-data:/var/lib/postgresql/data
      - ./db/migrations:/docker-entrypoint-initdb.d:ro

  mosquitto:
    image: eclipse-mosquitto:2
    command: mosquitto -c /mosquitto/config/dev.conf
    ports:
      - "1883:1883"
    volumes:
      - ./deploy/mosquitto:/mosquitto/config:ro

  backend:
    image: golang:1.26-alpine
    working_dir: /app
    command: go run ./cmd/server
    env_file: .env
    environment:
      DATABASE_URL: postgres://aquarium:${POSTGRES_PASSWORD}@db:5432/aquarium
      MQTT_URL: tcp://mosquitto:1883
    ports:
      - "8080:8080"
    volumes:
      - ./backend:/app
      - go-cache:/go
    depends_on:
      - db
      - mosquitto

  simulator:
    profiles: [tools]
    build: ./tools/simulator
    environment:
      MQTT_HOST: mosquitto
    depends_on:
      - mosquitto

volumes:
  db-data:
  go-cache:
```

## README

Tulis ringkas: prasyarat, cara menjalankan, dan tautan ke `AGENTS.md` serta `MULAI-DI-SINI.md`.

## Pemeriksaan

```bash
copy .env.example .env
docker compose config -q
docker compose up -d
docker compose ps
curl http://localhost:8080/api/health
docker compose run --rm --no-deps backend go test ./...
cd frontend && npm run build
```

Di macOS atau Linux, ganti `copy` dengan `cp`. Fondasi selesai kalau health menjawab `{"status":"ok"}` dan build frontend berhasil. Build pertama backend lebih lama karena Docker mengunduh image dan dependensi.
