# Panduan Deploy Produksi (VM Azure)

Dokumen ini menjelaskan langkah-langkah penyiapan dan pembaruan layanan produksi pada VM menggunakan Docker Compose, Caddy, DuckDNS, dan Mosquitto dengan autentikasi serta ACL.

## 1. Persiapan Variabel Lingkungan (`.env`)

Berkas `.env` produksi diletakkan di **root repositori** pada VM (bukan di dalam folder `deploy/`). Salin dari `.env.example` lalu sesuaikan nilainya:

- `DOMAIN=nama.duckdns.org` (sesuai domain DuckDNS yang telah diarahkan ke IP publik VM)
- `POSTGRES_PASSWORD` diisi password database produksi yang kuat
- `JWT_SECRET` diisi minimal 32 karakter acak
- `ADMIN_USERNAME` dan `ADMIN_PASSWORD` diisi kredensial admin awal
- `MQTT_USERNAME=backend`
- `MQTT_PASSWORD` diisi sama dengan nilai `PASS_BACKEND` yang dibuat pada langkah di bawah
- `TELEGRAM_BOT_TOKEN` diisi token bot Telegram (atau kosongkan jika hanya dicatat di log)

## 2. Pembuatan File Password Mosquitto & Deploy di VM

Jalankan perintah berikut dari **root repositori** di terminal VM. Ganti `PASS_BACKEND` dan `PASS_NODE` hanya di terminal VM (jangan menuliskan password asli ke dalam kode atau Git):

```bash
docker run --rm -v "$PWD/deploy/mosquitto:/mosquitto/config" eclipse-mosquitto:2 sh -c "touch /mosquitto/config/passwd && mosquitto_passwd -b /mosquitto/config/passwd backend PASS_BACKEND && mosquitto_passwd -b /mosquitto/config/passwd node PASS_NODE && chown 1883:1883 /mosquitto/config/passwd /mosquitto/config/acl && chmod 600 /mosquitto/config/passwd /mosquitto/config/acl"
sh deploy/deploy.sh
```

Untuk pembaruan berikutnya setelah perubahan kode masuk ke repositori, cukup jalankan:

```bash
sh deploy/deploy.sh
```

## 3. Pemeriksaan Setelah VM Siap

Pastikan seluruh poin berikut terpenuhi setelah layanan berjalan di VM:

- [ ] `https://DOMAIN/api/health` menjawab `{"status":"ok"}`.
- [ ] `mosquitto_pub` dengan password benar berhasil, dan password salah ditolak.
- [ ] Port `5432` dan `8080` tidak bisa diakses dari internet.
- [ ] Layanan menyala lagi secara otomatis setelah VM restart.
- [ ] Publish ke topik node lain muncul sebagai `Denied` di `docker compose --env-file .env -f deploy/docker-compose.yml logs mosquitto`.

