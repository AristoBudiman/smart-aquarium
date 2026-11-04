# AI Smart Aquarium — Laporan Bab 1–4 (Revisi Arsitektur Mandiri)

**Keterangan tanda:**
- **[DIUBAH]** = isi bagian ini berubah dari laporan terakhir
- **[BARU]** = bagian ini belum ada di laporan terakhir
- Tanpa tanda = sama seperti laporan terakhir, tidak perlu diubah
- `[isi ...]` = perlu dilengkapi kelompok

Daftar pustaka ada di akhir dokumen. Hasil pelatihan model di 3.1.4 masih berupa gambaran awal.

---

# BAB 1 IDENTIFIKASI MASALAH & LATAR BELAKANG

## 1.1 Latar Belakang Permasalahan [DIUBAH sedikit, paragraf terakhir]

Pemeliharaan ekosistem akuarium memerlukan pemantauan dan pengelolaan parameter kualitas air secara konsisten untuk menjaga kelangsungan hidup serta kesehatan biota di dalamnya. Parameter fisik dan kimia seperti suhu air dan tingkat kekeruhan (*turbidity*) merupakan indikator vital yang saling memengaruhi kestabilan lingkungan akuarium. Namun, metode pemantauan konvensional yang masih bergantung pada pengecekan manual secara berkala memiliki keterbatasan signifikan. Pemantauan manual tidak dapat dilakukan secara *real-time* dan terus-menerus selama 24 jam, sehingga fluktuasi atau anomali kondisi air sering kali terlambat terdeteksi oleh pemilik.

Keterlambatan dalam mendeteksi perubahan parameter air ini berpotensi menimbulkan dampak buruk bagi ekosistem akuarium. Sebagai contoh, lonjakan suhu yang mendadak dapat merusak keseimbangan ekologis serta memicu stres pada ikan. Selain itu, tingkat kekeruhan air yang tinggi menandakan penumpukan sisa pakan atau kotoran yang berisiko menurunkan kadar oksigen terlarut. Jika kondisi-kondisi anomali ini tidak segera ditangani melalui tindakan korektif—seperti sirkulasi air atau aerasi tambahan—risiko kematian biota akuarium akan meningkat secara signifikan.

Masalah utama dalam evaluasi kondisi akuarium terletak pada sifat keterkaitan antar-parameter (*multi-parameter dependency*). Status kesehatan lingkungan akuarium tidak dapat disimpulkan secara akurat hanya berdasarkan satu variabel tunggal. Sebaliknya, kombinasi dan interaksi dari beberapa parameter secara simultan memberikan gambaran yang jauh lebih holistik dan presisi mengenai kondisi akuarium sebenarnya. Sebagai contoh, suhu tinggi yang disertai kekeruhan meningkat membutuhkan respons penanganan yang berbeda dibandingkan dengan suhu tinggi pada kondisi air yang jernih. Oleh karena itu, pendekatan evaluasi berbasis variabel tunggal tidak lagi memadai untuk pemeliharaan akuarium secara modern.

Berdasarkan permasalahan tersebut, dikembangkanlah sebuah solusi inovatif berupa Smart Aquarium berbasis Internet of Things (IoT) dan Kecerdasan Buatan (*Machine Learning*). Sistem ini dirancang untuk melakukan akuisisi data parameter air (suhu, kekeruhan) secara otomatis dan berkelanjutan. Data *telemetry* tersebut ditransmisikan **ke server *cloud* yang dibangun secara mandiri** untuk diklasifikasikan menggunakan algoritma *Machine Learning* guna menentukan status kondisi akuarium secara cerdas (Normal, Caution, Warning, dan Severe) sekaligus menentukan tindakan aktuator yang perlu dijalankan. Melalui integrasi **antarmuka website** dan kontrol aktuator jarak jauh (*remote & automatic control*), sistem ini memungkinkan pemilik akuarium melakukan pemantauan *real-time* sekaligus mengambil tindakan pencegahan secara cepat dan tepat dari mana saja. **Sistem juga dirancang agar dapat menampung lebih dari satu akuarium (*multi-node*) tanpa perubahan pada sisi server.**

## 1.2 Rumusan Masalah [DIUBAH, tambah poin d]

Dengan berlandaskan pada latar belakang yang telah dipaparkan, rumusan masalah dari tugas besar / proyek ini adalah sebagai berikut.

a. Bagaimana cara memantau kondisi akuarium secara *real-time* menggunakan sensor?
b. Respon otomatis apa yang diperlukan ketika kondisi akuarium terdeteksi tidak normal?
c. Bagaimana cara agar pengguna dapat memantau dan mengontrol akuarium dari jarak jauh?
d. **Bagaimana merancang sistem agar dapat menampung penambahan node akuarium baru tanpa bergantung pada platform IoT pihak ketiga?**

## 1.3 Potensi Solusi IoT [DIUBAH pada poin b dan c, tambah poin f]

Pendekatan Internet of Things (IoT) menjadi solusi yang sangat tepat dan relevan untuk mengatasi kompleksitas pemantauan serta pengelolaan ekosistem akuarium. Keterbatasan pemantauan konvensional berbasis manusia—yang bersifat manual, episodik, dan rentan terhadap keterlambatan respon—dapat diatasi secara efektif melalui otomasi dan kontinuitas yang ditawarkan oleh teknologi IoT.

Berikut adalah potensi dan keunggulan utama dari penerapan solusi berbasis IoT dalam menyelesaikan permasalahan pemeliharaan akuarium:

**a. Pemantauan Real-Time dan Kontinu (*Continuous Real-Time Monitoring*)**
Integrasi mikrokontroler dengan berbagai sensor (seperti sensor suhu DS18B20 dan sensor kekeruhan air) memungkinkan akuisisi data parameter air dilakukan secara *real-time* dan terus-menerus selama 24/7 tanpa jeda. Hal ini menghilangkan ketergantungan pada pengecekan fisik secara berkala oleh pemilik akuarium.

**b. Deteksi Dini dan Respon Cepat terhadap Anomali (*Early Warning System*) [DIUBAH]**
Data *telemetry* sensor yang ditransmisikan secara rutin ke server *cloud* melalui **protokol MQTT** memungkinkan analisis kondisi air dilakukan secara instan. Ketika terjadi fluktuasi parameter yang berada di luar batas normal (misalnya lonjakan suhu mendadak), sistem dapat memberikan notifikasi peringatan secara cepat kepada pengguna **melalui website dan bot Telegram** sebelum kondisi memburuk dan membahayakan biota akuarium.

**c. Kemudahan Akses dan Kontrol Jarak Jauh (*Remote Control & Accessibility*) [DIUBAH]**
Melalui **website responsif yang dapat dibuka dari komputer maupun ponsel**, pengguna dapat memantau grafik tren kualitas air dari mana saja dan kapan saja. Selain pemantauan, teknologi IoT memungkinkan eksekusi kontrol jarak jauh terhadap aktuator—seperti menyalakan atau mematikan pemanas, kipas, dan pompa filter secara manual dari antarmuka pengguna saat pengguna berada jauh dari lokasi akuarium.

**d. Integrasi Analisis Cerdas Berbasis Machine Learning (*Smart Decision Support*)**
Penerapan IoT menyediakan fondasi alur data (*data pipeline*) yang terstruktur untuk diintegrasikan dengan algoritma kecerdasan buatan (*Machine Learning*) di platform *cloud*. Pengolahan parameter secara simultan memungkinkan pengklasifikasian kondisi akuarium ke dalam empat tingkat (Normal, Caution, Warning, Severe). Hasil klasifikasi ini langsung menjadi dasar pengambilan keputusan aktuator pada mode otomatis, sehingga tindakan korektif ditentukan oleh model, bukan oleh ambang batas yang ditetapkan manual.

**e. Efisiensi Operasional dan Pengurangan Risiko Kerugian**
Dengan adanya otomatisasi penanganan awal serta peringatan dini berbasis IoT, risiko kematian biota air akibat *human error* atau keterlambatan penanganan dapat ditekan secara signifikan. Hal ini tidak hanya menjaga kelestarian ekosistem akuarium, tetapi juga mengoptimalkan efisiensi energi dan penggunaan perangkat pendukung secara lebih terukur.

**f. Skalabilitas dan Kemandirian Sistem (*Scalable & Self-Hosted*) [BARU]**
Dengan pola komunikasi *publish/subscribe* MQTT dan topik yang dipisahkan per node, akuarium baru dapat ditambahkan cukup dengan memasang perangkat yang sama, tanpa mengubah kode server. Seluruh komponen *cloud* (broker, *backend*, basis data, dan website) dibangun dan dijalankan secara mandiri, sehingga sistem tidak bergantung pada batasan maupun keberlangsungan layanan platform IoT pihak ketiga.

---

# BAB 2 RUMUSAN KONSEP SOLUSI IoT

## 2.1 Deskripsi Konsep Solusi [DIUBAH]

Sistem AI Smart Aquarium merupakan sebuah prototipe solusi berbasis Internet of Things (IoT) dan *Machine Learning* yang dirancang untuk mengotomatisasi pemantauan serta pengelolaan kualitas air pada ekosistem akuarium secara *real-time* dan responsif. Sistem ini mengintegrasikan tiga lapisan utama pengembangan IoT, yaitu akuisisi data fisik (*hardware*/sensor), transmisi dan pemrosesan data cerdas berbasis *cloud* (AI/*Machine Learning*), serta antarmuka pemantauan dan pengendalian bagi pengguna (**website**).

Pada lapisan fisik, sistem memanfaatkan mikrokontroler ESP8266 sebagai unit pemroses utama di tingkat perangkat (*edge node*). ESP8266 secara berkala membaca *telemetry* dari 2 jenis sensor, yaitu:
- Sensor Suhu Air (DS18B20) mengukur temperatur air secara presisi.
- Sensor Kekeruhan Air (*Turbidity Sensor*) mengukur tingkat kejernihan air.

Data *telemetry* yang diperoleh dari kedua sensor tersebut kemudian dikirimkan secara otomatis **melalui protokol MQTT ke server *cloud* yang dibangun secara mandiri. Server terdiri atas broker MQTT Mosquitto, *backend* berbahasa Go, basis data PostgreSQL, dan website berbasis Vue.js, yang seluruhnya dijalankan dengan Docker Compose pada sebuah *virtual machine* (VM).**

Di tingkat *cloud*, data diolah oleh tiga model *Machine Learning* berbasis algoritma Random Forest, yaitu model suhu, model kekeruhan, dan model kondisi keseluruhan (*overall*). Setiap model mengklasifikasikan kondisi ke dalam empat tingkat: NORMAL, CAUTION, WARNING, dan SEVERE. Model suhu dan model kekeruhan menjadi dasar keputusan aktuator pada Mode Kontrol Otomatis, sedangkan model *overall* menentukan status utama akuarium dan notifikasi darurat. Keputusan tersebut ditransmisikan kembali ke ESP8266 untuk menggerakkan aktuator melalui relay, yaitu Pompa Air dengan Internal Filter (sirkulasi dan penyaringan), Pemanas (menaikkan suhu), dan Kipas (menurunkan suhu). Sistem juga menyediakan Mode Kontrol Manual.

**Setiap ESP8266 memiliki identitas unik (*node_id*) yang diambil dari *chip ID*-nya. Node baru yang terhubung ke broker akan terdaftar secara otomatis di basis data, sehingga satu server dapat melayani banyak akuarium sekaligus.**

## 2.2 Batasan Sistem & Target Pengguna [DIUBAH]

Sistem ini dibatasi pada pemantauan dan pengendalian dua parameter kualitas air, yaitu suhu yang diukur dengan sensor suhu air dan kekeruhan yang diukur dengan sensor kekeruhan (*turbidity*). Parameter lain seperti ketinggian air, pH, kadar oksigen terlarut, amonia, dan salinitas tidak termasuk dalam cakupan sistem, sehingga penilaian kondisi aquarium hanya didasarkan pada kombinasi data suhu dan kekeruhan. Sistem bekerja secara *real time* sehingga membutuhkan koneksi jaringan yang stabil agar data dapat ditampilkan dan perintah kontrol dapat dijalankan tanpa jeda berarti. Pengendalian dilakukan melalui perangkat keluaran (aktuator), misalnya pemanas atau kipas pendingin untuk mengatur suhu serta pompa atau perangkat penyaringan untuk menangani kekeruhan. Sistem ditujukan untuk aquarium skala kecil hingga menengah, dan ketepatan hasil pengukurannya bergantung pada akurasi serta kalibrasi sensor yang digunakan. Model dilatih menggunakan dataset kualitas air kolam ikan air tawar, sehingga penerapannya pada akuarium mengasumsikan karakteristik suhu dan kekeruhan yang serupa.

**Seluruh layanan server dijalankan pada satu VM, sehingga ketersediaan sistem bergantung pada VM tersebut. Prototipe diuji dengan satu node fisik, sedangkan kemampuan multi-node diuji menggunakan simulator node yang mengirimkan data dengan format yang sama.**

Target pengguna sistem ini adalah pemilik aquarium, baik penghobi ikan hias maupun pelaku budidaya ikan skala kecil, yang membutuhkan pemantauan kondisi air secara berkelanjutan tanpa harus berada di dekat aquarium, **termasuk pengguna yang memiliki lebih dari satu akuarium**. Pengguna dapat memantau suhu dan kekeruhan air secara *real time* melalui antarmuka yang mudah dipahami, serta mengontrol kondisi aquarium dengan dua mode operasi. Pada mode otomatis, sistem menjalankan aktuator berdasarkan hasil klasifikasi model *Machine Learning*, misalnya menyalakan pemanas ketika model suhu mendeteksi kondisi Warning pada suhu yang lebih rendah dari kondisi normal, atau mengaktifkan pompa ketika model kekeruhan mendeteksi kondisi Warning atau Severe. Pada mode manual, pengguna dapat menghidupkan dan mematikan perangkat secara langsung sesuai kebutuhannya. Dengan demikian, pengguna yang sibuk dapat mengandalkan mode otomatis, sedangkan pengguna yang menginginkan kendali penuh dapat memilih mode manual.

## 2.3 Spesifikasi Fitur Utama

Spesifikasi fitur utama pada sistem AI Smart Aquarium dirancang untuk memberikan visibilitas penuh terhadap kondisi ekosistem air serta memberikan fleksibilitas pengendalian aktuator. Fitur-fitur ini dikelompokkan ke dalam dua ranah fungsionalitas utama, yaitu Fitur Monitoring (pemantauan) dan Fitur Kontrol (pengendalian), yang dapat beroperasi dalam dua mode kerja, yakni Mode Otomatis dan Mode Manual.

### 2.3.1 Fitur Monitoring [DIUBAH]

Fitur pemantauan berfungsi untuk mengumpulkan, mengolah, dan menyajikan data parameter kualitas air akuarium secara *real-time* kepada pengguna melalui **website**. Parameter dan informasi yang dipantau mencakup:

- **Pemantauan Suhu Air (*Water Temperature Monitoring*)**
  Menggunakan sensor suhu DS18B20 (*waterproof probe*) yang dipasang secara terendam di dalam tangki akuarium. Sensor ini secara kontinu mengukur temperatur air dengan satuan derajat Celsius (°C) untuk mendeteksi fluktuasi suhu yang berpotensi memengaruhi metabolisme biota air.

- **Pemantauan Kekeruhan Air (*Turbidity Monitoring*)**
  Menggunakan *Turbidity Sensor* yang bekerja mengukur tingkat kekeruhan atau transparansi air akuarium. Nilai pembacaan sensor dikalibrasi ke dalam satuan *Nephelometric Turbidity Unit* (NTU) untuk mengindikasikan tingkat penumpukan partikel terdispersi, pakan, atau kotoran di dalam air.

- **Status Klasifikasi Kondisi Akuarium (*AI Condition Classification*)**
  Menampilkan hasil analisis model *Machine Learning* di server *cloud* yang mengklasifikasikan kondisi keseluruhan akuarium berdasarkan interaksi data suhu dan kekeruhan. Status dikategorikan ke dalam 4 tingkatan kelas:
  1. NORMAL: parameter berada dalam rentang ideal; tidak diperlukan tindakan.
  2. CAUTION: penyimpangan ringan; sistem memberikan notifikasi peringatan dini tanpa menyalakan aktuator.
  3. WARNING: penyimpangan yang memerlukan penanganan; sistem mengaktifkan aktuator yang sesuai.
  4. SEVERE: kondisi kritis; sistem mengaktifkan aktuator yang sesuai dan mengirimkan peringatan darurat.

- **Status Parameter Individual & Indikator Aktuator**
  Dashboard menyajikan hasil klasifikasi model suhu dan model kekeruhan (masing-masing Normal/Caution/Warning/Severe), status aktuator (Pemanas, Kipas, Pompa Filter: ON/OFF), serta mode kontrol yang aktif (Automatic/Manual).

- **Pemantauan Multi-Node [BARU]**
  Website menampilkan daftar seluruh node akuarium beserta status koneksinya (*online*/*offline*), nilai sensor terakhir, dan status kondisinya. Pengguna dapat membuka halaman detail untuk setiap node.

- **Riwayat Data dan Notifikasi [BARU]**
  Seluruh data sensor, hasil klasifikasi, dan riwayat aksi aktuator tersimpan di basis data PostgreSQL, sehingga dapat ditampilkan dalam grafik tren dengan rentang waktu yang dapat dipilih. Peringatan dikirim melalui bot Telegram dan ditampilkan pada website.

### 2.3.2 Fitur Kontrol dan Mode Operasional [DIUBAH sedikit]

Fitur kontrol memungkinkan tindakan korektif dieksekusi baik secara otomatis oleh sistem maupun secara langsung oleh pengguna melalui interaksi jarak jauh (*remote control*). **Mode operasional diatur untuk setiap node secara terpisah.**

- **Mode Kontrol Otomatis (*Automatic Action Mode*)**
  Pada mode otomatis, keputusan aktuator ditentukan oleh hasil klasifikasi model suhu dan model kekeruhan.
  1. Respons suhu (pemanas dan kipas): jika model suhu memprediksi Warning atau Severe, sistem menyalakan pemanas apabila suhu lebih rendah dari suhu acuan normal, atau menyalakan kipas apabila suhu lebih tinggi. Suhu acuan adalah median suhu kelas Normal pada dataset (27,48 °C). Jika prediksi Normal, pemanas dan kipas dimatikan. Jika prediksi Caution, sistem hanya mengirim notifikasi.
  2. Respons kekeruhan (pompa filter): jika model kekeruhan memprediksi Warning atau Severe, pompa filter dinyalakan. Jika Normal, pompa dimatikan. Jika Caution, sistem hanya mengirim notifikasi.
  3. Anomali ganda: kedua respons berjalan bersamaan, sehingga pemanas/kipas dan pompa dapat aktif sekaligus. Model *overall* menandai kondisi ini sebagai status utama, dan status Severe memicu peringatan darurat.
  4. Kondisi normal: seluruh aktuator berada dalam keadaan OFF.

- **Mode Kontrol Manual (*Manual Action Mode*)**
  Mode manual memberikan kendali penuh kepada pengguna untuk mengambil alih fungsi kontrol aktuator kapan saja (*override*) melalui **website**. Fitur kontrol manual yang tersedia meliputi:
  1. Peralihan Mode Operasional (*Mode Switch*): Tombol navigasi untuk mengalihkan mode sistem antara *Automatic Mode* dan *Manual Mode* secara *real-time*.
  2. Kontrol manual aktuator: tiga sakelar digital untuk menyalakan atau mematikan pemanas, kipas, dan pompa filter secara langsung dari antarmuka.

---

# BAB 3 STUDI LITERATUR & REFERENSI

## 3.1 Tinjauan Pustaka

### 3.1.1 Landasan Teori

**a. Kualitas Air Akuarium**
Kualitas air menentukan kesehatan dan kelangsungan hidup biota akuarium. Dua parameter fisik yang dipantau dalam sistem ini adalah suhu dan kekeruhan. Suhu memengaruhi metabolisme dan tingkat stres ikan, sehingga perlu dijaga dalam rentang yang sesuai dengan jenis ikan yang dipelihara. Ikan hias tropis seperti manfish membutuhkan suhu air antara 26 °C hingga 30 °C [5]. Kekeruhan disebabkan oleh partikel tersuspensi seperti sisa pakan, kotoran ikan, dan bahan organik yang menumpuk di akuarium. Kekeruhan dinyatakan dalam satuan NTU (*Nephelometric Turbidity Unit*). Penelitian Andini dkk. menggunakan batas 30 NTU sebagai ambang air keruh pada akuarium ikan mas [3].

**b. Internet of Things (IoT)**
IoT adalah konsep yang menghubungkan perangkat fisik seperti sensor, mikrokontroler, dan aktuator ke internet, sehingga data dapat dikumpulkan, dipertukarkan, dan perangkat dapat dikendalikan dari jarak jauh tanpa interaksi langsung manusia. Dalam sistem pemantauan akuarium, IoT memungkinkan kondisi air dipantau secara berkelanjutan dan aktuator dikendalikan melalui aplikasi.

**c. Protokol MQTT (*Message Queuing Telemetry Transport*) [DIUBAH, menggantikan Blynk IoT Platform]**
MQTT adalah protokol komunikasi ringan berbasis pola *publish/subscribe* yang dirancang untuk perangkat dengan sumber daya terbatas dan jaringan yang tidak selalu stabil [14]. Pada MQTT, perangkat tidak berkomunikasi secara langsung satu sama lain, melainkan melalui perantara yang disebut *broker*. Pengirim (*publisher*) mengirim pesan ke sebuah *topic*, lalu broker meneruskannya ke semua penerima (*subscriber*) yang berlangganan *topic* tersebut. Struktur *topic* yang hierarkis (misalnya `aquarium/node01/telemetry`) memungkinkan data dari banyak perangkat dipisahkan dan dilanggani sekaligus menggunakan *wildcard* (`aquarium/+/telemetry`).

MQTT menyediakan tiga tingkat jaminan pengiriman (*Quality of Service*, QoS): QoS 0 (paling banyak sekali), QoS 1 (minimal sekali), dan QoS 2 (tepat sekali). Fitur *retained message* membuat broker menyimpan pesan terakhir pada suatu *topic* sehingga *subscriber* baru langsung menerimanya, sedangkan *Last Will and Testament* (LWT) membuat broker otomatis mengirim pesan tertentu ketika koneksi sebuah perangkat terputus secara tidak normal. Kedua fitur ini dimanfaatkan untuk menyimpan mode operasi node dan mendeteksi status *online*/*offline* node.

**d. Random Forest**
Random Forest adalah algoritma *ensemble learning* yang membangun banyak pohon keputusan (*decision tree*). Setiap pohon dilatih menggunakan sampel data acak hasil *bootstrap* dan subset fitur acak pada setiap percabangan. Hasil klasifikasi akhir ditentukan dengan *majority voting* dari seluruh pohon [4]. Pendekatan ini membuat Random Forest lebih tahan terhadap *overfitting* dan derau dibandingkan satu pohon keputusan, sehingga cocok untuk data sensor yang mengandung fluktuasi.

Pada sistem ini digunakan tiga model Random Forest: model suhu (masukan: suhu), model kekeruhan (masukan: kekeruhan), dan model *overall* (masukan: suhu dan kekeruhan). Ketiganya dilatih menggunakan label tingkat degradasi kualitas air yang sama pada dataset, sehingga batas setiap tingkat kondisi dipelajari dari data, bukan ditetapkan secara manual.

**e. Relay**
Relay adalah sakelar elektromagnetik. Ketika kumparan relay dialiri arus dari sinyal kendali, kontak mekanis berpindah dan menyambungkan atau memutus rangkaian beban. Dengan relay, mikrokontroler bertegangan rendah dapat mengendalikan perangkat bertegangan tinggi seperti pemanas dan pompa 220 VAC.

### 3.1.2 Penelitian Terkait [DIUBAH, diperluas menjadi 10 penelitian]

**Tabel 3.1** Perbandingan penelitian terkait

| No | Peneliti | Metode / Teknologi | Parameter | Hasil | Keterbatasan |
|---|---|---|---|---|---|
| 1 | Nasir & Natasya (2020) [1] | Arduino Uno, Django, SQLite, laptop sebagai server melalui kabel USB | Kekeruhan, ketinggian air, pH, suhu (LM35) | Data tampil di website dalam bentuk grafik dan angka | Hanya pemantauan, tanpa kontrol aktuator; perangkat bergantung pada kabel USB ke laptop |
| 2 | Qolbu, Taqwa & Fadhli (2026) [2] | LightGBM, 710 data sensor, *train-test split* 80:20 | pH, suhu, amonia, TDS | Akurasi 99,30%, klasifikasi Baik/Waspada/Buruk | Model diuji secara *offline*, belum terintegrasi dengan sistem IoT dan aktuator; tidak memakai kekeruhan |
| 3 | Andini, Sunardi & Rachmansyah (2024) [3] | ESP32, Blynk, relay, driver motor, LCD | Kekeruhan, ketinggian air, pH | Pompa pengurasan aktif saat kekeruhan > 30 NTU; data terpantau di Blynk | Keputusan berbasis ambang batas tunggal per parameter; bergantung pada platform Blynk |
| 4 | Islam dkk. (2023) [7] | Arduino Uno, ThingSpeak (REST API); 10 algoritma klasifikasi di WEKA | pH, suhu, kekeruhan, konduktivitas | Random Forest terbaik dengan akurasi 94,42% untuk memprediksi kesesuaian kolam bagi spesies ikan | Klasifikasi dilakukan *offline*; tanpa kontrol aktuator; bergantung pada ThingSpeak |
| 5 | Adeleke, Nwulu & Ogbolumani (2023) [8] | Arduino Nano + ESP8266, ANN dan SVM, penyimpanan di ThingSpeak | Suhu, pH, kekeruhan, DO, TDS, ORP, EC | ANN memberikan akurasi tertinggi dan digunakan untuk prediksi status air secara *real-time* | Mekanisme pengolahan air aktif berdasarkan ambang ORP/pH, bukan hasil model; bergantung pada platform pihak ketiga |
| 6 | Olanubi, Akano & Asaolu (2024) [9] | ESP32, DS18B20, SEN0189, Firebase, dashboard web, notifikasi WhatsApp | Suhu, pH, kekeruhan | Pompa masuk dan pompa buang aktif otomatis; uji 21 hari pada kolam lele | Kontrol berbasis ambang batas tanpa *machine learning*; bergantung pada Firebase |
| 7 | Kanwal dkk. (2024) [10] | ESP32, DS18B20, sensor kekeruhan, ThingSpeak + Firebase; 11 algoritma ML | pH, suhu, kekeruhan | Decision Tree terbaik dengan akurasi sekitar 99,8% untuk memprediksi kelangsungan hidup ikan | Tanpa aktuator; model untuk analisis kesesuaian kolam, bukan pengendalian |
| 8 | Arepalli & Naik (2024) [11] | Arduino + NodeMCU, model *deep learning* AODEGRU | pH, DO, suhu, kekeruhan, amonia, nitrat, mangan | Akurasi validasi 99,85% untuk klasifikasi air tercemar/tidak tercemar | Fokus pada model prediksi; tanpa kontrol aktuator |
| 9 | Shete, Bongale & Dharrao (2024) [12] | Arduino Uno, modul GSM SIM800L, server *cloud*, aplikasi *mobile*, pelampung IP66 | DO, pH, suhu | Pengiriman data setiap 15 menit dengan cadangan baterai ≥ 48 jam; galat maksimum 4,87% | Tanpa klasifikasi *machine learning* dan tanpa aktuator fisik |
| 10 | Hridoy dkk. (2025) [13] | Pemantauan IoT kolam tilapia selama 6 bulan, *neural network* TensorFlow | Suhu, DO, pH, kekeruhan | Prediksi parameter dengan R² 0,83–0,85 | Tanpa aktuator; hasil hanya untuk analisis dan rekomendasi manual |

### 3.1.3 Posisi Penelitian [DIUBAH]

Dari sepuluh penelitian pada Tabel 3.1, terdapat tiga kesenjangan utama:

1. **Hasil *machine learning* belum menjadi dasar tindakan.** Penelitian [2], [7], [10], [11], dan [13] menghasilkan model dengan akurasi tinggi, tetapi model hanya digunakan untuk analisis, prediksi, atau penilaian kesesuaian kolam. Pada penelitian yang memiliki aktuator ([3], [8], [9]), aktuator dikendalikan oleh ambang batas per parameter, bukan oleh hasil model.
2. **Ketergantungan pada platform pihak ketiga.** Sebagian besar sistem bergantung pada layanan seperti Blynk [3], ThingSpeak [7], [8], [10], dan Firebase [9], [10], sehingga kapasitas, fitur, dan penyimpanan data dibatasi oleh kebijakan layanan tersebut.
3. **Rancangan untuk satu lokasi.** Sistem pada umumnya dirancang untuk satu kolam atau satu akuarium, tanpa mekanisme yang jelas untuk menambah perangkat baru.

Sistem Smart Aquarium ini menjawab kesenjangan tersebut melalui:
1. Tiga model Random Forest (suhu, kekeruhan, *overall*) yang hasil klasifikasinya langsung menentukan aksi aktuator, sehingga keputusan tindakan berasal dari model yang dilatih pada data, bukan dari ambang batas tunggal per parameter.
2. Server *cloud* yang dibangun secara mandiri (broker MQTT, *backend*, basis data, dan website), tanpa bergantung pada platform IoT pihak ketiga.
3. Rancangan multi-node berbasis topik MQTT per node dan pendaftaran node otomatis, sehingga akuarium baru dapat ditambahkan tanpa mengubah kode server.
4. Pengendalian tiga aktuator (pemanas, kipas, pompa filter) secara otomatis maupun manual, serta mekanisme *fail-safe* pada perangkat sehingga akuarium tetap terkendali ketika server tidak dapat dihubungi.

### 3.1.4 Dataset [DIUBAH, tabel hasil disesuaikan dengan model yang di-deploy]

Model dilatih menggunakan dataset kualitas air kolam tanah ikan air tawar yang terdiri atas 2.153 sampel dengan tujuh parameter (suhu, pH, oksigen terlarut, kekeruhan, konduktivitas listrik, TDS, dan ORP) [6]. Setiap sampel memiliki label tingkat degradasi kualitas air: Normal (638), Caution (670), Warning (478), dan Severe (367). Sesuai sensor yang tersedia, hanya fitur suhu (°C) dan kekeruhan (NTU) yang digunakan.

Data dibagi 80:20 secara *stratified* (1.722 data latih dan 431 data uji), lalu setiap model diuji juga dengan *5-fold cross validation*. **Agar model dapat dikonversi menjadi kode program Go (lihat 3.2.4), setiap model dibatasi menjadi 50 pohon dengan kedalaman maksimum 8.** Hasil awal pelatihan ditunjukkan pada Tabel 3.2.

**Tabel 3.2** Hasil awal pelatihan model

| Model | Fitur | Akurasi uji | F1 macro | CV 5-fold |
|---|---|---|---|---|
| Suhu | Suhu | 99,07% | 99,02% | 99,63 ± 0,24% |
| Kekeruhan | Kekeruhan | 99,77% | 99,78% | 99,91 ± 0,11% |
| Overall | Suhu + kekeruhan | 100,00% | 100,00% | 100,00 ± 0,00% |

Hasil ini merupakan gambaran awal kinerja model dan akan diperbarui pada tahap pengujian. Akurasi yang sangat tinggi menunjukkan batas antartingkat pada dataset sangat tegas. Hal ini menjadi keterbatasan, karena kondisi akuarium nyata kemungkinan memiliki batas yang lebih bervariasi.

Kelas Normal pada dataset berada pada rentang suhu 25–30 °C, sedangkan kondisi abnormal dapat berupa suhu di bawah maupun di atas rentang tersebut. Oleh karena itu, median suhu kelas Normal (27,48 °C) digunakan sebagai acuan arah untuk menentukan apakah pemanas atau kipas yang diaktifkan.

## 3.2 Spesifikasi Komponen Perangkat Keras & Lunak

### 3.2.1 Mikrokontroler / Processing Unit: NodeMCU ESP8266 (ESP-12E) [DIUBAH sedikit]

NodeMCU ESP8266 berperan sebagai *edge node*. Papan ini membaca sensor, mengirim data **ke broker MQTT**, menerima perintah aktuator, dan menjalankan logika *fail-safe* ketika server tidak aktif. **Setiap papan memakai *chip ID*-nya sebagai identitas node (*node_id*).**

**Tabel 3.3** Spesifikasi NodeMCU ESP8266

| Parameter | Spesifikasi |
|---|---|
| Prosesor | Tensilica L106 32-bit, 80/160 MHz |
| Memori flash | 4 MB |
| Konektivitas | WiFi 802.11 b/g/n, 2,4 GHz |
| Tegangan logika | 3,3 V |
| Pin analog | 1 buah (A0), ADC 10-bit (0–1023) |
| Sumber daya | 5 V melalui micro USB |

### 3.2.2 Sensor (Fungsi & Cara Kerja)

**a. Sensor Suhu DS18B20 (Waterproof Probe)**

*Fungsi:* mengukur suhu air akuarium dalam derajat Celsius (°C). Nilai suhu menjadi fitur masukan model klasifikasi. Masukan model suhu yang menentukan pengendalian pemanas serta kipas.

**Tabel 3.4** Spesifikasi DS18B20

| Parameter | Spesifikasi |
|---|---|
| Rentang ukur | −55 °C s.d. +125 °C |
| Akurasi | ±0,5 °C (−10 °C s.d. +85 °C) |
| Resolusi | 9–12 bit (0,0625 °C pada 12 bit) |
| Tegangan kerja | 3,0–5,5 V |
| Antarmuka | 1-Wire |

*Cara kerja:*
1. Elemen sensor di dalam chip menghasilkan nilai yang bergantung pada temperatur.
2. ADC internal mengubah nilai tersebut menjadi data digital 9–12 bit, lalu menyimpannya di *scratchpad memory*.
3. Mikrokontroler mengirim perintah konversi melalui protokol 1-Wire, menunggu konversi selesai (±750 ms pada 12 bit), lalu membaca hasilnya.
4. Sensor dipasang melalui modul adapter terminal yang sudah dilengkapi resistor *pull-up* 4,7 kΩ pada jalur data.

**b. Sensor Kekeruhan (Turbidity Sensor)**

*Fungsi:* mengukur tingkat kekeruhan air akibat partikel tersuspensi. Nilai kekeruhan (NTU) menjadi fitur masukan model klasifikasi dan dasar pengendalian pompa filter.

*Spesifikasi:* modul analog (sejenis DFRobot SEN0189) dengan keluaran analog dan digital, serta potensiometer untuk ambang keluaran digital. Pada sistem ini, sensor diberi daya 3,3 V supaya tegangan keluaran tidak melebihi batas pin A0.

*Cara kerja:*
1. Probe berisi LED inframerah dan fototransistor yang saling berhadapan.
2. Pada air jernih, cahaya yang diterima fototransistor besar sehingga tegangan keluaran tinggi.
3. Pada air keruh, partikel menyebarkan dan menyerap cahaya sehingga tegangan keluaran turun.
4. Tegangan dibaca melalui pin A0 (rata-rata 20 sampel untuk mengurangi derau), lalu dikonversi menjadi NTU.

*Kalibrasi:* rumus konversi bawaan pabrikan dirancang untuk tegangan 5 V, sehingga tidak berlaku langsung pada tegangan 3,3 V. Konversi dilakukan dengan kalibrasi dua titik, yaitu mengukur tegangan pada air jernih dan pada air dengan kekeruhan acuan, lalu membentuk persamaan linear NTU = a·V + b.

### 3.2.3 Aktuator (Fungsi & Mekanisme Penggerak)

Ketiga aktuator dikendalikan oleh modul relay 4 channel 3,3 V. Saat pin IN diberi sinyal dari NodeMCU, kumparan relay menarik kontak sehingga terminal COM tersambung ke NO dan aktuator menyala. Sisi kendali (DC 3,3 V) dan sisi beban (AC 220 V) dipisahkan oleh optocoupler, sehingga tegangan tinggi tidak masuk ke mikrokontroler. Kontak relay mampu menahan beban hingga 10 A 250 VAC.

Setiap aktuator dihubungkan ke stop kontak 1 lubang tersendiri. Kabel fasa dari setiap stop kontak dipotong dan dilewatkan pada terminal COM–NO relay, sedangkan kabel netral dan ground tetap utuh.

**Tabel 3.5** Spesifikasi dan fungsi aktuator

| Aktuator | Fungsi | Mekanisme penggerak | Channel / Pin |
|---|---|---|---|
| Pemanas akuarium 100 W, 220 VAC | Menaikkan suhu saat model suhu mendeteksi kondisi Warning/Severe di bawah suhu acuan | Elemen resistif di dalam tabung kaca mengubah energi listrik menjadi panas yang dirambatkan ke air | IN1 / D6 |
| Kipas DC 12 V 8 cm + adaptor | Menurunkan suhu saat model suhu mendeteksi kondisi Warning/Severe di atas suhu acuan | Motor DC memutar baling-baling dan meniupkan udara ke permukaan air sehingga terjadi pendinginan melalui penguapan | IN2 / D7 |
| Pompa internal filter 4 W, 220 VAC | Mensirkulasikan dan menyaring air saat model kekeruhan mendeteksi Warning/Severe | Motor pompa menghisap air melewati media filter lalu mengembalikannya ke akuarium | IN3 / D1 |

### 3.2.4 Protocol & Cloud Platform [DIUBAH TOTAL]

Seluruh layanan *cloud* dibangun secara mandiri tanpa platform IoT pihak ketiga. Komponennya adalah sebagai berikut.

**a. Protokol MQTT dan Broker Mosquitto**
Komunikasi antara node dan server menggunakan protokol MQTT. Broker yang digunakan adalah Eclipse Mosquitto, yaitu broker MQTT *open source* yang ringan. Setiap node terhubung menggunakan *username* dan *password*, dan aturan ACL (*access control list*) membatasi setiap node agar hanya dapat menulis dan membaca topik miliknya sendiri (`aquarium/{node_id}/...`). Rancangan topik dijelaskan pada subbab 4.3.3. Pada prototipe ini, koneksi MQTT dari node belum dienkripsi; penggunaan MQTT dengan TLS (port 8883) menjadi pengembangan selanjutnya.

**b. Backend (Go)**
*Backend* ditulis dalam bahasa Go. Tugasnya:
1. Berlangganan (*subscribe*) seluruh topik node menggunakan *wildcard*, lalu menyimpan data ke basis data.
2. Mendaftarkan node baru secara otomatis ketika node tersebut pertama kali mengirim data.
3. Menjalankan inferensi tiga model Random Forest dan menentukan aksi aktuator pada mode otomatis.
4. Mengirim perintah aktuator ke node melalui MQTT.
5. Menyediakan REST API dan WebSocket bagi website.
6. Mengirim *heartbeat* berkala ke seluruh node dan mengirim notifikasi melalui bot Telegram.

Bahasa Go dipilih karena ringan dalam penggunaan memori, menghasilkan satu file program yang mudah dijalankan, dan mendukung pemrosesan paralel (*goroutine*) yang sesuai untuk melayani banyak node sekaligus.

**c. Inferensi Model (m2cgen)**
Model Random Forest dilatih menggunakan Python (scikit-learn) di Google Colab. Model yang sudah dilatih kemudian dikonversi menjadi kode program Go menggunakan pustaka m2cgen, sehingga *backend* dapat menjalankan inferensi secara langsung tanpa layanan Python terpisah. Konversi ini membutuhkan model yang tidak terlalu besar, sehingga setiap model dibatasi menjadi 50 pohon dengan kedalaman maksimum 8. Nilai suhu acuan (27,48 °C) disimpan sebagai konstanta pada *backend*.

**d. Basis Data (PostgreSQL)**
PostgreSQL menyimpan data pengguna, data node, riwayat sensor, hasil klasifikasi, riwayat aksi aktuator, dan riwayat notifikasi. Berbeda dengan platform pihak ketiga, lama penyimpanan data tidak dibatasi oleh kebijakan layanan. Rancangan tabel dijelaskan pada subbab 4.3.5.

**e. Website (Vue.js)**
Antarmuka pengguna dibangun sebagai *single page application* menggunakan Vue 3 dan Vite, dengan Tailwind CSS untuk tampilan dan ECharts untuk grafik. Website mengambil data melalui REST API dan menerima pembaruan *real-time* melalui WebSocket. Akses website dilindungi dengan *login* berbasis JWT (*JSON Web Token*). Tampilan dibuat responsif sehingga dapat digunakan dari komputer maupun ponsel.

**f. Notifikasi (Bot Telegram)**
Peringatan dikirim melalui Telegram Bot API ketika status *overall* suatu node berubah menjadi selain NORMAL. Status SEVERE dikirim sebagai peringatan darurat.

**g. Hosting dan Deployment**
Broker Mosquitto, *backend* Go, PostgreSQL, dan Caddy dijalankan sebagai kontainer menggunakan Docker Compose pada sebuah *virtual machine* (VM) Ubuntu di Microsoft Azure (program Azure for Students). Caddy berperan sebagai *reverse proxy* yang menyajikan berkas website, meneruskan permintaan API dan WebSocket ke *backend*, serta mengurus sertifikat HTTPS secara otomatis. Nama domain diperoleh dari layanan DNS dinamis gratis DuckDNS (misalnya `[isi-nama].duckdns.org`).

**Tabel 3.6** Komponen server dalam Docker Compose

| Kontainer | Perangkat lunak | Port | Fungsi |
|---|---|---|---|
| `mosquitto` | Eclipse Mosquitto | 1883 (MQTT) | Broker pesan antara node dan *backend* |
| `backend` | Aplikasi Go | 8080 (internal) | Logika sistem, inferensi, REST API, WebSocket |
| `db` | PostgreSQL | 5432 (internal) | Penyimpanan data |
| `caddy` | Caddy | 80, 443 | *Reverse proxy*, HTTPS, penyaji website |

Port bertanda "internal" hanya dapat diakses di dalam jaringan Docker dan tidak dibuka ke internet.

**Tabel 3.7** Perangkat lunak yang digunakan

| Lapisan | Perangkat lunak | Kegunaan |
|---|---|---|
| *Firmware* | Arduino IDE + board ESP8266 | Pemrograman perangkat |
| | OneWire, DallasTemperature | Pembacaan DS18B20 |
| | PubSubClient | Koneksi MQTT |
| | ArduinoJson | Format pesan JSON |
| | WiFiManager | Pengaturan WiFi node melalui portal tanpa *flash* ulang |
| Pelatihan model | Python 3, scikit-learn, Google Colab | Melatih tiga model Random Forest |
| | m2cgen | Mengonversi model menjadi kode Go |
| *Backend* | Go, pustaka klien MQTT (Eclipse Paho) dan PostgreSQL untuk Go | Logika server, API, inferensi |
| Basis data | PostgreSQL | Penyimpanan data |
| *Frontend* | Vue 3, Vite, Tailwind CSS, ECharts | Website pemantauan dan kontrol |
| Infrastruktur | Docker, Docker Compose, Caddy, DuckDNS | Menjalankan layanan, HTTPS, nama domain |
| Notifikasi | Telegram Bot API | Pengiriman peringatan |

---

# BAB 4 PERANCANGAN DETAIL SISTEM

## 4.1 Diagram Proses Bisnis / Alur Kerja Sistem (Flowchart)

Alur kerja sistem dibagi menjadi dua flowchart, yaitu alur pada perangkat (NodeMCU) dan alur pada server.

### 4.1.1 Alur Kerja Perangkat [DIUBAH]

Perangkat diawali dengan inisialisasi sensor dan relay, dengan seluruh aktuator dalam keadaan OFF, lalu membaca *chip ID* sebagai *node_id*. Jika belum ada jaringan WiFi yang tersimpan, perangkat membuka portal konfigurasi WiFiManager sehingga pengguna dapat memasukkan nama dan kata sandi WiFi melalui ponsel. Setelah terhubung ke WiFi, perangkat terhubung ke broker MQTT dengan pesan LWT `offline`, mengirim status `online`, lalu berlangganan topik perintah (`cmd`), konfigurasi mode (`config`), dan *heartbeat* server.

Setiap 60 detik, perangkat membaca suhu dan kekeruhan, lalu mengirimkannya ke topik `aquarium/{node_id}/telemetry`. Interval ini dipilih karena suhu dan kekeruhan akuarium berubah lambat, sekaligus menjaga ukuran basis data. Ketika perintah aktuator diterima pada topik `cmd`, perangkat langsung menyalakan atau mematikan relay sesuai perintah, lalu melaporkan status relay aktual ke topik `state`.

Perangkat juga memantau *heartbeat* server. Jika mode node adalah otomatis dan *heartbeat* tidak diterima lebih dari 3 menit, atau koneksi ke broker terputus, perangkat beralih ke mode lokal (*fail-safe*) dan mengendalikan aktuator berdasarkan rentang kelas Normal pada dataset (suhu 25–30 °C, kekeruhan < 25 NTU). Selama koneksi terputus, perangkat terus mencoba terhubung kembali. Mode lokal hanya berfungsi sebagai cadangan saat server tidak dapat dihubungi; pada kondisi normal, seluruh keputusan aktuator berasal dari model.

> **Gambar 4.1 Flowchart alur kerja perangkat. [PERLU DIGAMBAR ULANG]**
> Urutan simbol:
> 1. (Oval) Mulai
> 2. (Proses) Inisialisasi sensor, relay OFF, baca *chip ID* sebagai node_id
> 3. (Keputusan) WiFi tersimpan & terhubung? Tidak → (Proses) Buka portal WiFiManager → kembali ke keputusan ini. Ya → lanjut
> 4. (Proses) Hubungkan ke broker MQTT (LWT `offline`), kirim status `online`, *subscribe* `cmd`, `config`, `server/heartbeat`
> 5. (Input) Baca suhu & kekeruhan
> 6. (Proses) Konversi ke °C & NTU (rata-rata 20 sampel)
> 7. (Output) *Publish* ke `aquarium/{node_id}/telemetry`
> 8. (Keputusan) Mode otomatis **dan** (*heartbeat* > 3 menit **atau** MQTT terputus)?
>    - Ya → (Proses) Mode lokal: relay sesuai rentang Normal (25–30 °C, < 25 NTU)
>    - Tidak → (Proses) Relay mengikuti perintah terakhir dari topik `cmd`
> 9. (Output) *Publish* status relay ke `aquarium/{node_id}/state`
> 10. (Proses) Tunggu 60 detik (sambil memproses pesan masuk dan menyambung ulang MQTT jika terputus) → kembali ke langkah 5
>
> Catatan: keputusan notifikasi di gambar lama ("V2 berubah & ≠ NORMAL?") **dihapus**, karena notifikasi sekarang dikirim oleh server.

### 4.1.2 Alur Kerja Server [DIUBAH]

Server bekerja berbasis *event*: setiap pesan *telemetry* yang masuk langsung diproses, tanpa *polling* berkala. Ketika *backend* menerima pesan pada topik `aquarium/{node_id}/telemetry`, *backend* memeriksa apakah node tersebut sudah terdaftar. Jika belum, node didaftarkan secara otomatis ke tabel `nodes` dengan nama bawaan. Data sensor kemudian disimpan, lalu dijadikan masukan tiga model Random Forest: model suhu, model kekeruhan, dan model *overall*. Hasil klasifikasi disimpan ke tabel `predictions`.

Jika mode node adalah otomatis, *backend* menentukan status aktuator langsung dari hasil klasifikasi sesuai Tabel 4.1. Perintah hanya dikirim ke topik `aquarium/{node_id}/cmd` ketika status aktuator berubah, lalu dicatat di tabel `actuator_logs`. Jika status *overall* berubah menjadi selain NORMAL, *backend* mengirim notifikasi melalui bot Telegram; status SEVERE dikirim sebagai peringatan darurat. Setiap perubahan data diteruskan ke website melalui WebSocket.

Secara paralel, *backend* mengirim *heartbeat* ke topik `server/heartbeat` setiap 30 detik dan memperbarui status *online*/*offline* node berdasarkan pesan pada topik `status`. Pada mode manual, *backend* tetap melakukan klasifikasi tetapi tidak mengirim perintah aktuator; perintah berasal dari pengguna melalui website.

**Tabel 4.1** Aturan aktuator berdasarkan hasil klasifikasi

| Kondisi | Pemanas | Kipas | Pompa filter |
|---|---|---|---|
| Suhu = Normal | OFF | OFF | – |
| Suhu = Caution | tetap (notifikasi) | tetap (notifikasi) | – |
| Suhu = Warning/Severe, suhu < 27,48 °C | ON | OFF | – |
| Suhu = Warning/Severe, suhu ≥ 27,48 °C | OFF | ON | – |
| Kekeruhan = Normal | – | – | OFF |
| Kekeruhan = Caution | – | – | tetap (notifikasi) |
| Kekeruhan = Warning/Severe | – | – | ON |

"Tetap" berarti aktuator mempertahankan status sebelumnya. Dengan begitu, pemanas yang sedang menyala saat suhu membaik dari Warning ke Caution tidak langsung mati, dan baru mati saat model memprediksi Normal. Ini mencegah aktuator menyala-mati berulang di batas kelas.

> **Gambar 4.2 Flowchart alur kerja server. [PERLU DIGAMBAR ULANG]**
> Urutan simbol:
> 1. (Oval) Mulai
> 2. (Proses) Muat model (kode Go hasil m2cgen), hubungkan ke PostgreSQL dan broker MQTT
> 3. (Proses) *Subscribe* `aquarium/+/telemetry`, `aquarium/+/state`, `aquarium/+/status`
> 4. (Input) Terima pesan *telemetry* dari node
> 5. (Keputusan) Node sudah terdaftar? Tidak → (Proses) Daftarkan node baru → lanjut. Ya → lanjut
> 6. (Proses) Simpan data sensor
> 7. (Proses) Prediksi kelas suhu, kekeruhan, dan *overall* → simpan hasil
> 8. (Keputusan) Mode node = Auto? Ya → (Proses) Tentukan aksi sesuai Tabel 4.1 → (Keputusan) Aksi berubah? Ya → (Output) *Publish* ke `aquarium/{node_id}/cmd` dan simpan log
> 9. (Keputusan) Status *overall* berubah & ≠ NORMAL? Ya → (Output) Kirim notifikasi Telegram
> 10. (Output) Kirim pembaruan ke website (WebSocket)
> 11. Kembali ke langkah 4 (menunggu pesan berikutnya)
>
> Di samping alur utama, gambarkan satu proses paralel kecil: "Setiap 30 detik → *publish* `server/heartbeat`".

## 4.2 Arsitektur Sistem IoT [DIUBAH]

Arsitektur sistem terdiri atas lima lapisan:

1. **IoT Devices:** sensor DS18B20 dan sensor kekeruhan sebagai masukan; modul relay 4 channel beserta pemanas, kipas, dan pompa filter sebagai keluaran. Setiap akuarium memiliki satu set perangkat ini.
2. **Gateway:** NodeMCU ESP8266 pada setiap akuarium berperan sebagai *edge gateway*. Perangkat ini mengolah data mentah (perata-rataan, konversi ke °C dan NTU), mengemasnya dalam format JSON, dan menjalankan logika *fail-safe* sebelum terhubung ke jaringan, sehingga tidak diperlukan perangkat gateway terpisah.
3. **Network:** router WiFi 2,4 GHz menghubungkan setiap node ke internet melalui TCP/IP. Data dikirim menggunakan protokol MQTT.
4. **Cloud:** sebuah VM yang menjalankan empat kontainer Docker, yaitu broker Mosquitto, *backend* Go (inferensi tiga model Random Forest, logika mode otomatis, API), basis data PostgreSQL, dan Caddy (*reverse proxy* dan HTTPS).
5. **Web UI:** website Vue.js yang diakses melalui *browser* komputer maupun ponsel untuk pemantauan dan kontrol, serta bot Telegram untuk notifikasi.

> **Gambar 4.3 Arsitektur sistem. [PERLU DIGAMBAR ULANG]**
> Lima kolom dari kiri ke kanan sesuai lapisan di atas.
> - Kolom IoT Devices dan Gateway: gambar **2–3 node bertumpuk** (Node 1, Node 2, Node n), masing-masing berisi sensor, ESP8266, relay, dan aktuator, untuk menunjukkan rancangan multi-node.
> - Kolom Network: router WiFi dan awan internet.
> - Kolom Cloud: satu kotak besar berlabel "VM (Docker Compose)" berisi empat kotak: Mosquitto, Backend Go (+ model RF), PostgreSQL, Caddy. Panah: Mosquitto ↔ Backend, Backend ↔ PostgreSQL, Caddy ↔ Backend.
> - Kolom UI: *browser* (website Vue) dan ikon Telegram.
> - Label panah: Node → Mosquitto "MQTT: telemetry"; Mosquitto → Node "MQTT: cmd"; Caddy ↔ *Browser* "HTTPS + WebSocket"; Backend → Telegram "Bot API".
> - Hapus semua label Blynk.

## 4.3 Topologi Jaringan & Protokol Komunikasi

### 4.3.1 Topologi Jaringan [DIUBAH]

Sistem menggunakan **topologi bintang (*star*)** pada dua tingkat. Pada tingkat lokal, setiap node terhubung ke router WiFi di lokasi akuariumnya sebagai *client* WiFi. Pada tingkat internet, seluruh node terhubung ke broker Mosquitto di VM, sehingga broker menjadi titik pusat pertukaran pesan. *Browser* pengguna terhubung ke VM melalui Caddy. Di dalam VM, keempat kontainer berkomunikasi melalui jaringan internal Docker, dan hanya port 1883 (MQTT), 80, serta 443 (HTTPS) yang dibuka ke internet. Node dan *browser* tidak pernah berkomunikasi secara langsung; seluruh perintah melewati *backend*.

> **Gambar 4.4 Topologi jaringan. [PERLU DIGAMBAR ULANG]**
> - Kiri: 2–3 node (ESP8266), masing-masing terhubung dengan garis putus-putus (WiFi) ke router. Boleh dua node di satu router dan satu node di router lain, untuk menunjukkan lokasi berbeda.
> - Tengah: awan internet.
> - Kanan: kotak VM berisi Mosquitto, Backend, PostgreSQL, Caddy, dengan keterangan port yang terbuka (1883, 80/443).
> - Bawah: ponsel dan laptop pengguna terhubung ke internet → Caddy. Tambahkan juga server Telegram yang menerima panah dari Backend.

### 4.3.2 Protokol Komunikasi [DIUBAH]

**Tabel 4.2** Protokol per jalur komunikasi

| Jalur | Protokol | Pola komunikasi |
|---|---|---|
| DS18B20 ↔ NodeMCU | 1-Wire | Digital |
| Sensor kekeruhan → NodeMCU | Sinyal analog (ADC) | – |
| NodeMCU → Relay | Sinyal digital GPIO | – |
| NodeMCU ↔ Router | WiFi IEEE 802.11 b/g/n | – |
| NodeMCU ↔ Mosquitto | MQTT di atas TCP/IP (port 1883, *username/password*) | *Publish/subscribe*, koneksi persisten |
| Backend ↔ Mosquitto | MQTT (jaringan internal Docker) | *Publish/subscribe* |
| Backend ↔ PostgreSQL | Protokol PostgreSQL di atas TCP (internal) | *Request–response* |
| *Browser* ↔ Caddy ↔ Backend | HTTPS (REST API) dan WSS (WebSocket) | REST: *request–response*; WebSocket: *push* dua arah |
| Backend → Telegram | HTTPS (Telegram Bot API) | *Request–response* |

MQTT dipilih untuk jalur node karena ringan, mendukung komunikasi dua arah melalui koneksi yang terus terbuka, dan memisahkan data setiap node melalui topik. WebSocket dipilih untuk website agar pembaruan data dan status aktuator tampil tanpa *refresh* halaman.

### 4.3.3 Rancangan Topik MQTT [BARU, menggantikan Rancangan Datastream Blynk]

**Tabel 4.3** Rancangan topik MQTT

| Topik | Arah | Isi pesan (JSON) | QoS | Retained |
|---|---|---|---|---|
| `aquarium/{node_id}/telemetry` | Node → Backend | `{"temp":27.31,"ntu":12.4}` | 0 | Tidak |
| `aquarium/{node_id}/state` | Node → Backend | `{"heater":0,"fan":1,"pump":0,"failsafe":false}` | 1 | Tidak |
| `aquarium/{node_id}/status` | Node/Broker → Backend | `online` atau `offline` (LWT) | 1 | Ya |
| `aquarium/{node_id}/cmd` | Backend → Node | `{"heater":0,"fan":1,"pump":0}` | 1 | Tidak |
| `aquarium/{node_id}/config` | Backend → Node | `{"mode":"auto"}` | 1 | Ya |
| `server/heartbeat` | Backend → semua node | `{"ts":1762171200}` | 0 | Tidak |

QoS 1 dipakai untuk perintah, status, dan konfigurasi supaya pesan dijamin sampai minimal satu kali. Data *telemetry* cukup QoS 0 karena dikirim berkala. Topik `config` dan `status` menggunakan *retained message* sehingga node yang baru menyala langsung menerima mode terakhirnya, dan *backend* langsung mengetahui status setiap node.

### 4.3.4 Rancangan REST API dan WebSocket [BARU]

**Tabel 4.4** Endpoint REST API

| Metode | Endpoint | Fungsi |
|---|---|---|
| POST | `/api/auth/login` | *Login*, menghasilkan JWT |
| GET | `/api/nodes` | Daftar node beserta status terakhir |
| GET | `/api/nodes/{id}` | Detail node (sensor, kelas, aktuator, mode) |
| PATCH | `/api/nodes/{id}` | Mengubah nama node |
| GET | `/api/nodes/{id}/telemetry?from=&to=` | Riwayat sensor dan kelas untuk grafik |
| PUT | `/api/nodes/{id}/mode` | Mengubah mode Auto/Manual |
| PUT | `/api/nodes/{id}/actuators` | Perintah aktuator manual |
| GET | `/api/nodes/{id}/logs` | Riwayat aksi aktuator dan notifikasi |
| GET | `/ws` | WebSocket untuk pembaruan *real-time* |

Seluruh *endpoint* kecuali *login* memerlukan JWT. Perintah manual dari website tidak dikirim langsung ke node, tetapi melalui *backend*, yang memeriksa mode node, mencatat log, lalu meneruskannya ke topik `cmd`.

### 4.3.5 Rancangan Basis Data [BARU]

**Tabel 4.5** Rancangan tabel basis data

| Tabel | Kolom utama | Keterangan |
|---|---|---|
| `users` | id, username, password_hash, telegram_chat_id | Akun pengguna website |
| `nodes` | node_id (PK), name, mode, online, last_seen, created_at | Daftar node; terisi otomatis saat node baru terhubung |
| `telemetry` | id, node_id (FK), ts, temp, ntu | Riwayat data sensor |
| `predictions` | id, node_id (FK), ts, class_temp, class_turbidity, class_overall | Hasil klasifikasi tiga model |
| `actuator_logs` | id, node_id (FK), ts, heater, fan, pump, source (auto/manual/failsafe) | Riwayat perubahan aktuator |
| `alerts` | id, node_id (FK), ts, level, message, sent | Riwayat notifikasi |

Tabel `telemetry` dan `predictions` diberi indeks pada kolom (node_id, ts) agar pengambilan data grafik per node tetap cepat ketika data bertambah.

> **Gambar 4.5 Entity Relationship Diagram (ERD). [BARU, PERLU DIGAMBAR]**
> Enam entitas sesuai Tabel 4.5. Relasi: `nodes` 1–N `telemetry`, `nodes` 1–N `predictions`, `nodes` 1–N `actuator_logs`, `nodes` 1–N `alerts`. Tabel `users` berdiri sendiri (seluruh pengguna dapat mengakses seluruh node pada prototipe ini). Bisa dibuat di dbdiagram.io atau draw.io.

### 4.3.6 Mekanisme Penambahan Node [BARU]

Penambahan akuarium baru dilakukan tanpa mengubah kode server:
1. *Firmware* yang sama di-*flash* ke ESP8266 baru. Kredensial MQTT perangkat sudah tertanam di *firmware*, sedangkan ACL broker menggunakan pola berbasis *client ID* sehingga setiap node hanya dapat mengakses topik `aquarium/{node_id}/...` miliknya sendiri.
2. Saat pertama kali menyala, node membuka portal WiFiManager. Pengguna memasukkan WiFi lokasi akuarium melalui ponsel.
3. Node terhubung ke broker dengan *client ID* berupa *node_id*, lalu mulai mengirim *telemetry*.
4. *Backend* mendeteksi *node_id* yang belum terdaftar dan menambahkannya ke tabel `nodes` dengan mode bawaan Manual, sehingga aktuator tidak langsung bergerak sebelum pengguna memeriksanya.
5. Node baru muncul di website. Pengguna dapat memberi nama dan mengaktifkan mode otomatis.

Kemampuan ini diuji dengan menjalankan satu node fisik bersama beberapa node simulasi, yaitu program yang mengirim *telemetry* dengan format topik yang sama.

## 4.4 Skematik Rangkaian Elektronika (Wiring Diagram)

**Tabel 4.6** Pemetaan pin sisi kendali

| Komponen | Pin komponen | Terhubung ke |
|---|---|---|
| Breadboard | Jalur + / Jalur − | 3V3 / GND NodeMCU |
| DS18B20 (modul) | VCC / GND / DAT | Jalur + / Jalur − / D5 |
| Sensor kekeruhan (modul) | V / G / A | Jalur + / Jalur − / A0 |
| Relay 4 channel | DC+ / DC− | Jalur + / Jalur − |
| | IN1 / IN2 / IN3 | D6 (pemanas) / D7 (kipas) / D1 (pompa) |
| | Jumper trigger | Posisi Low |

**Tabel 4.7** Jalur sisi beban (220 VAC)

| Aktuator | Jalur |
|---|---|
| Pemanas | Steker → fasa → COM1 → NO1 → stop kontak 1 lubang → pemanas |
| Kipas | Steker → fasa → COM2 → NO2 → stop kontak 1 lubang → adaptor DC → kipas |
| Pompa filter | Steker → fasa → COM3 → NO3 → stop kontak 1 lubang → pompa |

Kabel netral dan ground setiap stop kontak tidak diputus. Terminal NC tidak digunakan, sehingga aktuator dalam keadaan mati ketika relay tidak aktif.

NodeMCU mendapat daya 5 V melalui micro USB. Regulator internalnya menyediakan 3,3 V untuk kedua sensor dan sisi kendali relay. Aktuator memperoleh daya langsung dari jaringan listrik 220 VAC melalui stop kontak masing-masing.

> **Gambar 4.6 Wiring diagram.** Gambar tetap sama seperti laporan terakhir; hanya nomornya bergeser dari 4.5 menjadi 4.6.

## 4.5 Desain Antarmuka Aplikasi Web / Mobile (UI/UX Mockup) [DIUBAH TOTAL]

Antarmuka berupa website Vue.js yang responsif, sehingga satu desain dapat digunakan pada komputer maupun ponsel. Website terdiri atas halaman berikut:

**Tabel 4.8** Rancangan halaman website

| No | Halaman | Isi |
|---|---|---|
| 1 | *Login* | Formulir *username* dan *password* |
| 2 | Daftar Node (beranda) | Kartu untuk setiap node: nama, indikator *online*/*offline*, suhu dan NTU terakhir, label status *overall* berwarna (NORMAL hijau, CAUTION kuning, WARNING oranye, SEVERE merah), dan waktu pembaruan terakhir |
| 3 | Detail Node | *Gauge* suhu (0–40 °C) dan kekeruhan (0–200 NTU); label kelas suhu, kekeruhan, dan *overall*; sakelar mode Auto/Manual; tiga sakelar aktuator (pemanas, kipas, pompa) yang hanya menampilkan status pada mode Auto; indikator jika node sedang dalam mode *fail-safe* |
| 4 | Grafik Riwayat | Grafik garis ECharts untuk suhu dan NTU dengan pilihan rentang waktu (1 jam, 24 jam, 7 hari, 30 hari), dengan warna latar menandai kelas kondisi |
| 5 | Log & Notifikasi | Tabel riwayat aksi aktuator (waktu, aktuator, ON/OFF, sumber: auto/manual/*fail-safe*) dan riwayat peringatan |
| 6 | Pengaturan | Mengubah nama node dan menghubungkan akun Telegram untuk menerima notifikasi |

> **Gambar yang perlu dibuat untuk 4.5 [BARU, menggantikan mockup Blynk]:**
> - **Gambar 4.7** Mockup halaman Daftar Node (tampilan desktop), minimal 3 kartu node dengan status berbeda (misalnya NORMAL, WARNING, *offline*).
> - **Gambar 4.8** Mockup halaman Detail Node (tampilan desktop), termasuk *gauge*, tiga label kelas, sakelar mode, tiga sakelar aktuator, dan grafik riwayat.
> - **Gambar 4.9** Mockup halaman Detail Node (tampilan ponsel, *frame* 390×844), susunan komponen menjadi satu kolom.
> - **Gambar 4.10** Contoh pesan notifikasi bot Telegram, misalnya: "⚠️ WARNING – Akuarium 1: suhu 33,1 °C (Warning), kekeruhan 18 NTU (Normal). Kipas dinyalakan otomatis."
> - Alat yang disarankan: Figma. Mockup Blynk lama (Gambar 4.6 dan 4.7 di laporan terakhir) **dihapus**.

---

# DAFTAR PUSTAKA

[1] M. Nasir dan N. Natasya, "Sistem Monitoring Akuarium Berbasis Mikrokontroler dan Django Web Framework," *Jurnal Elektrika Borneo (JEB)*, vol. 6, no. 1, hlm. 25–28, Apr. 2020.

[2] M. Qolbu, A. Taqwa, dan M. Fadhli, "Penerapan Algoritma LightGBM untuk Penentuan Kondisi Air Akuarium Berdasarkan Parameter Kualitas Air," *MALCOM: Indonesian Journal of Machine Learning and Computer Science*, vol. 6, no. 4, hlm. 1896–1907, Okt. 2026, doi: 10.57152/malcom.v6i4.2933.

[3] P. O. Andini, H. Sunardi, dan Rachmansyah, "Perancangan Smart Aquarium Menggunakan Sensor Turbidity dan Sensor Ultrasonik pada Akuarium Ikan Mas Berbasis IoT," *Journal of Intelligent Networks and IoT Global*, vol. 2, no. 2, hlm. 85–98, 2024.

[4] L. Breiman, "Random Forests," *Machine Learning*, vol. 45, no. 1, hlm. 5–32, 2001.

[5] Y. T. Revantianto, P. Handoko, dan N. Uddin, "Pengembangan Prototype Sistem Pengatur Suhu Air Aquarium Berbasis IoT," Skripsi, Program Studi Informatika, Universitas Pembangunan Jaya, `[isi tahun]`.

[6] `[isi pembuat dataset]`, "`[isi judul dataset]`," `[isi platform, misalnya Kaggle]`, `[isi tahun]`. `[isi URL]`.

[7] M. M. Islam, M. A. Kashem, S. A. Alyami, dan M. A. Moni, "Monitoring water quality metrics of ponds with IoT sensors and machine learning to predict fish species survival," *Microprocessors and Microsystems*, vol. 102, art. no. 104930, 2023, doi: 10.1016/j.micpro.2023.104930.

[8] I. A. Adeleke, N. I. Nwulu, dan O. A. Ogbolumani, "A hybrid machine learning and embedded IoT-based water quality monitoring system," *Internet of Things*, vol. 22, art. no. 100774, 2023, doi: 10.1016/j.iot.2023.100774.

[9] O. O. Olanubi, T. T. Akano, dan O. S. Asaolu, "Design and development of an IoT-based intelligent water quality management system for aquaculture," *Journal of Electrical Systems and Information Technology*, vol. 11, art. no. 15, 2024, doi: 10.1186/s43067-024-00139-z.

[10] S. Kanwal, M. Abdullah, S. Kumar, S. Arshad, M. Shahroz, D. Zhang, dan D. Kumar, "An Optimal Internet of Things-Driven Intelligent Decision-Making System for Real-Time Fishpond Water Quality Monitoring and Species Survival," *Sensors*, vol. 24, no. 23, art. no. 7842, 2024, doi: 10.3390/s24237842.

[11] P. G. Arepalli dan K. J. Naik, "Water contamination analysis in IoT enabled aquaculture using deep learning based AODEGRU," *Ecological Informatics*, vol. 79, art. no. 102405, 2024, doi: 10.1016/j.ecoinf.2023.102405.

[12] R. P. Shete, A. M. Bongale, dan D. Dharrao, "IoT-enabled effective real-time water quality monitoring method for aquaculture," *MethodsX*, vol. 13, art. no. 102906, 2024, doi: 10.1016/j.mex.2024.102906.

[13] M. A. A. M. Hridoy, C. Bordin, A. Masood, dan K. Masood, "Predictive modelling of aquaculture water quality using IoT and advanced machine learning algorithms," *Results in Chemistry*, vol. 16, art. no. 102456, 2025, doi: 10.1016/j.rechem.2025.102456.

[14] A. Banks, E. Briggs, K. Borgendale, dan R. Gupta (Eds.), *MQTT Version 5.0*, OASIS Standard, Mar. 2019.
