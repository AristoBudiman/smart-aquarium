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

OneWire oneWire(PIN_SUHU);
DallasTemperature sensorSuhu(&oneWire);

bool stateHeater = false;
bool stateFan = false;
bool statePump = false;

unsigned long lastSensorReadMs = 0;

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

void printHelp() {
  Serial.println("=== Perintah Tes Relay (Serial Monitor) ===");
  Serial.println("  h1 / heater on  : Nyalakan pemanas (IN1 / D6)");
  Serial.println("  h0 / heater off : Matikan pemanas  (IN1 / D6)");
  Serial.println("  f1 / fan on     : Nyalakan kipas   (IN2 / D7)");
  Serial.println("  f0 / fan off    : Matikan kipas    (IN2 / D7)");
  Serial.println("  p1 / pump on    : Nyalakan pompa   (IN3 / D1)");
  Serial.println("  p0 / pump off   : Matikan pompa    (IN3 / D1)");
  Serial.println("  off             : Matikan semua relay");
  Serial.println("  status / ?      : Tampilkan status relay & bantuan");
  Serial.println("===========================================");
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

  Serial.println();
  Serial.println("=== Smart Aquarium - Tahap 1: Sensor & Relay ===");
  Serial.print("Jumlah DS18B20 terdeteksi: ");
  Serial.println(sensorSuhu.getDeviceCount());
  printRelayState();
  printHelp();
}

void loop() {
  handleSerialCommand();

  unsigned long now = millis();
  if (now - lastSensorReadMs >= SENSOR_PRINT_INTERVAL_MS) {
    lastSensorReadMs = now;
    readAndPrintSensors();
  }
}

