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
