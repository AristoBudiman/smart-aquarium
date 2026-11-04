# Peran 3: ML dan logika backend

Target: satu telemetry dari simulator berakhir sebagai data tersimpan, keputusan aktuator, notifikasi, dan event WebSocket.

## Tugas

- [ ] Buat skrip ekspor model dan latih ulang.
- [ ] Tambahkan `model.go` dan tes kesamaan Python–Go.
- [ ] Buat `internal/decision` beserta unit test Tabel kontrak.
- [ ] Buat `internal/mqtt`, termasuk `Commander` dan heartbeat.
- [ ] Buat `internal/ingest`.
- [ ] Buat `internal/notify`.
- [ ] Sambungkan bagian milikmu di `cmd/server/main.go` lewat pull request bersama Peran 1.

## Model

Dataset ada di `ml/dataset/fishpond_dataset_multiclass_2153.csv`. Simpan skrip ini sebagai `ml/train_export.py`:

```python
"""Latih 3 model Random Forest, ekspor ke Go, dan tulis data uji.
Jalankan dari root repo: python ml/train_export.py
"""
import json
import re
import sys
from pathlib import Path

import m2cgen as m2c
import numpy as np
import pandas as pd
from sklearn.ensemble import RandomForestClassifier
from sklearn.metrics import accuracy_score, f1_score
from sklearn.model_selection import StratifiedKFold, cross_val_score, train_test_split

sys.setrecursionlimit(100000)

ROOT = Path(__file__).resolve().parents[1]
DATASET = ROOT / "ml/dataset/fishpond_dataset_multiclass_2153.csv"
OUT = ROOT / "backend/internal/model"
OUT.mkdir(parents=True, exist_ok=True)

KELAS = ["Normal", "Caution", "Warning", "Severe"]
LABEL = {k: i for i, k in enumerate(KELAS)}
FITUR = {
    "Suhu": ["temp_C"],
    "Kekeruhan": ["turbidity_NTU"],
    "Overall": ["temp_C", "turbidity_NTU"],
}

df = pd.read_csv(DATASET)[["temp_C", "turbidity_NTU", "class"]].dropna()
df["label"] = df["class"].map(LABEL)
latih, uji = train_test_split(df, test_size=0.2, stratify=df["label"], random_state=42)
cv = StratifiedKFold(5, shuffle=True, random_state=42)

models = {}
baris = ["| Model | Fitur | Akurasi uji | F1 macro | CV 5-fold |", "|---|---|---|---|---|"]
for nama, kolom in FITUR.items():
    model = RandomForestClassifier(n_estimators=50, max_depth=8, random_state=42, n_jobs=-1)
    model.fit(latih[kolom].values, latih["label"])
    pred = model.predict(uji[kolom].values)
    skor = cross_val_score(model, df[kolom].values, df["label"], cv=cv)
    baris.append(
        f"| {nama} | {', '.join(kolom)} | {accuracy_score(uji['label'], pred) * 100:.2f}% "
        f"| {f1_score(uji['label'], pred, average='macro') * 100:.2f}% "
        f"| {skor.mean() * 100:.2f} ± {skor.std() * 100:.2f}% |"
    )
    models[nama] = model

    kode = m2c.export_to_go(model, function_name=f"Score{nama}")
    kode = re.sub(r"\baddVectors\b", f"addVectors{nama}", kode)
    kode = re.sub(r"\bmulVectorNumber\b", f"mulVectorNumber{nama}", kode)
    (OUT / f"{nama.lower()}_gen.go").write_text("// Kode dibuat oleh ml/train_export.py. Jangan diedit.\npackage model\n\n" + kode)

(ROOT / "ml/hasil_pelatihan.md").write_text("\n".join(baris) + "\n")

rng = np.random.default_rng(0)
titik = []
for suhu, ntu in zip(rng.uniform(10, 40, 1000), rng.uniform(0, 700, 1000)):
    titik.append({
        "temp": float(suhu),
        "ntu": float(ntu),
        "suhu": int(models["Suhu"].predict([[suhu]])[0]),
        "kekeruhan": int(models["Kekeruhan"].predict([[ntu]])[0]),
        "overall": int(models["Overall"].predict([[suhu, ntu]])[0]),
    })
(OUT / "testdata.json").write_text(json.dumps(titik))
print("\n".join(baris))
```

Jalankan dari root repo:

```bash
docker run --rm -v "${PWD}:/repo" -w /repo python:3.13-slim sh -c "pip install -q scikit-learn==1.9.1 pandas==3.0.5 numpy==2.5.3 m2cgen==0.10.0 && python ml/train_export.py"
```

Hasilnya:

- `backend/internal/model/suhu_gen.go`, `kekeruhan_gen.go`, dan `overall_gen.go`.
- `backend/internal/model/testdata.json`.
- `ml/hasil_pelatihan.md` untuk Tabel 3.2.

Simpan `backend/internal/model/model.go`. Ganti `NAMA-AKUN` sesuai `go.mod`.

```go
package model

import "github.com/NAMA-AKUN/smart-aquarium/backend/internal/core"

var kelas = [...]core.Class{core.ClassNormal, core.ClassCaution, core.ClassWarning, core.ClassSevere}

func indeksTerbesar(v []float64) int {
	terbesar := 0
	for i := range v {
		if v[i] > v[terbesar] {
			terbesar = i
		}
	}
	return terbesar
}

// Predict menjalankan ketiga model Random Forest.
func Predict(temp, ntu float64) core.Classes {
	suhu := kelas[indeksTerbesar(ScoreSuhu([]float64{temp}))]
	kekeruhan := kelas[indeksTerbesar(ScoreKekeruhan([]float64{ntu}))]
	overall := kelas[indeksTerbesar(ScoreOverall([]float64{temp, ntu}))]
	return core.Classes{
		Temp:      suhu,
		Turbidity: kekeruhan,
		Overall:   overall,
		Status:    core.Worst(suhu, kekeruhan, overall),
	}
}
```

Simpan `backend/internal/model/model_test.go`:

```go
package model

import (
	"encoding/json"
	"os"
	"testing"
)

type titik struct {
	Temp      float64 `json:"temp"`
	NTU       float64 `json:"ntu"`
	Suhu      int     `json:"suhu"`
	Kekeruhan int     `json:"kekeruhan"`
	Overall   int     `json:"overall"`
}

func TestSamaDenganPython(t *testing.T) {
	data, err := os.ReadFile("testdata.json")
	if err != nil {
		t.Fatal(err)
	}
	var daftar []titik
	if err := json.Unmarshal(data, &daftar); err != nil {
		t.Fatal(err)
	}
	for _, x := range daftar {
		h := Predict(x.Temp, x.NTU)
		if h.Temp != kelas[x.Suhu] || h.Turbidity != kelas[x.Kekeruhan] || h.Overall != kelas[x.Overall] {
			t.Fatalf("berbeda pada suhu=%v ntu=%v: %+v", x.Temp, x.NTU, h)
		}
	}
}
```

Jalankan:

```bash
docker compose run --rm --no-deps backend go test ./internal/model
```

Hasil pelatihan saat dokumen ini dibuat:

| Model | Akurasi uji | F1 macro | CV 5-fold |
| --- | --- | --- | --- |
| Suhu | 99,07% | 99,02% | 99,63 ± 0,24% |
| Kekeruhan | 99,77% | 99,78% | 99,91 ± 0,11% |
| Overall | 100,00% | 100,00% | 100,00 ± 0,00% |

Kalau hasilnya berbeda, laporkan dan pakai angka baru untuk laporan.

## Keputusan aktuator

Buat fungsi murni tanpa database dan MQTT:

```go
func Decide(prev core.Actuators, c core.Classes, temp float64) core.Actuators
```

Unit test harus mencakup seluruh tabel aturan di kontrak, termasuk Caution yang mempertahankan status lama dan suhu tepat 27,48 °C.

## Alur telemetry

Urutan proses per pesan:

1. Baca `node_id` dari topik dan validasi JSON.
2. Tolak suhu di luar rentang kontrak dan buat alert sekali per perubahan kondisi.
3. Panggil `EnsureNode`. Kalau node baru, kirim `config` manual retained lewat `Commander.SetMode`.
4. Ambil `Snapshot`, karena `LastStatus` harus dibaca sebelum prediksi baru disimpan.
5. Simpan telemetry.
6. Jalankan `model.Predict`, lalu simpan prediksi.
7. Kalau mode Auto, jalankan `Decide`. Pakai `prev` nol jika `HasState` false.
8. Kirim `cmd` dan `LogActuators` dengan source `auto` kalau hasilnya berbeda, atau kalau `HasState` false.
9. Kalau status utama naik dan bukan NORMAL, simpan alert dan kirim notifikasi.
10. Kirim `TelemetryEvent` dan, bila ada, `AlertEvent`.

Topik `state` memanggil `SaveState` dan mengirim `StateEvent`. Topik `status` memanggil `SetOnline` dan mengirim `NodeStatusEvent`. Saat `failsafe` berubah menjadi true, catat `LogActuators` dengan source `failsafe`.

## MQTT

- Baca `MQTT_URL`, `MQTT_USERNAME`, dan `MQTT_PASSWORD` dari environment.
- Aktifkan auto-reconnect dan subscribe ulang di `OnConnect`.
- Subscribe `aquarium/+/telemetry`, `aquarium/+/state`, dan `aquarium/+/status`.
- Jangan memblokir callback Paho dengan proses panjang.
- `SetMode` menulis `config` QoS 1 retained. `SendActuators` menulis `cmd` QoS 1.
- Kirim heartbeat setiap 30 detik sampai context selesai.

## Telegram

- Kirim `sendMessage` di goroutine terpisah, dan panggil `MarkAlertSent` setelah berhasil.
- Kalau `TELEGRAM_BOT_TOKEN` kosong, tulis ke log dan lanjutkan.
- Jangan menulis token ke log.
- Contoh pesan: `⚠️ WARNING – Akuarium 1: suhu 33,1 °C (WARNING), kekeruhan 18 NTU (NORMAL). Kipas dinyalakan otomatis.`
- Level SEVERE memakai awalan `🚨 DARURAT`.

## Pemeriksaan

```bash
docker compose run --rm --no-deps backend go vet ./...
docker compose run --rm --no-deps backend go test ./...
docker compose up -d
docker compose exec mosquitto mosquitto_pub -t "aquarium/aq-test01/telemetry" -m '{"temp":33,"ntu":10}'
docker compose logs --tail=50 backend
```

Node `aq-test01` harus terdaftar, dan prediksi status utamanya harus `WARNING`.
