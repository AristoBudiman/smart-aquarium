# Peran 5: frontend

Target: enam halaman berfungsi di laptop dan ponsel. Mulai dengan data palsu, lalu ganti ke API asli tanpa mengubah komponen halaman.

## Tugas

- [ ] Buat tipe data dan lapisan `src/api/`.
- [ ] Buat data palsu dengan `VITE_MOCK=1`.
- [ ] Buat store, router, guard login, dan layout.
- [ ] Buat Login, Daftar Node, Detail Node, Grafik Riwayat, Log & Notifikasi, dan Pengaturan.
- [ ] Sambungkan WebSocket.
- [ ] Rapikan tampilan ponsel 390×844 dan desktop.
- [ ] Ambil gambar untuk Gambar 4.7 sampai 4.9.

## Struktur

```text
frontend/src/
├─ api/          client.ts, mock.ts, nodes.ts, auth.ts, ws.ts
├─ components/
├─ router/
├─ stores/
├─ types/
└─ views/
```

Komponen tidak boleh memanggil `axios` langsung. Saat `VITE_MOCK=1`, lapisan `api` memakai data palsu dan WebSocket palsu. Sediakan minimal tiga node: NORMAL, WARNING, dan offline.

File `.env.development`:

```text
VITE_MOCK=1
```

Hapus baris tersebut atau isi `0` saat memakai backend asli. Proxy `/api` dan `/ws` sudah disiapkan di Vite.

## Halaman

| Halaman | Isi |
| --- | --- |
| Login | Username, password, pesan error |
| Daftar Node | Kartu berisi nama, online/offline, suhu, NTU, status utama, waktu terakhir |
| Detail Node | Gauge suhu 0–40 °C, gauge kekeruhan 0–200 NTU, kelas suhu, kekeruhan, overall, status utama, sakelar mode, tiga sakelar aktuator, tanda fail-safe |
| Grafik Riwayat | Grafik suhu dan NTU, pilihan 1 jam, 24 jam, 7 hari, 30 hari |
| Log & Notifikasi | Riwayat aktuator dan alert |
| Pengaturan | Ubah nama node dan isi Telegram chat ID |

Warna status: NORMAL hijau, CAUTION kuning, WARNING oranye, SEVERE merah. Selain warna, selalu tampilkan teks status.

## Perilaku

- Simpan token di `localStorage`, dan arahkan ke Login saat menerima `401`.
- Tampilkan waktu dalam WIB dengan zona `Asia/Jakarta`.
- Sakelar aktuator nonaktif saat mode Auto.
- Setelah `PUT actuators`, tampilkan status menunggu sampai event `state` datang. Jangan menganggap relay sudah berubah.
- Tampilkan `409` sebagai pesan bahwa node harus berada di mode Manual.
- WebSocket tersambung ke `/ws?token=JWT` dan tersambung ulang otomatis dengan jeda bertahap.
- Tampilkan state kosong, loading, dan error pada setiap halaman.
- Impor ECharts dari `echarts/core` dan daftarkan komponen yang dipakai saja.
- Teks antarmuka memakai bahasa Indonesia.

## Pemeriksaan

```bash
cd frontend
npm run type-check
npm run build
```

Cek manual:

- [ ] Semua halaman berjalan dengan `VITE_MOCK=1`.
- [ ] Refresh di halaman Detail Node tidak mengeluarkan pengguna.
- [ ] Tidak ada scroll horizontal pada lebar 390 px.
- [ ] Data WebSocket mengubah kartu dan detail tanpa refresh.
