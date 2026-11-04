# Peran 1B: infrastruktur

Agent membuat file deploy di branch `feat/infra`. Pembuatan akun Azure, VM, DuckDNS, password, dan login SSH dilakukan pengguna. Jangan meminta pengguna menempelkan password atau token ke chat.

## Tugas

- [ ] Buat `backend/Dockerfile` dan `frontend/Dockerfile`.
- [ ] Buat `deploy/docker-compose.yml`, `deploy/Caddyfile`, `deploy/mosquitto/mosquitto.conf`, dan `deploy/mosquitto/acl`.
- [ ] Buat `deploy/deploy.sh`.
- [ ] Jalankan `docker compose --env-file .env -f deploy/docker-compose.yml config -q`.
- [ ] Tulis panduan singkat di `deploy/README.md`.

## File

`backend/Dockerfile`:

```dockerfile
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata
COPY --from=build /server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

`frontend/Dockerfile`:

```dockerfile
FROM node:22-alpine AS build
WORKDIR /src
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build-only

FROM caddy:2-alpine
COPY --from=build /src/dist /srv
```

Tambahkan `.dockerignore` berisi `node_modules` dan `dist` di `frontend/`.

`deploy/docker-compose.yml`:

```yaml
services:
  db:
    image: postgres:17-alpine
    restart: unless-stopped
    environment:
      POSTGRES_DB: aquarium
      POSTGRES_USER: aquarium
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    volumes:
      - db-data:/var/lib/postgresql/data
      - ../db/migrations:/docker-entrypoint-initdb.d:ro

  mosquitto:
    image: eclipse-mosquitto:2
    restart: unless-stopped
    ports:
      - "1883:1883"
    volumes:
      - ./mosquitto:/mosquitto/config:ro
      - mosquitto-data:/mosquitto/data

  backend:
    build: ../backend
    restart: unless-stopped
    environment:
      DATABASE_URL: postgres://aquarium:${POSTGRES_PASSWORD}@db:5432/aquarium
      MQTT_URL: tcp://mosquitto:1883
      MQTT_USERNAME: ${MQTT_USERNAME}
      MQTT_PASSWORD: ${MQTT_PASSWORD}
      JWT_SECRET: ${JWT_SECRET}
      ADMIN_USERNAME: ${ADMIN_USERNAME}
      ADMIN_PASSWORD: ${ADMIN_PASSWORD}
      TELEGRAM_BOT_TOKEN: ${TELEGRAM_BOT_TOKEN}
    depends_on:
      - db
      - mosquitto

  caddy:
    build: ../frontend
    restart: unless-stopped
    environment:
      DOMAIN: ${DOMAIN}
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - caddy-data:/data
      - caddy-config:/config
    depends_on:
      - backend

volumes:
  db-data:
  mosquitto-data:
  caddy-data:
  caddy-config:
```

`deploy/Caddyfile`:

```text
{$DOMAIN} {
    handle /api/* {
        reverse_proxy backend:8080
    }
    handle /ws {
        reverse_proxy backend:8080
    }
    handle {
        root * /srv
        try_files {path} /index.html
        file_server
    }
}
```

`deploy/mosquitto/mosquitto.conf`:

```text
listener 1883
allow_anonymous false
password_file /mosquitto/config/passwd
acl_file /mosquitto/config/acl
persistence true
persistence_location /mosquitto/data/
log_dest stdout
log_type all
```

`deploy/mosquitto/acl`. Baris `topic` harus berada di bawah `user`:

```text
user backend
topic readwrite aquarium/#
topic write server/heartbeat

user node
topic read server/heartbeat
pattern readwrite aquarium/%c/#
```

`deploy/deploy.sh`:

```bash
#!/bin/sh
set -e
git pull
docker compose --env-file .env -f deploy/docker-compose.yml build backend
docker compose --env-file .env -f deploy/docker-compose.yml build caddy
docker compose --env-file .env -f deploy/docker-compose.yml up -d
docker image prune -f
```

## Langkah manusia di VM

Tulis langkah ini di `deploy/README.md`. Ganti `PASS_BACKEND` dan `PASS_NODE` hanya di terminal VM.

```bash
docker run --rm -v "$PWD/deploy/mosquitto:/mosquitto/config" eclipse-mosquitto:2 sh -c "touch /mosquitto/config/passwd && mosquitto_passwd -b /mosquitto/config/passwd backend PASS_BACKEND && mosquitto_passwd -b /mosquitto/config/passwd node PASS_NODE && chown 1883:1883 /mosquitto/config/passwd /mosquitto/config/acl && chmod 600 /mosquitto/config/passwd /mosquitto/config/acl"
sh deploy/deploy.sh
```

File `.env` produksi berada di root repo pada VM. Gunakan `DOMAIN=nama.duckdns.org`, `MQTT_USERNAME=backend`, dan `MQTT_PASSWORD` yang sama dengan `PASS_BACKEND`.

## Pemeriksaan setelah VM siap

- [ ] `https://DOMAIN/api/health` menjawab `{"status":"ok"}`.
- [ ] `mosquitto_pub` dengan password benar berhasil, dan password salah ditolak.
- [ ] Port 5432 dan 8080 tidak bisa diakses dari internet.
- [ ] Layanan menyala lagi setelah VM restart.
- [ ] Publish ke topik node lain muncul sebagai `Denied` di `docker compose logs mosquitto`.
