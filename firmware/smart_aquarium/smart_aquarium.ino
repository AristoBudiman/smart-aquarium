#include <ESP8266WiFi.h>
#include <WiFiManager.h>
#include <PubSubClient.h>
#include <OneWire.h>
#include <DallasTemperature.h>

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

#define SENSOR_PRINT_INTERVAL_MS 2000UL
#define MQTT_RETRY_INTERVAL_MS   5000UL
#define WIFI_PORTAL_TIMEOUT_S    180

OneWire oneWire(PIN_SUHU);
DallasTemperature sensorSuhu(&oneWire);

WiFiClient espClient;
PubSubClient mqttClient(espClient);

char chipIdHex[16];
char nodeId[24];
char portalSsid[32];

char statusTopic[64];
char cmdTopic[64];
char configTopic[64];
const char* heartbeatTopic = "server/heartbeat";

bool stateHeater = false;
bool stateFan = false;
bool statePump = false;

unsigned long lastSensorReadMs = 0;
unsigned long lastMqttRetryMs = 0;

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

void initNodeIdentity() {
  snprintf(chipIdHex, sizeof(chipIdHex), "%06x", ESP.getChipId());
  snprintf(nodeId, sizeof(nodeId), "aq-%s", chipIdHex);
  snprintf(portalSsid, sizeof(portalSsid), "SmartAquarium-%s", chipIdHex);

  snprintf(statusTopic, sizeof(statusTopic), "aquarium/%s/status", nodeId);
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
    Serial.println("[WiFi] Gagal tersambung (timeout portal habis). Akan mencoba ulang di background.");
  } else {
    Serial.print("[WiFi] Tersambung! IP: ");
    Serial.println(WiFi.localIP());
  }

  WiFi.setAutoReconnect(true);
  WiFi.persistent(true);
}

void mqttCallback(char* topic, byte* payload, unsigned int length) {
  Serial.print("[MQTT Masuk] Topik: ");
  Serial.print(topic);
  Serial.print(" | Payload: ");
  for (unsigned int i = 0; i < length; i++) {
    Serial.print((char)payload[i]);
  }
  Serial.println();
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
  } else {
    Serial.print("[MQTT] Gagal, rc=");
    Serial.print(mqttClient.state());
    Serial.println(" (coba lagi dalam 5 detik)");
  }
}

int bacaTurbidity(int sampel) {
  long total = 0;
  for (int i = 0; i < sampel; i++) {
    total += analogRead(PIN_TURBIDITY); // 0 - 1023
    delay(10);
  }
  return total / sampel;
}

void printRelayState() {
  Serial.print("Status Relay -> Pemanas (D6): ");
  Serial.print(stateHeater ? "ON" : "OFF");
  Serial.print(" | Kipas (D7): ");
  Serial.print(stateFan ? "ON" : "OFF");
  Serial.print(" | Pompa (D1): ");
  Serial.println(statePump ? "ON" : "OFF");
}

void printConnectionInfo() {
  Serial.print("Node ID: ");
  Serial.print(nodeId);
  Serial.print(" | WiFi: ");
  if (WiFi.status() == WL_CONNECTED) {
    Serial.print("Terhubung (");
    Serial.print(WiFi.localIP());
    Serial.print(")");
  } else {
    Serial.print("Terputus");
  }
  Serial.print(" | MQTT: ");
  Serial.println(mqttClient.connected() ? "Terhubung" : "Terputus");
}

void printHelp() {
  Serial.println("=== Perintah Tes Relay & Info (Serial Monitor) ===");
  Serial.println("  h1 / heater on  : Nyalakan pemanas (IN1 / D6)");
  Serial.println("  h0 / heater off : Matikan pemanas  (IN1 / D6)");
  Serial.println("  f1 / fan on     : Nyalakan kipas   (IN2 / D7)");
  Serial.println("  f0 / fan off    : Matikan kipas    (IN2 / D7)");
  Serial.println("  p1 / pump on    : Nyalakan pompa   (IN3 / D1)");
  Serial.println("  p0 / pump off   : Matikan pompa    (IN3 / D1)");
  Serial.println("  off             : Matikan semua relay");
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

  if (cmd == "h1" || cmd == "heater on") {
    stateHeater = true;
    setRelay(PIN_HEATER, true);
    Serial.println("[CMD] Pemanas (IN1/D6) -> ON");
    printRelayState();
  } else if (cmd == "h0" || cmd == "heater off") {
    stateHeater = false;
    setRelay(PIN_HEATER, false);
    Serial.println("[CMD] Pemanas (IN1/D6) -> OFF");
    printRelayState();
  } else if (cmd == "f1" || cmd == "fan on") {
    stateFan = true;
    setRelay(PIN_FAN, true);
    Serial.println("[CMD] Kipas (IN2/D7) -> ON");
    printRelayState();
  } else if (cmd == "f0" || cmd == "fan off") {
    stateFan = false;
    setRelay(PIN_FAN, false);
    Serial.println("[CMD] Kipas (IN2/D7) -> OFF");
    printRelayState();
  } else if (cmd == "p1" || cmd == "pump on") {
    statePump = true;
    setRelay(PIN_PUMP, true);
    Serial.println("[CMD] Pompa (IN3/D1) -> ON");
    printRelayState();
  } else if (cmd == "p0" || cmd == "pump off") {
    statePump = false;
    setRelay(PIN_PUMP, false);
    Serial.println("[CMD] Pompa (IN3/D1) -> OFF");
    printRelayState();
  } else if (cmd == "off" || cmd == "all0") {
    stateHeater = false;
    stateFan = false;
    statePump = false;
    setRelay(PIN_HEATER, false);
    setRelay(PIN_FAN, false);
    setRelay(PIN_PUMP, false);
    Serial.println("[CMD] Semua relay -> OFF");
    printRelayState();
  } else if (cmd == "status" || cmd == "?") {
    printConnectionInfo();
    printRelayState();
    printHelp();
  } else {
    Serial.print("[CMD] Perintah tidak dikenal: ");
    Serial.println(cmd);
    printHelp();
  }
}

void readAndPrintSensors() {
  // Suhu
  sensorSuhu.requestTemperatures();
  float suhu = sensorSuhu.getTempCByIndex(0);

  // Turbidity
  int rawTurb = bacaTurbidity(20);
  float tegangan = rawTurb * (3.3f / 1023.0f);
  float ntu = (NTU_A * tegangan) + NTU_B;
  if (ntu < 0.0f) {
    ntu = 0.0f;
  }

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
  Serial.print(" | Tegangan: ");
  Serial.print(tegangan, 2);
  Serial.print(" V | NTU: ");
  Serial.print(ntu, 2);
  Serial.println(" NTU");
}

void setup() {
  initRelays();

  Serial.begin(115200);
  delay(200);
  sensorSuhu.begin();

  initNodeIdentity();

  Serial.println();
  Serial.println("=== Smart Aquarium - Tahap 2 & 3: WiFi & MQTT ===");
  Serial.print("Node ID : ");
  Serial.println(nodeId);
  Serial.print("Jumlah DS18B20 terdeteksi: ");
  Serial.println(sensorSuhu.getDeviceCount());
  printRelayState();

  setupWiFi();
  setupMqtt();
  connectMqtt();
  lastMqttRetryMs = millis();

  printConnectionInfo();
  printHelp();
}

void loop() {
  if (WiFi.status() == WL_CONNECTED) {
    if (!mqttClient.connected()) {
      unsigned long now = millis();
      if (now - lastMqttRetryMs >= MQTT_RETRY_INTERVAL_MS) {
        lastMqttRetryMs = now;
        connectMqtt();
      }
    } else {
      mqttClient.loop();
    }
  }

  handleSerialCommand();

  unsigned long now = millis();
  if (now - lastSensorReadMs >= SENSOR_PRINT_INTERVAL_MS) {
    lastSensorReadMs = now;
    readAndPrintSensors();
  }
}
