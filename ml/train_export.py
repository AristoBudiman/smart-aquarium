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

