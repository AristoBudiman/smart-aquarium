# Skenario Uji

Batas berikut dari model 50 pohon, kedalaman 8. Nilainya bisa sedikit bergeser setelah dilatih ulang.

| Nilai | Kelas | Mode Auto |
| --- | --- | --- |
| Suhu di bawah 20,25 °C | SEVERE | Pemanas ON, alert darurat |
| Suhu 20,25–23 °C | WARNING | Pemanas ON, alert |
| Suhu 23–25 °C | CAUTION | Aktuator tetap, alert |
| Suhu 25–30 °C | NORMAL | Pemanas dan kipas OFF |
| Suhu 30–32,25 °C | CAUTION | Aktuator tetap, alert |
| Suhu 32,25–34 °C | WARNING | Kipas ON, alert |
| Suhu di atas 34 °C | SEVERE | Kipas ON, alert darurat |
| Kekeruhan di bawah 25 NTU | NORMAL | Pompa OFF |
| Kekeruhan 25–51 NTU | CAUTION | Pompa tetap, alert |
| Kekeruhan 51–151 NTU | WARNING | Pompa ON, alert |
| Kekeruhan 151 NTU ke atas | SEVERE | Pompa ON, alert darurat |
