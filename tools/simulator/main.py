"""
Simulator akuarium untuk proyek Smart Aquarium.
Menyimulasikan node sensor yang mengirim data ke broker MQTT.

Fitur:
- Koneksi MQTT dengan LWT, status online/offline.
- Pengiriman telemetry JSON setiap interval.
- Memproses perintah aktuator dari topik cmd dan membalas ke topik state.
- Mendukung skenario: normal, panas, dingin, keruh, ganda, mati.
"""

import argparse
import json
import logging
import os
import threading
import time

import paho.mqtt.client as mqtt
from paho.mqtt.enums import CallbackAPIVersion

# ---------------------------------------------------------------------------
# Konfigurasi logging: tampilkan waktu, nama modul, level, dan pesan
# ---------------------------------------------------------------------------
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(name)s] %(levelname)s %(message)s",
    datefmt="%H:%M:%S",
)

# ---------------------------------------------------------------------------
# Nilai awal sensor berdasarkan skenario (kontrak peran-2-simulator-uji.md)
# ---------------------------------------------------------------------------
NILAI_AWAL_SKENARIO: dict[str, dict[str, float]] = {
    "normal": {"temp": 27.0, "ntu": 10.0},
    "panas":  {"temp": 33.0, "ntu": 10.0},
    "dingin": {"temp": 22.0, "ntu": 10.0},
    "keruh":  {"temp": 27.0, "ntu": 80.0},
    "ganda":  {"temp": 17.0, "ntu": 300.0},
    # Skenario "mati": nilai awal sama dengan normal, lalu koneksi diputus paksa.
    "mati":   {"temp": 27.0, "ntu": 10.0},
}


# ---------------------------------------------------------------------------
# Laju perubahan sensor per interval — mudah disesuaikan jika perlu
# ---------------------------------------------------------------------------
# Suhu (°C per interval):
LAJU_NAIK_SUHU: float = 0.2   # heater ON, fan OFF
LAJU_TURUN_SUHU: float = 0.2  # fan ON, heater OFF
# Kekeruhan (NTU per interval):
LAJU_TURUN_NTU: float = 0.5   # pump ON

# Batas nilai sensor yang diizinkan (kontrak validasi telemetry)
BATAS_TEMP_MIN: float = 0.0
BATAS_TEMP_MAX: float = 50.0
BATAS_NTU_MIN: float = 0.0


# ---------------------------------------------------------------------------
# Kelas NodeSimulator: satu objek = satu akuarium virtual
# ---------------------------------------------------------------------------
class NodeSimulator:
    """
    Mewakili satu node akuarium virtual.

    Tanggung jawab:
    - Menyambung ke broker MQTT dengan LWT yang benar.
    - Menerbitkan status 'online' saat tersambung.
    - Mengirim telemetry JSON setiap interval.
    - Subscribe ke topik cmd, config, dan server/heartbeat.
    - Memproses perintah cmd dan membalas dengan state terbaru.
    - Mendukung pemutusan koneksi paksa untuk skenario 'mati'.
    """

    def __init__(
        self,
        node_id: str,
        skenario: str,
        interval: int,
        host: str,
        port: int,
        username: str,
        password: str,
    ) -> None:
        self.node_id = node_id
        self.skenario = skenario
        self.interval = interval
        self.host = host
        self.port = port
        self.log = logging.getLogger(node_id)

        # Nilai sensor: disalin dari tabel skenario
        nilai_awal = NILAI_AWAL_SKENARIO.get(skenario, NILAI_AWAL_SKENARIO["normal"])
        self.temp: float = nilai_awal["temp"]
        self.ntu: float = nilai_awal["ntu"]

        # Status aktuator awal (semua mati)
        self.heater: int = 0
        self.fan: int = 0
        self.pump: int = 0
        self.failsafe: bool = False

        # Tanda apakah node sudah berhasil tersambung ke broker
        self._tersambung = threading.Event()

        # ---------------------------------------------------------------------------
        # Inisialisasi klien MQTT dengan CallbackAPIVersion.VERSION2 (paho-mqtt v2)
        # ---------------------------------------------------------------------------
        self._client = mqtt.Client(
            callback_api_version=CallbackAPIVersion.VERSION2,
            client_id=node_id,
        )

        # Pasang kredensial jika tersedia
        if username:
            self._client.username_pw_set(username, password)

        # ---------------------------------------------------------------------------
        # Daftarkan Last Will and Testament (LWT)
        # Jika koneksi terputus tiba-tiba, broker akan menerbitkan pesan ini.
        # Topik  : aquarium/{node_id}/status
        # Payload: "offline"
        # QoS    : 1 (sesuai kontrak)
        # Retain : True (sesuai kontrak, agar klien baru langsung tahu statusnya)
        # ---------------------------------------------------------------------------
        topik_status = f"aquarium/{node_id}/status"
        self._client.will_set(
            topic=topik_status,
            payload="offline",
            qos=1,
            retain=True,
        )

        # Pasang fungsi callback
        self._client.on_connect = self._on_connect
        self._client.on_disconnect = self._on_disconnect
        self._client.on_message = self._on_message

    # ------------------------------------------------------------------
    # Properti pembantu: nama topik
    # ------------------------------------------------------------------

    @property
    def _topik_telemetry(self) -> str:
        return f"aquarium/{self.node_id}/telemetry"

    @property
    def _topik_status(self) -> str:
        return f"aquarium/{self.node_id}/status"

    @property
    def _topik_state(self) -> str:
        return f"aquarium/{self.node_id}/state"

    @property
    def _topik_cmd(self) -> str:
        return f"aquarium/{self.node_id}/cmd"

    @property
    def _topik_config(self) -> str:
        return f"aquarium/{self.node_id}/config"

    # ------------------------------------------------------------------
    # Metode internal: proses perintah aktuator dari topik cmd
    # ------------------------------------------------------------------

    def _proses_cmd(self, payload_str: str) -> None:
        """
        Parsing payload JSON dari topik cmd, perbarui state aktuator,
        lalu publikasikan state terbaru ke topik state.

        Format perintah (kontrak): {"heater":0,"fan":1,"pump":0}
        Nilai 1 = nyala, 0 = mati.
        """
        try:
            data = json.loads(payload_str)
        except json.JSONDecodeError as exc:
            self.log.warning("Payload cmd tidak valid (bukan JSON): %s — %s", payload_str, exc)
            return

        # Ambil nilai tiap aktuator; abaikan field yang tidak dikenal
        berubah = False
        for aktuator in ("heater", "fan", "pump"):
            if aktuator not in data:
                continue
            nilai_baru = int(data[aktuator])
            if nilai_baru not in (0, 1):
                self.log.warning("Nilai %s tidak valid: %s (harus 0 atau 1)", aktuator, nilai_baru)
                continue
            nilai_lama = getattr(self, aktuator)
            if nilai_baru != nilai_lama:
                setattr(self, aktuator, nilai_baru)
                self.log.info(
                    "Aktuator berubah: %s %d -> %d",
                    aktuator, nilai_lama, nilai_baru,
                )
                berubah = True

        if berubah:
            self.log.info("State aktuator diperbarui — mengirim state terbaru.")
        else:
            self.log.info("Tidak ada perubahan aktuator — tetap mengirim state sebagai konfirmasi.")

        # Kirim state terbaru; selalu dilakukan seperti perilaku firmware asli
        self._kirim_state()

    # ------------------------------------------------------------------
    # Callback: dipanggil saat koneksi ke broker berhasil
    # ------------------------------------------------------------------

    def _on_connect(
        self,
        client: mqtt.Client,
        userdata,
        connect_flags,
        reason_code,
        properties,
    ) -> None:
        if reason_code.is_failure:
            self.log.error("Gagal terhubung ke broker: %s", reason_code)
            return

        self.log.info("Terhubung ke broker %s:%d", self.host, self.port)

        # 1. Kirim status "online" (retained, sesuai kontrak)
        client.publish(
            topic=self._topik_status,
            payload="online",
            qos=0,
            retain=True,
        )
        self.log.info("Status diterbitkan: online -> %s", self._topik_status)

        # 2. Kirim state aktuator awal (semua 0, belum ada perintah)
        self._kirim_state()

        # 3. Subscribe ke topik yang perlu didengarkan
        #    - cmd    : perintah aktuator dari backend (QoS 1 sesuai kontrak)
        #    - config : konfigurasi mode auto/manual
        #    - server/heartbeat : detak jantung backend
        topik_langganan = [
            (self._topik_cmd, 1),
            (self._topik_config, 1),
            ("server/heartbeat", 0),
        ]
        client.subscribe(topik_langganan)
        self.log.info(
            "Berlangganan topik: %s",
            ", ".join(t for t, _ in topik_langganan),
        )

        # Beri tahu thread utama bahwa koneksi sudah siap
        self._tersambung.set()

    # ------------------------------------------------------------------
    # Callback: dipanggil saat koneksi terputus
    # ------------------------------------------------------------------

    def _on_disconnect(
        self,
        client: mqtt.Client,
        userdata,
        disconnect_flags,
        reason_code,
        properties,
    ) -> None:
        # Paho v2 menyediakan dua cara membedakan jenis disconnect:
        # - disconnect_flags.is_disconnect_packet_from_server: True jika broker
        #   yang mengirim paket DISCONNECT (misalnya kick/timeout dari sisi broker).
        # - reason_code.is_failure: True jika bukan disconnect yang bersih (rc != 0).
        #
        # Kombinasi keduanya:
        #   is_failure=False, is_disconnect_packet_from_server=False
        #     → kita yang memanggil disconnect() secara normal (clean shutdown).
        #   is_failure=True atau is_disconnect_packet_from_server=True
        #     → disconnect tidak direncanakan, kemungkinan akan direkoneksi otomatis.
        if reason_code.is_failure:
            self.log.warning(
                "Terputus tidak bersih dari broker (kode: %s, dari_server: %s)",
                reason_code,
                disconnect_flags.is_disconnect_packet_from_server,
            )
        else:
            self.log.info("Terputus dari broker dengan bersih.")
        self._tersambung.clear()

    # ------------------------------------------------------------------
    # Callback: dipanggil saat ada pesan masuk dari topik yang dilanggani
    # ------------------------------------------------------------------

    def _on_message(
        self,
        client: mqtt.Client,
        userdata,
        message: mqtt.MQTTMessage,
    ) -> None:
        payload_str = message.payload.decode("utf-8", errors="replace")
        self.log.info(
            "Pesan diterima <- %s : %s",
            message.topic,
            payload_str,
        )

        # Routing berdasarkan topik yang diterima
        if message.topic == self._topik_cmd:
            # Perintah aktuator dari backend: perbarui state dan balas ke topik state
            self._proses_cmd(payload_str)

        elif message.topic == self._topik_config:
            # Konfigurasi mode (auto/manual) — dicatat untuk kelengkapan log
            # Logika mode penuh akan diimplementasi jika dibutuhkan oleh pengujian integrasi
            self.log.info("Config diterima: %s", payload_str)

        elif message.topic == "server/heartbeat":
            # Detak jantung backend — cukup dicatat, tidak perlu tindakan
            self.log.debug("Heartbeat diterima: %s", payload_str)

    # ------------------------------------------------------------------
    # Metode internal: perbarui nilai sensor berdasarkan status aktuator
    # ------------------------------------------------------------------

    def _perbarui_sensor(self) -> None:
        """
        Hitung perubahan temp dan NTU berdasarkan status aktuator saat ini,
        lalu terapkan ke self.temp dan self.ntu.

        Aturan (dari spesifikasi Peran 2):
          - heater ON, fan OFF  → temp naik LAJU_NAIK_SUHU per interval
          - fan ON, heater OFF  → temp turun LAJU_TURUN_SUHU per interval
          - heater ON, fan ON   → tidak ada perubahan temp (saling menetralkan)
          - heater OFF, fan OFF → tidak ada perubahan temp
          - pump ON             → NTU turun LAJU_TURUN_NTU per interval
          - pump OFF            → tidak ada perubahan NTU

        Nilai sensor dibatasi sesuai rentang valid kontrak:
          - temp: 0–50 °C
          - NTU : ≥ 0
        """
        temp_lama = self.temp
        ntu_lama = self.ntu

        # --- Perubahan suhu ---
        if self.heater == 1 and self.fan == 0:
            self.temp += LAJU_NAIK_SUHU
        elif self.fan == 1 and self.heater == 0:
            self.temp -= LAJU_TURUN_SUHU
        # heater ON + fan ON, atau keduanya OFF → suhu tidak berubah akibat aktuator

        # --- Perubahan kekeruhan ---
        if self.pump == 1:
            self.ntu -= LAJU_TURUN_NTU

        # --- Terapkan batas nilai sensor ---
        self.temp = max(BATAS_TEMP_MIN, min(BATAS_TEMP_MAX, self.temp))
        self.ntu = max(BATAS_NTU_MIN, self.ntu)

        # --- Logging perubahan (hanya jika ada yang berubah) ---
        delta_temp = round(self.temp - temp_lama, 4)
        delta_ntu = round(self.ntu - ntu_lama, 4)

        if delta_temp != 0 or delta_ntu != 0:
            self.log.info(
                "Sensor diperbarui: temp %.2f -> %.2f (%+.2f), ntu %.2f -> %.2f (%+.2f) "
                "[heater=%d fan=%d pump=%d]",
                temp_lama, self.temp, delta_temp,
                ntu_lama, self.ntu, delta_ntu,
                self.heater, self.fan, self.pump,
            )

    # ------------------------------------------------------------------
    # Metode internal: kirim state aktuator saat ini
    # ------------------------------------------------------------------

    def _kirim_state(self) -> None:
        """
        Kirim status aktuator aktual ke topik state.
        Dipanggil saat: pertama kali tersambung dan setiap relay berubah akibat cmd.
        """
        payload = json.dumps({
            "heater": self.heater,
            "fan": self.fan,
            "pump": self.pump,
            "failsafe": self.failsafe,
        })
        self._client.publish(
            topic=self._topik_state,
            payload=payload,
            qos=0,
            retain=False,
        )
        self.log.info("State dikirim -> %s : %s", self._topik_state, payload)

    # ------------------------------------------------------------------
    # Metode internal: kirim satu paket telemetry
    # ------------------------------------------------------------------

    def _kirim_telemetry(self) -> None:
        """
        Buat payload JSON dan terbitkan ke topik telemetry.
        Format payload (kontrak): {"temp":27.31,"ntu":12.4}
        Nilai dibulatkan ke 2 desimal agar tidak terlalu panjang.
        """
        payload = json.dumps({
            "temp": round(self.temp, 2),
            "ntu": round(self.ntu, 2),
        })
        self._client.publish(
            topic=self._topik_telemetry,
            payload=payload,
            qos=0,
            retain=False,
        )
        self.log.info("Telemetry -> %s : %s", self._topik_telemetry, payload)

    # ------------------------------------------------------------------
    # Metode publik: jalankan node ini
    # ------------------------------------------------------------------

    def jalankan(self) -> None:
        """
        Sambungkan ke broker lalu mulai loop pengiriman telemetry.
        Fungsi ini berjalan hingga program dihentikan (Ctrl+C).
        """
        self.log.info(
            "Memulai koneksi ke %s:%d (skenario=%s, interval=%ds)",
            self.host, self.port, self.skenario, self.interval,
        )

        # Sambungkan; paho akan memanggil _on_connect secara asinkron
        self._client.connect(host=self.host, port=self.port, keepalive=60)

        # Jalankan loop jaringan MQTT di thread terpisah agar tidak memblokir
        self._client.loop_start()

        # Tunggu hingga benar-benar terhubung sebelum mengirim data
        if not self._tersambung.wait(timeout=10):
            self.log.error(
                "Gagal terhubung dalam 10 detik. Periksa alamat broker."
            )
            self._client.loop_stop()
            return

        # Jika skenario "mati", putus koneksi paksa setelah 2 detik
        # supaya LWT otomatis dikirim oleh broker.
        if self.skenario == "mati":
            self.log.info("Skenario 'mati': memutus koneksi paksa dalam 2 detik...")
            time.sleep(2)
            # loop_stop(force=True) memutus tanpa DISCONNECT yang bersih
            self._client.loop_stop(force=True)
            self.log.info("Koneksi diputus paksa. LWT 'offline' akan dikirim broker.")
            return

        # Loop utama: perbarui sensor lalu kirim telemetry setiap interval
        try:
            while True:
                self._perbarui_sensor()
                self._kirim_telemetry()
                time.sleep(self.interval)
        except KeyboardInterrupt:
            self.log.info("Dihentikan oleh pengguna.")
        finally:
            # Urutan shutdown yang benar:
            # 1. disconnect() mengantrekan paket DISCONNECT ke buffer internal paho.
            # 2. loop_stop() bersifat blocking — menunggu thread jaringan selesai.
            #    Thread tersebut memflush paket DISCONNECT ke socket sebelum berhenti.
            # Dengan urutan ini, broker menerima DISCONNECT yang bersih
            # sehingga tidak memicu reconnect otomatis yang tidak diinginkan.
            self._client.disconnect()
            self._client.loop_stop()
            self.log.info("Koneksi ditutup dengan bersih.")


# ---------------------------------------------------------------------------
# Fungsi utama: parsing argumen dan menjalankan semua node di thread terpisah
# ---------------------------------------------------------------------------

def parse_argumen() -> argparse.Namespace:
    """Parsing argumen baris perintah sesuai spesifikasi kontrak."""
    parser = argparse.ArgumentParser(
        description="Simulator node akuarium untuk proyek Smart Aquarium.",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument(
        "--nodes",
        type=int,
        default=3,
        metavar="N",
        help="Jumlah node yang disimulasikan.",
    )
    parser.add_argument(
        "--scenario",
        default="normal",
        choices=list(NILAI_AWAL_SKENARIO.keys()),
        help="Skenario simulasi.",
    )
    parser.add_argument(
        "--interval",
        type=int,
        default=5,
        metavar="DETIK",
        help="Jarak waktu antarkirim telemetry (detik).",
    )
    parser.add_argument(
        "--host",
        default=os.environ.get("MQTT_HOST", "localhost"),
        help="Alamat broker MQTT. Bawaan: env MQTT_HOST atau 'localhost'.",
    )
    parser.add_argument(
        "--port",
        type=int,
        default=1883,
        help="Port broker MQTT.",
    )
    parser.add_argument(
        "--username",
        default=os.environ.get("MQTT_USERNAME", ""),
        help="Username MQTT. Bawaan: env MQTT_USERNAME.",
    )
    parser.add_argument(
        "--password",
        default=os.environ.get("MQTT_PASSWORD", ""),
        help="Password MQTT. Bawaan: env MQTT_PASSWORD.",
    )
    return parser.parse_args()


def main() -> None:
    args = parse_argumen()

    # Buat node_id: aq-sim001, aq-sim002, ...
    node_ids = [f"aq-sim{i:03d}" for i in range(1, args.nodes + 1)]

    logging.info(
        "Memulai %d node | skenario=%s | interval=%ds | broker=%s:%d",
        args.nodes, args.scenario, args.interval, args.host, args.port,
    )

    # Jika hanya 1 node, jalankan langsung di thread utama
    if args.nodes == 1:
        NodeSimulator(
            node_id=node_ids[0],
            skenario=args.scenario,
            interval=args.interval,
            host=args.host,
            port=args.port,
            username=args.username,
            password=args.password,
        ).jalankan()
        return

    # Untuk lebih dari 1 node, jalankan setiap node di thread terpisah
    threads: list[threading.Thread] = []
    for node_id in node_ids:
        sim = NodeSimulator(
            node_id=node_id,
            skenario=args.scenario,
            interval=args.interval,
            host=args.host,
            port=args.port,
            username=args.username,
            password=args.password,
        )
        t = threading.Thread(target=sim.jalankan, name=node_id, daemon=True)
        threads.append(t)
        t.start()
        # Jeda kecil antarnoda agar koneksi tidak serentak
        time.sleep(0.2)

    # Tunggu semua thread selesai (atau Ctrl+C)
    try:
        for t in threads:
            t.join()
    except KeyboardInterrupt:
        logging.info("Program dihentikan. Menunggu semua node selesai...")


if __name__ == "__main__":
    main()
