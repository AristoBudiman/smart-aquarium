# Peran 2: simulator, pengujian, dokumentasi

Target: orang lain bisa menguji sistem tanpa alat fisik. Versi awal simulator harus siap paling lambat Kamis 5 November.

## Tugas

- [ ] Buat simulator Python di `tools/simulator/`.
- [ ] Uji simulator terhadap broker lokal.
- [ ] Buat `docs/laporan/hasil-uji.md` dari daftar uji.
- [ ] Bantu pengujian multi-node pada 8–9 November.
- [ ] Buat Gambar 4.3 dan 4.4.
- [ ] Lengkapi `[isi ...]` di laporan.
- [ ] Satukan laporan dan PPT paling lambat Selasa 10 November.

## Spesifikasi simulator

File:

```text
tools/simulator/main.py
tools/simulator/requirements.txt
tools/simulator/Dockerfile
tools/simulator/README.md
```

Pakai `paho-mqtt` versi 2 dengan `CallbackAPIVersion.VERSION2`. `Dockerfile`:

```dockerfile
FROM python:3.13-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install -r requirements.txt
COPY main.py .
ENTRYPOINT ["python", "main.py"]
```

Argumen:

| Argumen | Bawaan | Fungsi |
| --- | --- | --- |
| `--nodes` | 3 | Jumlah node |
| `--scenario` | `normal` | `normal`, `panas`, `dingin`, `keruh`, `ganda`, `mati` |
| `--interval` | 5 | Detik antarpesan |
| `--host` | env `MQTT_HOST` | Alamat broker |
| `--port` | 1883 | Port broker |
| `--username`, `--password` | env | Kredensial, kosong di lokal |

Perilaku:

- Client ID dan `node_id` memakai `aq-sim001`, `aq-sim002`, dan seterusnya.
- Pasang LWT `offline`, kirim `online`, lalu subscribe `cmd`, `config`, dan `server/heartbeat`.
- Setiap `cmd` diterima, ubah status aktuator dan kirim `state` persis seperti firmware.
- Saat aktuator menyala, nilai bergerak pelan: pemanas menaikkan suhu, kipas menurunkan suhu, pompa menurunkan NTU.
- Skenario `mati` memutus koneksi secara paksa supaya LWT terkirim.

Nilai skenario:

| Skenario | Suhu awal | NTU awal |
| --- | --- | --- |
| `normal` | 27 °C | 10 |
| `panas` | 33 °C | 10 |
| `dingin` | 22 °C | 10 |
| `keruh` | 27 °C | 80 |
| `ganda` | 17 °C | 300 |

## Pemeriksaan

```bash
docker compose up -d mosquitto
docker compose build simulator
docker compose run --rm simulator --nodes 3 --scenario normal --interval 2
```

Di terminal lain:

```bash
docker compose exec mosquitto mosquitto_sub -t "aquarium/#" -v
```

Telemetry dari tiga node harus muncul. Kirim satu perintah uji:

```bash
docker compose exec mosquitto mosquitto_pub -t "aquarium/aq-sim001/cmd" -q 1 -m '{"heater":1,"fan":0,"pump":0}'
```

Node `aq-sim001` harus membalas `state`.

## Daftar uji

Buat tabel berkolom skenario, hasil, tanggal, penguji, dan bukti:

- Batas kelas suhu dan kekeruhan pada `docs/referensi/skenario-uji.md`.
- Suhu rendah dan kekeruhan tinggi menyalakan pemanas dan pompa.
- Caution mempertahankan status aktuator.
- Mode Manual menggerakkan relay, dan model tidak mengubahnya.
- Node baru muncul sebagai Manual.
- Satu node fisik dan tiga node simulasi tampil bersamaan.
- Node mati berubah offline.
- Node restart saat Auto kembali mengikuti keputusan backend.
- DS18B20 lepas tidak menghasilkan data palsu.
- Backend mati lebih dari 3 menit memicu fail-safe.
- Restart VM mengembalikan layanan dan data.
