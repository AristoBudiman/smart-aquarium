#include <ESP8266WiFi.h>
#include <WiFiManager.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>
#include <OneWire.h>
#include <DallasTemperature.h>
#include <math.h>

#if __has_include("config.h")
#include "config.h"
#else
#include "config.example.h"
#endif

#define PIN_SUHU        D5
#define PIN_TURBIDITY   A0
#define PIN_HEATER      D6
#define PIN_FAN         D7
#define PIN_PUMP        D1

#define SENSOR_READ_INTERVAL_MS  2000UL
#define DS18B20_CONV_WAIT_MS     800UL
#define TEMP_INVALID_CUTOFF_MS   10000UL
#define TURB_SAMPLE_COUNT        20
#define TURB_SAMPLE_INTERVAL_MS  20UL
#define MQTT_RETRY_INTERVAL_MS   5000UL
#define WIFI_RETRY_INTERVAL_MS   10000UL
#define WIFI_PORTAL_TIMEOUT_S    180

OneWire oneWire(PIN_SUHU);
DallasTemperature sensorSuhu(&oneWire);

WiFiClient espClient;
PubSubClient mqttClient(espClient);

char chipIdHex[16];
char nodeId[24];
char portalSsid[32];

char statusTopic[64];
char telemetryTopic[64];
char stateTopic[64];
char cmdTopic[64];
char configTopic[64];
const char* heartbeatTopic = "server/heartbeat";

// Status aktuator & mode
bool stateHeater = false;
bool stateFan = false;
bool statePump = false;
bool isAutoMode = false;      // Mode awal Manual sampai config diterima
bool failsafeActive = false;

// Pembacaan sensor non-blocking
bool sensorCycleActive = false;
unsigned long lastSensorCycleMs = 0;
unsigned long tempRequestMs = 0;
unsigned long lastTurbSampleMs = 0;
long turbSampleSum = 0;
int turbSampleIndex = 0;

bool firstTempReadAfterBoot = true;
bool hasSensorReading = false;
bool latestTempValid = false;
bool tempInvalidActive = false;
unsigned long tempInvalidSinceMs = 0;
float latestTemp = DEVICE_DISCONNECTED_C;
int latestRawTurb = 0;
float latestVoltage = 0.0f;
float latestNtu = 0.0f;

// Timer MQTT, WiFi, Telemetry, State, dan Heartbeat
unsigned long lastMqttRetryMs = 0;
unsigned long lastWifiRetryMs = 0;
unsigned long lastTelemetryMs = 0;
unsigned long lastStatePublishMs = 0;
unsigned long lastHeartbeatMs = 0;
bool firstTelemetrySent = false;
bool lastWifiConnected = false;

void setRelay(uint8_t pin, bool on) {
  digitalWrite(pin, on ? RELAY_ON : RELAY_OFF);
}

void initRelays() {
  // Tulis RELAY_OFF ke pin sebelum pinMode(pin, OUTPUT) agar relay tidak berkedip saat boot
  digitalWrite(PIN_HEATER, RELAY_OFF);
  digitalWrite(PIN_FAN, RELAY_OFF);
  digitalWrite(PIN_PUMP, RELAY_OFF);

  pinMode(PIN_HEATER, OUTPUT);
  pinMode(PIN_FAN, OUTPUT);
  pinMode(PIN_PUMP, OUTPUT);

  stateHeater = false;
  stateFan = false;
  statePump = false;
}

bool applyRelays(bool heaterOn, bool fanOn, bool pumpOn) {
  bool changed = (stateHeater != heaterOn) || (stateFan != fanOn) || (statePump != pumpOn);
  stateHeater = heaterOn;
  stateFan = fanOn;
  statePump = pumpOn;

  setRelay(PIN_HEATER, stateHeater);
  setRelay(PIN_FAN, stateFan);
  setRelay(PIN_PUMP, statePump);
  return changed;
}

void initNodeIdentity() {
  snprintf(chipIdHex, sizeof(chipIdHex), "%06x", ESP.getChipId());
  snprintf(nodeId, sizeof(nodeId), "aq-%s", chipIdHex);
  snprintf(portalSsid, sizeof(portalSsid), "SmartAquarium-%s", chipIdHex);

  snprintf(statusTopic, sizeof(statusTopic), "aquarium/%s/status", nodeId);
  snprintf(telemetryTopic, sizeof(telemetryTopic), "aquarium/%s/telemetry", nodeId);
  snprintf(stateTopic, sizeof(stateTopic), "aquarium/%s/state", nodeId);
  snprintf(cmdTopic, sizeof(cmdTopic), "aquarium/%s/cmd", nodeId);
  snprintf(configTopic, sizeof(configTopic), "aquarium/%s/config", nodeId);
}

void setupWiFi() {
  WiFi.mode(WIFI_STA);
  WiFiManager wm;
  wm.setConfigPortalTimeout(WIFI_PORTAL_TIMEOUT_S);

  Serial.print("[WiFi] Menghubungkan WiFi atau membuka portal AP: ");
  Serial.println(portalSsid);

  bool res = wm.autoConnect(portalSsid);
  if (!res) {
    Serial.println("[WiFi] Timeout portal habis. Akan mencoba sambung ulang di background.");
  } else {
    Serial.print("[WiFi] Tersambung! IP: ");
    Serial.println(WiFi.localIP());
    lastWifiConnected = true;
  }

  WiFi.setAutoReconnect(true);
  WiFi.persistent(true);
}

void printRelayState() {
  Serial.print("Status Relay -> Pemanas (D6): ");
  Serial.print(stateHeater ? "ON" : "OFF");
  Serial.print(" | Kipas (D7): ");
  Serial.print(stateFan ? "ON" : "OFF");
  Serial.print(" | Pompa (D1): ");
  Serial.print(statePump ? "ON" : "OFF");
  Serial.print(" | Mode: ");
  Serial.print(isAutoMode ? "AUTO" : "MANUAL");
  Serial.print(" | Fail-safe: ");
  Serial.println(failsafeActive ? "AKTIF" : "TIDAK");
}

void publishState() {
  if (!mqttClient.connected()) {
    return;
  }

  char payload[96];
  snprintf(
    payload,
    sizeof(payload),
    "{\"heater\":%d,\"fan\":%d,\"pump\":%d,\"failsafe\":%s}",
    stateHeater ? 1 : 0,
    stateFan ? 1 : 0,
    statePump ? 1 : 0,
    failsafeActive ? "true" : "false"
  );

  if (mqttClient.publish(stateTopic, payload, false)) {
    lastStatePublishMs = millis();
    Serial.print("[MQTT Publish] ");
    Serial.print(stateTopic);
    Serial.print(" -> ");
    Serial.println(payload);
  } else {
    Serial.println("[MQTT Publish] Gagal mengirim state");
  }
}

void publishTelemetry() {
  if (!mqttClient.connected()) {
    return;
  }

  if (!hasSensorReading || !latestTempValid) {
    Serial.println("[Telemetry] Dilewati: pembacaan suhu DS18B20 belum/tidak valid");
    return;
  }

  char payload[64];
  snprintf(
    payload,
    sizeof(payload),
    "{\"temp\":%.2f,\"ntu\":%.2f}",
    latestTemp,
    latestNtu
  );

  if (mqttClient.publish(telemetryTopic, payload, false)) {
    Serial.print("[MQTT Publish] ");
    Serial.print(telemetryTopic);
    Serial.print(" -> ");
    Serial.println(payload);
  } else {
    Serial.println("[MQTT Publish] Gagal mengirim telemetry");
  }
}

bool applyFailsafeRelays() {
  bool targetHeater = false;
  bool targetFan = false;
  bool targetPump = false;

  if (hasSensorReading && latestTempValid) {
    targetHeater = (latestTemp < 25.0f);
    targetFan = (latestTemp > 30.0f);
  } else {
    // Saat fail-safe aktif dan sensor suhu gagal, pemanas dan kipas harus OFF
    targetHeater = false;
    targetFan = false;
  }

  if (hasSensorReading) {
    targetPump = (latestNtu >= 25.0f);
  }

  if (applyRelays(targetHeater, targetFan, targetPump)) {
    Serial.println("[Fail-safe] Relay disesuaikan oleh aturan lokal fail-safe:");
    printRelayState();
    publishState();
    return true;
  }
  return false;
}

void evaluateFailsafe() {
  if (!isAutoMode) {
    if (failsafeActive) {
      failsafeActive = false;
      applyRelays(false, false, false);
      Serial.println("[Fail-safe] Nonaktif karena mode berpindah ke MANUAL. Semua relay OFF.");
      publishState();
    }
  return;
}

  unsigned long now = millis();
  bool mqttDown = !mqttClient.connected();
  bool heartbeatTimeout = (now - lastHeartbeatMs) > HEARTBEAT_TIMEOUT_MS;
  bool shouldFailsafe = mqttDown || heartbeatTimeout;

  if (shouldFailsafe && !failsafeActive) {
    failsafeActive = true;
    Serial.print("[Fail-safe] AKTIF! Alasan: ");
    if (mqttDown) {
      Serial.println("MQTT terputus");
    } else {
      Serial.println("Heartbeat server hilang > 180 detik");
    }
    if (!applyFailsafeRelays()) {
      publishState();
    }
  } else if (!shouldFailsafe && failsafeActive) {
    failsafeActive = false;
    Serial.println("[Fail-safe] PULIH! Kembali mengikuti server");
    publishState();
  } else if (failsafeActive) {
    applyFailsafeRelays();
  }
}

// Mode Manual tidak diubah: pengguna tetap memegang kendali penuh atas relay.
void enforceTempSensorCutoff() {
  if (!isAutoMode || failsafeActive || !tempInvalidActive) {
    return;
  }

  if ((millis() - tempInvalidSinceMs) >= TEMP_INVALID_CUTOFF_MS) {
    if (stateHeater || stateFan) {
      applyRelays(false, false, statePump);
      Serial.println("[Sensor] DS18B20 tidak terbaca > 10 dtk: pemanas & kipas dimatikan");
      publishState();
    }
  }
}

void handleCmdMessage(const byte* payload, unsigned int length) {
  JsonDocument doc;
  DeserializationError err = deserializeJson(doc, payload, length);
  if (err) {
    Serial.print("[MQTT] Gagal parse JSON cmd: ");
    Serial.println(err.c_str());
    return;
  }

  if (failsafeActive) {
    Serial.println("[MQTT] cmd diabaikan sementara karena fail-safe masih aktif (menunggu heartbeat)");
    return;
  }

  bool targetHeater = stateHeater;
  bool targetFan = stateFan;
  bool targetPump = statePump;

  if (doc["heater"].is<int>()) {
    targetHeater = (doc["heater"].as<int>() != 0);
  } else if (doc["heater"].is<bool>()) {
    targetHeater = doc["heater"].as<bool>();
  }

  if (doc["fan"].is<int>()) {
    targetFan = (doc["fan"].as<int>() != 0);
  } else if (doc["fan"].is<bool>()) {
    targetFan = doc["fan"].as<bool>();
  }

  if (doc["pump"].is<int>()) {
    targetPump = (doc["pump"].as<int>() != 0);
  } else if (doc["pump"].is<bool>()) {
    targetPump = doc["pump"].as<bool>();
  }

  // Mode Manual tidak dibatasi: pengguna tetap memegang kendali penuh.
  if (isAutoMode && hasSensorReading && !latestTempValid) {
    targetHeater = false;
    targetFan = false;
    Serial.println("[MQTT] Mode AUTO & sensor suhu gagal: pemanas & kipas dipaksa OFF, pompa mengikuti cmd");
  }

  applyRelays(targetHeater, targetFan, targetPump);
  Serial.println("[MQTT] cmd diterapkan ke relay:");
  printRelayState();
  publishState();
}

void handleConfigMessage(const byte* payload, unsigned int length) {
  JsonDocument doc;
  DeserializationError err = deserializeJson(doc, payload, length);
  if (err) {
    Serial.print("[MQTT] Gagal parse JSON config: ");
    Serial.println(err.c_str());
    return;
  }

  if (!doc["mode"].is<const char*>()) {
    Serial.println("[MQTT] Field 'mode' tidak ditemukan pada config");
    return;
  }

  const char* modeStr = doc["mode"].as<const char*>();
  if (strcmp(modeStr, "auto") == 0) {
    if (!isAutoMode) {
      isAutoMode = true;
      // Mulai hitung batas waktu heartbeat sejak masuk mode Auto jika belum ada heartbeat baru
      if (mqttClient.connected() && (millis() - lastHeartbeatMs > HEARTBEAT_TIMEOUT_MS)) {
        lastHeartbeatMs = millis();
      }
      Serial.println("[Config] Mode diubah ke AUTO");
    } else {
      Serial.println("[Config] Mode tetap AUTO");
    }
    evaluateFailsafe();
  } else if (strcmp(modeStr, "manual") == 0) {
    isAutoMode = false;
    Serial.println("[Config] Mode diubah ke MANUAL");
    evaluateFailsafe();
  } else {
    Serial.print("[Config] Mode tidak dikenal: ");
    Serial.println(modeStr);
  }
}

void handleHeartbeatMessage(const byte* payload, unsigned int length) {
  JsonDocument doc;
  DeserializationError err = deserializeJson(doc, payload, length);
  if (err) {
    Serial.print("[MQTT] Gagal parse JSON heartbeat: ");
    Serial.println(err.c_str());
    return;
  }

  lastHeartbeatMs = millis();
  if (failsafeActive) {
    evaluateFailsafe();
  }
}

void mqttCallback(char* topic, byte* payload, unsigned int length) {
  Serial.print("[MQTT Masuk] Topik: ");
  Serial.print(topic);
  Serial.print(" | Payload: ");
  for (unsigned int i = 0; i < length; i++) {
    Serial.print((char)payload[i]);
  }
  Serial.println();

  if (strcmp(topic, cmdTopic) == 0) {
    handleCmdMessage(payload, length);
  } else if (strcmp(topic, configTopic) == 0) {
    handleConfigMessage(payload, length);
  } else if (strcmp(topic, heartbeatTopic) == 0) {
    handleHeartbeatMessage(payload, length);
  }
}

void setupMqtt() {
  mqttClient.setServer(MQTT_HOST, MQTT_PORT);
  mqttClient.setBufferSize(512);
  mqttClient.setCallback(mqttCallback);
}

void connectMqtt() {
  if (WiFi.status() != WL_CONNECTED) {
    return;
  }
  if (mqttClient.connected()) {
    return;
  }

  Serial.print("[MQTT] Menghubungkan ke ");
  Serial.print(MQTT_HOST);
  Serial.print(":");
  Serial.print(MQTT_PORT);
  Serial.print(" sebagai ");
  Serial.print(nodeId);
  Serial.println(" ...");

  const char* user = (strlen(MQTT_USER) > 0) ? MQTT_USER : nullptr;
  const char* pass = (strlen(MQTT_USER) > 0) ? MQTT_PASS : nullptr;

  if (mqttClient.connect(nodeId, user, pass, statusTopic, 1, true, "offline")) {
    Serial.println("[MQTT] Tersambung ke broker!");
    mqttClient.publish(statusTopic, "online", true);
    mqttClient.subscribe(cmdTopic, 1);
    mqttClient.subscribe(configTopic, 1);
    mqttClient.subscribe(heartbeatTopic, 0);

    Serial.print("[MQTT] Publish retained 'online' ke: ");
    Serial.println(statusTopic);
    Serial.print("[MQTT] Subscribed (QoS 1): ");
    Serial.print(cmdTopic);
    Serial.print(", ");
    Serial.println(configTopic);
    Serial.print("[MQTT] Subscribed (QoS 0): ");
    Serial.println(heartbeatTopic);

    // Reset timer heartbeat saat koneksi broker berhasil dipulihkan
    lastHeartbeatMs = millis();
    if (failsafeActive) {
      failsafeActive = false;
      Serial.println("[Fail-safe] PULIH setelah koneksi broker MQTT kembali!");
    }

    // Kirim state saat tersambung sesuai kontrak
    publishState();
  } else {
    Serial.print("[MQTT] Gagal, rc=");
    Serial.print(mqttClient.state());
    Serial.println(" (coba lagi dalam 5 detik)");
  }
}

void updateSensorsNonBlocking() {
  unsigned long now = millis();

  if (!sensorCycleActive) {
    if (now - lastSensorCycleMs >= SENSOR_READ_INTERVAL_MS) {
      sensorCycleActive = true;
      sensorSuhu.requestTemperatures();
      tempRequestMs = now;
      turbSampleSum = 0;
      turbSampleIndex = 0;
      lastTurbSampleMs = now;
    }
    return;
  }

  // Ambil sampel turbidity secara bertahap tanpa delay()
  if (turbSampleIndex < TURB_SAMPLE_COUNT && (now - lastTurbSampleMs >= TURB_SAMPLE_INTERVAL_MS)) {
    lastTurbSampleMs = now;
    turbSampleSum += analogRead(PIN_TURBIDITY);
    turbSampleIndex++;
  }

  // Tunggu konversi DS18B20 selesai (~750ms) dan 20 sampel turbidity terkumpul
  if (turbSampleIndex >= TURB_SAMPLE_COUNT && (now - tempRequestMs >= DS18B20_CONV_WAIT_MS)) {
    float suhu = sensorSuhu.getTempCByIndex(0);

    bool valid = true;
    bool skippedFirst85 = false;
    if (suhu == DEVICE_DISCONNECTED_C) {
      valid = false;
      if (!tempInvalidActive) {
        tempInvalidActive = true;
        tempInvalidSinceMs = now;
      }
    } else if (firstTempReadAfterBoot && fabs(suhu - 85.0f) < 0.01f) {
      valid = false;
      skippedFirst85 = true;
      firstTempReadAfterBoot = false;
    } else {
      firstTempReadAfterBoot = false;
      tempInvalidActive = false;
    }

    latestTempValid = valid;
    if (valid) {
      latestTemp = suhu;
    }

    latestRawTurb = turbSampleSum / TURB_SAMPLE_COUNT;
    latestVoltage = latestRawTurb * (3.3f / 1023.0f);
    float ntu = (NTU_A * latestVoltage) + NTU_B;
    if (ntu < 0.0f) {
      ntu = 0.0f;
    }
    latestNtu = ntu;
    hasSensorReading = true;

    // Tampilkan pembacaan di Serial Monitor
    Serial.print("Suhu: ");
    if (!latestTempValid) {
      if (skippedFirst85) {
        Serial.print("SKIP (85.00 C pembacaan pertama setelah boot)");
      } else {
        Serial.print("ERROR (cek kabel DS18B20)");
      }
    } else {
      Serial.print(latestTemp, 2);
      Serial.print(" C");
    }

    Serial.print(" | Turbidity raw: ");
    Serial.print(latestRawTurb);
    Serial.print(" | Tegangan: ");
    Serial.print(latestVoltage, 2);
    Serial.print(" V | NTU: ");
    Serial.print(latestNtu, 2);
    Serial.print(" NTU | Mode: ");
    Serial.print(isAutoMode ? "AUTO" : "MANUAL");
    if (failsafeActive) {
      Serial.print(" [FAILSAFE]");
    }
    Serial.println();

    sensorCycleActive = false;
    lastSensorCycleMs = now;

    // Evaluasi aturan fail-safe setiap kali pembacaan sensor diperbarui
    if (failsafeActive) {
      applyFailsafeRelays();
    }

    // Kirim telemetry pertama segera setelah ada pembacaan suhu valid saat sudah terkoneksi
    if (mqttClient.connected() && latestTempValid && !firstTelemetrySent) {
      firstTelemetrySent = true;
      lastTelemetryMs = now;
      publishTelemetry();
    }
  }
}

void printConnectionInfo() {
  Serial.print("Node ID: ");
  Serial.print(nodeId);
  Serial.print(" | Mode: ");
  Serial.print(isAutoMode ? "AUTO" : "MANUAL");
  Serial.print(" | Fail-safe: ");
  Serial.print(failsafeActive ? "AKTIF" : "TIDAK");
  Serial.print(" | WiFi: ");
  if (WiFi.status() == WL_CONNECTED) {
    Serial.print("Terhubung (");
    Serial.print(WiFi.localIP());
    Serial.print(")");
  } else {
    Serial.print("Terputus");
  }
  Serial.print(" | MQTT: ");
  Serial.print(mqttClient.connected() ? "Terhubung" : "Terputus");
  Serial.print(" | Usia Heartbeat: ");
  Serial.print((millis() - lastHeartbeatMs) / 1000UL);
  Serial.println(" dtk");
}

void printHelp() {
  Serial.println("=== Perintah Tes Relay & Info (Serial Monitor) ===");
  Serial.println("  h1 / heater on  : Nyalakan pemanas (IN1 / D6) [Hanya MANUAL]");
  Serial.println("  h0 / heater off : Matikan pemanas  (IN1 / D6) [Hanya MANUAL]");
  Serial.println("  f1 / fan on     : Nyalakan kipas   (IN2 / D7) [Hanya MANUAL]");
  Serial.println("  f0 / fan off    : Matikan kipas    (IN2 / D7) [Hanya MANUAL]");
  Serial.println("  p1 / pump on    : Nyalakan pompa   (IN3 / D1) [Hanya MANUAL]");
  Serial.println("  p0 / pump off   : Matikan pompa    (IN3 / D1) [Hanya MANUAL]");
  Serial.println("  off             : Matikan semua relay         [Hanya MANUAL]");
  Serial.println("  status / ?      : Tampilkan status node, relay & bantuan");
  Serial.println("==================================================");
}

void handleSerialCommand() {
  if (!Serial.available()) {
    return;
  }

  String cmd = Serial.readStringUntil('\n');
  cmd.trim();
  cmd.toLowerCase();

  if (cmd.length() == 0) {
    return;
  }

  if (cmd == "status" || cmd == "?") {
    printConnectionInfo();
    printRelayState();
    printHelp();
    return;
  }

  bool isRelayCmd =
    (cmd == "h1" || cmd == "heater on" ||
     cmd == "h0" || cmd == "heater off" ||
     cmd == "f1" || cmd == "fan on" ||
     cmd == "f0" || cmd == "fan off" ||
     cmd == "p1" || cmd == "pump on" ||
     cmd == "p0" || cmd == "pump off" ||
     cmd == "off" || cmd == "all0");

  if (isRelayCmd && isAutoMode) {
    Serial.println("[CMD] Ditolak: Kontrol relay lewat Serial dinonaktifkan saat mode AUTO.");
    return;
  }

  if (cmd == "h1" || cmd == "heater on") {
    applyRelays(true, stateFan, statePump);
    Serial.println("[CMD] Pemanas (IN1/D6) -> ON");
    printRelayState();
    publishState();
  } else if (cmd == "h0" || cmd == "heater off") {
    applyRelays(false, stateFan, statePump);
    Serial.println("[CMD] Pemanas (IN1/D6) -> OFF");
    printRelayState();
    publishState();
  } else if (cmd == "f1" || cmd == "fan on") {
    applyRelays(stateHeater, true, statePump);
    Serial.println("[CMD] Kipas (IN2/D7) -> ON");
    printRelayState();
    publishState();
  } else if (cmd == "f0" || cmd == "fan off") {
    applyRelays(stateHeater, false, statePump);
    Serial.println("[CMD] Kipas (IN2/D7) -> OFF");
    printRelayState();
    publishState();
  } else if (cmd == "p1" || cmd == "pump on") {
    applyRelays(stateHeater, stateFan, true);
    Serial.println("[CMD] Pompa (IN3/D1) -> ON");
    printRelayState();
    publishState();
  } else if (cmd == "p0" || cmd == "pump off") {
    applyRelays(stateHeater, stateFan, false);
    Serial.println("[CMD] Pompa (IN3/D1) -> OFF");
    printRelayState();
    publishState();
  } else if (cmd == "off" || cmd == "all0") {
    applyRelays(false, false, false);
    Serial.println("[CMD] Semua relay -> OFF");
    printRelayState();
    publishState();
  } else {
    Serial.print("[CMD] Perintah tidak dikenal: ");
    Serial.println(cmd);
    printHelp();
  }
}

void maintainConnections() {
  unsigned long now = millis();
  bool wifiConnected = (WiFi.status() == WL_CONNECTED);

  if (wifiConnected && !lastWifiConnected) {
    lastWifiConnected = true;
    Serial.print("[WiFi] Tersambung kembali! IP: ");
    Serial.println(WiFi.localIP());
  } else if (!wifiConnected && lastWifiConnected) {
    lastWifiConnected = false;
    Serial.println("[WiFi] Koneksi terputus! Mencoba menyambung ulang...");
  }

  if (!wifiConnected) {
    if (now - lastWifiRetryMs >= WIFI_RETRY_INTERVAL_MS) {
      lastWifiRetryMs = now;
      WiFi.reconnect();
    }
    return;
  }

  if (!mqttClient.connected()) {
    if (now - lastMqttRetryMs >= MQTT_RETRY_INTERVAL_MS) {
      lastMqttRetryMs = now;
      connectMqtt();
    }
  } else {
    mqttClient.loop();
  }
}

void setup() {
  initRelays();

  Serial.begin(115200);
  delay(200);

  sensorSuhu.begin();
  sensorSuhu.setWaitForConversion(false); // Non-blocking DS18B20 conversion

  initNodeIdentity();

  Serial.println();
  Serial.println("=== Smart Aquarium - Firmware Lengkap (Tahap 1-7 + sensor cutoff) ===");
  Serial.print("Node ID : ");
  Serial.println(nodeId);
  Serial.print("Jumlah DS18B20 terdeteksi: ");
  Serial.println(sensorSuhu.getDeviceCount());
  printRelayState();

  setupWiFi();
  setupMqtt();

  unsigned long now = millis();
  lastHeartbeatMs = now;
  lastMqttRetryMs = now;
  lastWifiRetryMs = now;
  lastTelemetryMs = now;
  lastStatePublishMs = now;

  connectMqtt();

  printConnectionInfo();
  printHelp();
}

void loop() {
  maintainConnections();
  handleSerialCommand();
  updateSensorsNonBlocking();
  evaluateFailsafe();
  enforceTempSensorCutoff();

  unsigned long now = millis();
  if (mqttClient.connected()) {
    if (now - lastTelemetryMs >= TELEMETRY_INTERVAL_MS) {
      lastTelemetryMs = now;
      publishTelemetry();
    }

    if (now - lastStatePublishMs >= STATE_INTERVAL_MS) {
      lastStatePublishMs = now;
      publishState();
    }
  }
}
