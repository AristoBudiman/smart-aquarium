# Mulai di sini

Panduan ini untuk manusia. Agent membaca `AGENTS.md`, `docs/kontrak.md`, dan instruksi perannya.

## Aristo: buat repo pertama kali

1. Buat repo private `smart-aquarium` di GitHub tanpa README, lalu undang keempat anggota.
2. Clone repo tersebut:

    ```bash
    git clone https://github.com/NAMA-AKUN/smart-aquarium.git
    cd smart-aquarium
    ```

3. Ekstrak isi zip paket langsung ke folder `smart-aquarium`. File `AGENTS.md` harus berada di root repo.
4. Push paket sebagai commit pertama:

    ```bash
    git add .
    git commit -m "docs: tambah panduan agent dan kontrak"
    git branch -M main
    git push -u origin main
    ```

5. Buka folder `smart-aquarium` di Antigravity, nyalakan Docker Desktop, lalu pakai prompt Peran 1A di bawah.
6. Setelah fondasi lolos pemeriksaan dan masuk `main`, beri tahu anggota lain.

## Anggota lain

Mulai setelah Aristo mengumumkan fondasi sudah masuk `main`. Peran 3 boleh mulai lebih awal dari langkah pelatihan model.

```bash
git clone https://github.com/NAMA-AKUN/smart-aquarium.git
cd smart-aquarium
git checkout -b feat/NAMA-BAGIAN
```

Isi `NAMA-BAGIAN` dengan `sim`, `ml`, `api`, atau `web`. Salin `.env.example` menjadi `.env`, buka folder di Antigravity, lalu pakai prompt peran masing-masing.

## Yang perlu menyala

| Peran | Aplikasi |
| --- | --- |
| 1, 2, 3, 4 | Git, Docker Desktop |
| 5 | Git, Node 22.18 ke atas; Docker Desktop saat memakai backend asli |
| 1 untuk firmware | Arduino IDE dengan NodeMCU tersambung USB |

Docker Desktop saja tidak cukup untuk menyelesaikan tugas. Setiap orang tetap perlu menyetujui perintah agent, mengecek hasil, membuat commit, dan membuka pull request.

## Prompt untuk agent

Kirim satu prompt pada awal sesi. Pakai sesi baru untuk tugas besar berikutnya.

| Peran | Prompt |
| --- | --- |
| 1A Fondasi | `Baca AGENTS.md, docs/kontrak.md, dan docs/agent/peran-1a-fondasi.md. Kamu membantu Peran 1. Buat fondasi repo sesuai instruksi, tunjukkan rencana dulu, lalu jalankan semua pemeriksaannya.` |
| 1B Infrastruktur | `Baca AGENTS.md, docs/kontrak.md, dan docs/agent/peran-1b-infrastruktur.md. Kamu membantu Peran 1. Buat file deploy, tetapi jangan menjalankan perintah di VM kecuali aku minta.` |
| 1C Firmware | `Baca AGENTS.md, docs/kontrak.md, dan docs/agent/peran-1c-firmware.md. Kamu membantu Peran 1. Kerjakan tahap firmware yang aku sebutkan, lalu tulis cara mengujinya di alat.` |
| 2 Simulator | `Baca AGENTS.md, docs/kontrak.md, dan docs/agent/peran-2-simulator-uji.md. Kamu membantu Peran 2. Kerjakan tugas pertama yang belum selesai dan jalankan pemeriksaannya.` |
| 3 ML dan logika | `Baca AGENTS.md, docs/kontrak.md, dan docs/agent/peran-3-ml-logika.md. Kamu membantu Peran 3. Kerjakan tugas pertama yang belum selesai dan jalankan pemeriksaannya.` |
| 4 API | `Baca AGENTS.md, docs/kontrak.md, dan docs/agent/peran-4-api.md. Kamu membantu Peran 4. Kerjakan tugas pertama yang belum selesai dan jalankan pemeriksaannya.` |
| 5 Frontend | `Baca AGENTS.md, docs/kontrak.md, dan docs/agent/peran-5-frontend.md. Kamu membantu Peran 5. Kerjakan tugas pertama yang belum selesai dan jalankan pemeriksaannya.` |

## Rutinitas harian

```bash
git checkout main
git pull
git checkout feat/NAMA-BAGIAN
git merge main
```

Setelah satu tugas selesai, commit lalu push:

```bash
git add NAMA-FILE
git commit -m "feat: jelaskan perubahan"
git push -u origin feat/NAMA-BAGIAN
```

Buka pull request ke `main`, lalu minta satu anggota mereviu.
