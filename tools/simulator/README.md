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

### 3. Jalankan unit test simulator (tanpa broker)

```bash
docker compose run --rm --no-deps --entrypoint python simulator -m unittest discover -v
```

### 4. Jalankan simulator

```bash
docker compose run --rm simulator --nodes 3 --scenario normal --interval 5
```

### 5. Pantau pesan MQTT (di terminal lain)

```bash
docker compose exec mosquitto mosquitto_sub -t "aquarium/#" -v
```

## Argumen dan Variabel Lingkungan

Argumen baris perintah selalu mengalahkan variabel lingkungan:

| Argumen | Variabel Lingkungan | Bawaan | Keterangan |
| --- | --- | --- | --- |
| `--nodes` | `SIM_NODES` | `3` | Jumlah node virtual yang dijalankan (`aq-sim001`, `aq-sim002`, dst.). |
| `--scenario` | `SIM_SCENARIO` | `normal` | Skenario awal sensor: `normal`, `panas`, `dingin`, `keruh`, `ganda`, `mati`. |
| `--interval` | `SIM_INTERVAL` | `5` | Jeda waktu antarkirim telemetry dan state berkala dalam detik. |
| `--host` | `MQTT_HOST` | `localhost` | Alamat broker MQTT (`mosquitto` di Compose). |
| `--port` | `MQTT_PORT` | `1883` | Port broker MQTT. |
| `--username` | `MQTT_USERNAME` | `""` | Username MQTT (kosong di lokal, `node` di broker VM produksi). |
| `--password` | `MQTT_PASSWORD` | `""` | Password MQTT (kosong di lokal). |

## Skenario dan Perubahan Sensor

| Skenario | Suhu Awal | NTU Awal | Keterangan |
| --- | --- | --- | --- |
| `normal` | 27 °C | 10 | Kondisi ideal akuarium. |
| `panas` | 33 °C | 10 | Suhu tinggi, menguji kipas. |
| `dingin` | 22 °C | 10 | Suhu rendah, menguji pemanas. |
| `keruh` | 27 °C | 80 | Air keruh, menguji pompa filter. |
| `ganda` | 17 °C | 300 | Suhu sangat rendah dan air sangat keruh (WARNING/SEVERE). |
| `mati` | 27 °C | 10 | Simulasi node mati — koneksi TCP diputus tanpa paket `DISCONNECT` agar LWT `offline` dikirim oleh broker. |

Saat aktuator menyala, nilai sensor berubah setiap siklus interval:
- **Pemanas (`heater=1`, `fan=0`)**: suhu naik `+0.3 °C` per siklus.
- **Kipas (`fan=1`, `heater=0`)**: suhu turun `-0.3 °C` per siklus.
- **Pompa (`pump=1`)**: kekeruhan turun `-5.0 NTU` per siklus (minimal `0.0 NTU`).

## Topik MQTT

| Topik | Arah | QoS | Retained | Isi |
| --- | --- | --- | --- | --- |
| `aquarium/{node_id}/telemetry` | Simulator → Broker | 0 | Tidak | `{"temp":27.0,"ntu":10.0}` |
| `aquarium/{node_id}/state` | Simulator → Broker | 0 | Tidak | `{"heater":0,"fan":0,"pump":0,"failsafe":false}` |
| `aquarium/{node_id}/status` | Simulator → Broker | 0 | Ya | `online` |
| `aquarium/{node_id}/status` | LWT Broker | 1 | Ya | `offline` |
| `aquarium/{node_id}/cmd` | Broker → Simulator | 1 | Tidak | `{"heater":0,"fan":1,"pump":0}` |
| `aquarium/{node_id}/config` | Broker → Simulator | 1 | Ya | `{"mode":"auto"}` atau `{"mode":"manual"}` |
| `server/heartbeat` | Broker → Simulator | 0 | Tidak | `{"ts":1762171200}` |

## Pemeriksaan Dasar

```bash
# Terminal 1: jalankan simulator 3 node
docker compose run --rm simulator --nodes 3 --scenario panas --interval 2

# Terminal 2: pantau semua topik
docker compose exec mosquitto mosquitto_sub -t "aquarium/#" -v

# Terminal 3: kirim perintah uji ke node pertama
docker compose exec mosquitto mosquitto_pub -t "aquarium/aq-sim001/cmd" -q 1 -m '{"heater":0,"fan":1,"pump":0}'
```
