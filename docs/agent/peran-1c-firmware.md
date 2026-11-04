# Peran 1C: firmware

Agent menulis kode, sedangkan upload dan pengujian dilakukan pengguna di alat. Setelah setiap tahap, berikan langkah uji di alat dan hasil yang diharapkan. Jangan mengaku firmware sudah teruji sebelum pengguna melaporkannya.

## File

```text
firmware/smart_aquarium/
├─ smart_aquarium.ino
├─ config.example.h
└─ config.h             tidak masuk Git
```

Board Arduino IDE: `NodeMCU 1.0 (ESP-12E Module)`. Library: OneWire, DallasTemperature, PubSubClient, ArduinoJson versi 7, WiFiManager.

`config.example.h`:

```cpp
#pragma once
#define MQTT_HOST "192.168.1.10"
#define MQTT_PORT 1883
#define MQTT_USER "node"
#define MQTT_PASS "ganti"
#define TELEMETRY_INTERVAL_MS 60000UL
#define STATE_INTERVAL_MS 60000UL
#define HEARTBEAT_TIMEOUT_MS 180000UL
#define NTU_A 0.0f
#define NTU_B 0.0f
#define RELAY_ON LOW
#define RELAY_OFF HIGH
```

## Pin dan titik awal

| Fungsi | Pin |
| --- | --- |
| DS18B20 | D5 |
| Turbidity | A0 |
| Relay pemanas IN1 | D6 |
| Relay kipas IN2 | D7 |
| Relay pompa IN3 | D1 |

Kode pembacaan sensor yang sudah berjalan:

```cpp
#include <OneWire.h>
#include <DallasTemperature.h>

#define PIN_SUHU        D5
#define PIN_TURBIDITY   A0

OneWire oneWire(PIN_SUHU);
DallasTemperature sensorSuhu(&oneWire);

void setup() {
  Serial.begin(115200);
  delay(200);
  sensorSuhu.begin();

  Serial.println();
  Serial.print("Jumlah DS18B20 terdeteksi: ");
  Serial.println(sensorSuhu.getDeviceCount());
}

int bacaTurbidity(int sampel) {
  long total = 0;
  for (int i = 0; i < sampel; i++) {
    total += analogRead(PIN_TURBIDITY);  // 0 - 1023
    delay(10);
  }
  return total / sampel;
}

void loop() {
  // Suhu
  sensorSuhu.requestTemperatures();
  float suhu = sensorSuhu.getTempCByIndex(0);

  // Turbidity
  int rawTurb = bacaTurbidity(20);
  float tegangan = rawTurb * (3.3 / 1023.0);

  // Tampilkan
  Serial.print("Suhu: ");
  if (suhu == DEVICE_DISCONNECTED_C) {
    Serial.print("ERROR (cek kabel DS18B20)");
  } else {
    Serial.print(suhu, 2);
    Serial.print(" C");
  }

  Serial.print(" | Turbidity raw: ");
  Serial.print(rawTurb);
  Serial.print(" (");
  Serial.print(tegangan, 2);
  Serial.println(" V)");

  delay(1000);
}
```

## Tahap kerja

Kerjakan satu tahap per sesi.

1. **Sensor dan relay.** Hitung `NTU = NTU_A * tegangan + NTU_B`, batasi minimal 0, dan sediakan perintah Serial untuk mengetes relay.
2. **WiFi.** Pakai WiFiManager dengan nama portal `SmartAquarium-{chip ID}` dan timeout portal 180 detik.
3. **MQTT.** Pakai format `node_id` dari kontrak, `setBufferSize(512)`, LWT QoS 1 retained, lalu `online` retained. Subscribe `cmd` dan `config` QoS 1, serta `server/heartbeat` QoS 0.
4. **Telemetry.** Pakai timer `millis()`, bukan `delay()` panjang. Lewati pengiriman kalau suhu `DEVICE_DISCONNECTED_C` atau 85 °C pada pembacaan pertama setelah boot.
5. **Perintah.** Terima `cmd`, ubah relay, kirim `state`. Kirim `state` juga saat tersambung dan setiap `STATE_INTERVAL_MS`.
6. **Fail-safe.** Aktif saat mode Auto dan MQTT putus atau heartbeat hilang lebih dari 180 detik.
7. **Sambung ulang.** WiFi dan MQTT tersambung lagi tanpa reset, dan `client.loop()` tetap dipanggil sesering mungkin.

Contoh koneksi LWT:

```cpp
client.connect(nodeId, MQTT_USER, MQTT_PASS, statusTopic, 1, true, "offline");
client.publish(statusTopic, "online", true);
```

## Aturan keselamatan

- Tulis `RELAY_OFF` ke pin sebelum `pinMode(pin, OUTPUT)`.
- Semua relay OFF saat boot.
- PubSubClient hanya mengirim publish QoS 0. Jangan mencoba publish QoS 1.
- Fail-safe: pemanas ON kalau suhu di bawah 25 °C, kipas ON kalau suhu di atas 30 °C, pompa ON kalau kekeruhan minimal 25 NTU.
- Saat fail-safe aktif dan sensor suhu gagal, pemanas dan kipas harus OFF.
- Mode Manual tidak memakai fail-safe.
- Kirim `failsafe:true` pada `state` selama mode lokal aktif, lalu kembali mengikuti server saat heartbeat kembali.

## Uji di alat

- [ ] Saat boot, tidak ada aktuator yang berkedip.
- [ ] Setiap relay mengikuti `cmd`.
- [ ] Node muncul setelah restart tanpa flash ulang.
- [ ] DS18B20 dilepas: tidak ada telemetry suhu palsu.
- [ ] Broker dimatikan saat Auto: fail-safe aktif dan pulih setelah broker kembali.
- [ ] Hotspot dimatikan lalu dinyalakan: node tersambung sendiri.
