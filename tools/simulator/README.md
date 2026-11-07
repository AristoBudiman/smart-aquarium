# Simulator Akuarium

Simulator ini mewakili node akuarium virtual yang terhubung ke broker MQTT.
Digunakan untuk menguji backend dan integrasi sistem tanpa memerlukan alat fisik (NodeMCU ESP8266).

## Persyaratan

- Docker dan Docker Compose sudah terpasang.
- Broker Mosquitto sudah berjalan (via Docker Compose).

## Cara Menjalankan

### 1. Jalankan broker Mosquitto

```bash
docker compose up -d mosquitto
```

### 2. Build image simulator

```bash
docker compose build simulator
```

### 3. Jalankan simulator

```bash
docker compose run --rm simulator --nodes 3 --scenario normal --interval 5
```

### 4. Pantau pesan MQTT (di terminal lain)

```bash
docker compose exec mosquitto mosquitto_sub -t "aquarium/#" -v
```

## Argumen

| Argumen | Bawaan | Keterangan |
| --- | --- | --- |
| `--nodes` | `3` | Jumlah node virtual yang dijalankan. |
| `--scenario` | `normal` | Skenario awal sensor. Pilihan: `normal`, `panas`, `dingin`, `keruh`, `ganda`, `mati`. |
| `--interval` | `5` | Jeda waktu antarkirim telemetry dalam detik. |
| `--host` | env `MQTT_HOST` atau `localhost` | Alamat broker MQTT. |
| `--port` | `1883` | Port broker MQTT. |
| `--username` | env `MQTT_USERNAME` | Username MQTT (kosong di lokal). |
| `--password` | env `MQTT_PASSWORD` | Password MQTT (kosong di lokal). |

## Skenario

| Skenario | Suhu Awal | NTU Awal | Keterangan |
| --- | --- | --- | --- |
| `normal` | 27 °C | 10 | Kondisi ideal akuarium. |
| `panas` | 33 °C | 10 | Suhu tinggi, menguji kipas. |
| `dingin` | 22 °C | 10 | Suhu rendah, menguji pemanas. |
| `keruh` | 27 °C | 80 | Air keruh, menguji pompa filter. |
| `ganda` | 17 °C | 300 | Suhu sangat rendah dan air sangat keruh (WARNING/SEVERE). |
| `mati` | 27 °C | 10 | Simulasi node mati — koneksi diputus paksa agar LWT `offline` dikirim. |

## Topik MQTT

| Topik | Arah | Isi |
| --- | --- | --- |
| `aquarium/{node_id}/telemetry` | Simulator → Broker | `{"temp":27.0,"ntu":10.0}` |
| `aquarium/{node_id}/state` | Simulator → Broker | `{"heater":0,"fan":0,"pump":0,"failsafe":false}` |
| `aquarium/{node_id}/status` | Simulator → Broker | `online` / `offline` (retained) |
| `aquarium/{node_id}/cmd` | Broker → Simulator | `{"heater":1,"fan":0,"pump":0}` |
| `aquarium/{node_id}/config` | Broker → Simulator | `{"mode":"auto"}` |
| `server/heartbeat` | Broker → Simulator | `{"ts":1762171200}` |

## Contoh ID Node

- `--nodes 1` → `aq-sim001`
- `--nodes 3` → `aq-sim001`, `aq-sim002`, `aq-sim003`

## Pemeriksaan Dasar

```bash
# Terminal 1: jalankan simulator 3 node
docker compose run --rm simulator --nodes 3 --scenario normal --interval 2

# Terminal 2: pantau semua topik
docker compose exec mosquitto mosquitto_sub -t "aquarium/#" -v

# Terminal 3: kirim perintah uji ke node pertama
docker compose exec mosquitto mosquitto_pub -t "aquarium/aq-sim001/cmd" -q 1 -m '{"heater":1,"fan":0,"pump":0}'
```
