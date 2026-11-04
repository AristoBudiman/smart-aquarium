# AI Smart Aquarium

Sistem pemantauan dan pengendalian akuarium pintar berbasis IoT (multi-node) dengan klasifikasi kondisi menggunakan tiga model Random Forest (suhu, kekeruhan, dan kondisi keseluruhan) serta pengendalian otomatis pemanas, kipas, dan pompa filter.

## Prasyarat

- **Git**
- **Docker Desktop** (untuk menjalankan PostgreSQL, Mosquitto, Backend Go, dan Simulator)
- **Node.js 22.18+** (untuk pengembangan frontend Vue 3)
- **Arduino IDE** (untuk pengembangan firmware ESP8266 pada Peran 1)

> Catatan: Go tidak perlu dipasang di laptop karena seluruh proses kompilasi dan pengujian backend dijalankan melalui container Docker.

## Cara Menjalankan

1. Salin berkas variabel lingkungan:
   ```powershell
   copy .env.example .env
   ```
   *(Gunakan `cp .env.example .env` pada macOS/Linux).*

2. Jalankan layanan utama menggunakan Docker Compose:
   ```powershell
   docker compose up -d
   ```

3. Periksa status layanan dan endpoint kesehatan backend:
   ```powershell
   docker compose ps
   curl http://localhost:8080/api/health
   ```

4. Untuk menjalankan server pengembangan frontend:
   ```powershell
   cd frontend
   npm install
   npm run dev
   ```

## Panduan Lebih Lanjut

- Aturan kerja, pembagian peran, dan perintah harian: [AGENTS.md](./AGENTS.md)
- Panduan memulai untuk setiap anggota tim: [MULAI-DI-SINI.md](./MULAI-DI-SINI.md)
- Kontrak sistem (MQTT, REST API, WebSocket, Database): [docs/kontrak.md](./docs/kontrak.md)

