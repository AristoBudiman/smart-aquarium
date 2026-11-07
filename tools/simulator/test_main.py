"""
Unit test simulator Smart Aquarium tanpa broker MQTT asli.
"""

import json
import os
import unittest
from unittest.mock import MagicMock, patch

import main as sim_mod
from main import (
    BATAS_NTU_MIN,
    BATAS_TEMP_MAX,
    BATAS_TEMP_MIN,
    LAJU_NAIK_SUHU,
    LAJU_TURUN_NTU,
    LAJU_TURUN_SUHU,
    NILAI_AWAL_SKENARIO,
    NodeSimulator,
    buat_id_node,
    parse_argumen,
)


class TestIdentitasDanKonfigurasi(unittest.TestCase):
    def test_buat_id_node(self) -> None:
        self.assertEqual(buat_id_node(1), ["aq-sim001"])
        self.assertEqual(buat_id_node(3), ["aq-sim001", "aq-sim002", "aq-sim003"])

    @patch("main.mqtt.Client")
    def test_client_id_lwt_dan_kredensial(self, mock_client_cls: MagicMock) -> None:
        mock_instance = MagicMock()
        mock_client_cls.return_value = mock_instance

        node = NodeSimulator(
            node_id="aq-sim001",
            skenario="normal",
            interval=5,
            host="broker.example.com",
            port=1883,
            username="node",
            password="secret",
        )

        mock_client_cls.assert_called_once_with(
            callback_api_version=sim_mod.CallbackAPIVersion.VERSION2,
            client_id="aq-sim001",
        )
        mock_instance.username_pw_set.assert_called_once_with("node", "secret")
        mock_instance.will_set.assert_called_once_with(
            topic="aquarium/aq-sim001/status",
            payload="offline",
            qos=1,
            retain=True,
        )
        self.assertEqual(node.temp, 27.0)
        self.assertEqual(node.ntu, 10.0)

    @patch("main.mqtt.Client")
    def test_on_connect_online_state_dan_subscribe(self, mock_client_cls: MagicMock) -> None:
        mock_instance = MagicMock()
        mock_client_cls.return_value = mock_instance

        node = NodeSimulator(
            node_id="aq-sim002",
            skenario="normal",
            interval=5,
            host="localhost",
            port=1883,
            username="",
            password="",
        )

        reason_code = MagicMock()
        reason_code.is_failure = False

        node._on_connect(mock_instance, None, None, reason_code, None)

        # Pastikan online retained dikirim
        mock_instance.publish.assert_any_call(
            topic="aquarium/aq-sim002/status",
            payload="online",
            qos=0,
            retain=True,
        )
        # Pastikan state awal dikirim saat tersambung
        expected_state = json.dumps({
            "heater": 0,
            "fan": 0,
            "pump": 0,
            "failsafe": False,
        })
        mock_instance.publish.assert_any_call(
            topic="aquarium/aq-sim002/state",
            payload=expected_state,
            qos=0,
            retain=False,
        )
        # Pastikan subscribe cmd QoS 1, config QoS 1, server/heartbeat QoS 0
        mock_instance.subscribe.assert_called_once_with([
            ("aquarium/aq-sim002/cmd", 1),
            ("aquarium/aq-sim002/config", 1),
            ("server/heartbeat", 0),
        ])
        self.assertTrue(node._tersambung.is_set())


class TestSkenarioDanSensor(unittest.TestCase):
    @patch("main.mqtt.Client")
    def test_enam_skenario_nilai_awal(self, _mock_client_cls: MagicMock) -> None:
        ekspektasi = {
            "normal": (27.0, 10.0),
            "panas": (33.0, 10.0),
            "dingin": (22.0, 10.0),
            "keruh": (27.0, 80.0),
            "ganda": (17.0, 300.0),
            "mati": (27.0, 10.0),
        }
        self.assertEqual(set(NILAI_AWAL_SKENARIO.keys()), set(ekspektasi.keys()))
        for skenario, (temp_awal, ntu_awal) in ekspektasi.items():
            node = NodeSimulator("aq-sim001", skenario, 5, "localhost", 1883, "", "")
            self.assertEqual(node.temp, temp_awal, f"temp salah untuk {skenario}")
            self.assertEqual(node.ntu, ntu_awal, f"ntu salah untuk {skenario}")

    @patch("main.mqtt.Client")
    def test_perubahan_sensor_saat_aktuator_menyala(self, _mock_client_cls: MagicMock) -> None:
        node = NodeSimulator("aq-sim001", "normal", 2, "localhost", 1883, "", "")

        # Pemanas menaikkan suhu (+0.3)
        node.heater = 1
        node.fan = 0
        node.pump = 0
        node._perbarui_sensor()
        self.assertAlmostEqual(node.temp, 27.0 + LAJU_NAIK_SUHU)
        self.assertAlmostEqual(node.ntu, 10.0)

        # Kipas menurunkan suhu (-0.3)
        node.heater = 0
        node.fan = 1
        node._perbarui_sensor()
        self.assertAlmostEqual(node.temp, 27.0)

        # Pemanas + kipas menyala bersamaan saling meniadakan
        node.heater = 1
        node.fan = 1
        node._perbarui_sensor()
        self.assertAlmostEqual(node.temp, 27.0)

        # Pompa menurunkan NTU (-5.0) dan dibatasi minimal 0
        node.heater = 0
        node.fan = 0
        node.pump = 1
        node._perbarui_sensor()
        self.assertAlmostEqual(node.ntu, 10.0 - LAJU_TURUN_NTU)
        node._perbarui_sensor()
        self.assertAlmostEqual(node.ntu, 0.0)
        node._perbarui_sensor()
        self.assertEqual(node.ntu, BATAS_NTU_MIN)

        # Batas suhu min dan max
        node.temp = 49.9
        node.heater = 1
        node.fan = 0
        node._perbarui_sensor()
        self.assertEqual(node.temp, BATAS_TEMP_MAX)

        node.temp = 0.1
        node.heater = 0
        node.fan = 1
        node._perbarui_sensor()
        self.assertEqual(node.temp, BATAS_TEMP_MIN)


class TestCmdDanState(unittest.TestCase):
    @patch("main.mqtt.Client")
    def test_cmd_mengubah_aktuator_dan_publish_state(self, mock_client_cls: MagicMock) -> None:
        mock_instance = MagicMock()
        mock_client_cls.return_value = mock_instance

        node = NodeSimulator("aq-sim001", "panas", 2, "localhost", 1883, "", "")
        node._proses_cmd('{"heater":0,"fan":1,"pump":1}')

        self.assertEqual(node.heater, 0)
        self.assertEqual(node.fan, 1)
        self.assertEqual(node.pump, 1)

        mock_instance.publish.assert_called_once_with(
            topic="aquarium/aq-sim001/state",
            payload=json.dumps({"heater": 0, "fan": 1, "pump": 1, "failsafe": False}),
            qos=0,
            retain=False,
        )

    @patch("main.mqtt.Client")
    def test_cmd_tidak_valid_aman_tanpa_crash(self, mock_client_cls: MagicMock) -> None:
        mock_instance = MagicMock()
        mock_client_cls.return_value = mock_instance

        node = NodeSimulator("aq-sim001", "normal", 2, "localhost", 1883, "", "")
        node.heater = 1
        node.fan = 0
        node.pump = 1

        # Bukan JSON -> diabaikan tanpa publish
        node._proses_cmd("bukan-json")
        mock_instance.publish.assert_not_called()

        # Bukan objek JSON (list / angka) -> diabaikan tanpa publish
        node._proses_cmd("[1, 2, 3]")
        node._proses_cmd("123")
        mock_instance.publish.assert_not_called()

        # Nilai bukan 0/1 (mis. 2, "abc", True, None) diabaikan, field lain yang valid diterapkan,
        # field yang tidak dikirim tetap memakai nilai lama.
        node._proses_cmd('{"heater":2,"fan":"nyala","pump":0}')
        self.assertEqual(node.heater, 1)  # tetap nilai lama karena 2 tidak valid
        self.assertEqual(node.fan, 0)     # tetap nilai lama karena string tidak valid
        self.assertEqual(node.pump, 0)    # berubah dari 1 ke 0

        # Boolean True/False pada JSON harus diabaikan (harus int 0/1)
        node._proses_cmd('{"heater":false,"fan":true}')
        self.assertEqual(node.heater, 1)
        self.assertEqual(node.fan, 0)


class TestSkenarioMatiDanLoop(unittest.TestCase):
    @patch("main.time.sleep", return_value=None)
    @patch("main.mqtt.Client")
    def test_skenario_mati_tutup_socket_tanpa_disconnect(
        self, mock_client_cls: MagicMock, _mock_sleep: MagicMock
    ) -> None:
        mock_instance = MagicMock()
        mock_sock = MagicMock()
        mock_instance.socket.return_value = mock_sock
        mock_client_cls.return_value = mock_instance

        node = NodeSimulator("aq-sim001", "mati", 2, "localhost", 1883, "", "")
        node._tersambung.set()

        node.jalankan()

        # Pastikan disconnect() TIDAK dipanggil agar LWT terpicu di broker
        mock_instance.disconnect.assert_not_called()
        # Pastikan loop_stop() dipanggil tanpa argumen force
        mock_instance.loop_stop.assert_called_once_with()
        # Pastikan socket ditutup paksa
        mock_sock.close.assert_called_once_with()

    @patch("main.mqtt.Client")
    def test_telemetry_pertama_nilai_awal_dan_state_berkala(
        self, mock_client_cls: MagicMock
    ) -> None:
        mock_instance = MagicMock()
        mock_client_cls.return_value = mock_instance

        node = NodeSimulator("aq-sim001", "panas", 2, "localhost", 1883, "", "")
        node.fan = 1  # Kipas menyala
        node._tersambung.set()

        # Jalankan 2 siklus lalu hentikan dengan KeyboardInterrupt pada sleep ke-2
        with patch("main.time.sleep", side_effect=[None, KeyboardInterrupt()]):
            node.jalankan()

        # Ambil semua panggilan publish ke telemetry dan state
        telemetry_calls = [
            json.loads(call.kwargs["payload"])
            for call in mock_instance.publish.call_args_list
            if call.kwargs.get("topic") == "aquarium/aq-sim001/telemetry"
        ]
        state_calls = [
            json.loads(call.kwargs["payload"])
            for call in mock_instance.publish.call_args_list
            if call.kwargs.get("topic") == "aquarium/aq-sim001/state"
        ]

        # Siklus 1 mengirim nilai awal 33.0, siklus 2 mengirim 32.7
        self.assertEqual(len(telemetry_calls), 2)
        self.assertEqual(telemetry_calls[0], {"temp": 33.0, "ntu": 10.0})
        self.assertEqual(telemetry_calls[1], {"temp": 32.7, "ntu": 10.0})
        # State berkala dikirim setiap siklus
        self.assertEqual(len(state_calls), 2)


class TestArgumenDanEnv(unittest.TestCase):
    def test_bawaan_tanpa_env(self) -> None:
        with patch.dict(os.environ, {}, clear=True):
            args = parse_argumen([])
            self.assertEqual(args.nodes, 3)
            self.assertEqual(args.scenario, "normal")
            self.assertEqual(args.interval, 5)
            self.assertEqual(args.host, "localhost")
            self.assertEqual(args.port, 1883)
            self.assertEqual(args.username, "")
            self.assertEqual(args.password, "")

    def test_baca_dari_env_dan_cli_mengalahkan_env(self) -> None:
        env = {
            "MQTT_HOST": "vm.duckdns.org",
            "MQTT_PORT": "1884",
            "MQTT_USERNAME": "node",
            "MQTT_PASSWORD": "rahasia",
            "SIM_NODES": "4",
            "SIM_SCENARIO": "keruh",
            "SIM_INTERVAL": "10",
        }
        with patch.dict(os.environ, env, clear=True):
            args_env = parse_argumen([])
            self.assertEqual(args_env.host, "vm.duckdns.org")
            self.assertEqual(args_env.port, 1884)
            self.assertEqual(args_env.username, "node")
            self.assertEqual(args_env.password, "rahasia")
            self.assertEqual(args_env.nodes, 4)
            self.assertEqual(args_env.scenario, "keruh")
            self.assertEqual(args_env.interval, 10)

            args_cli = parse_argumen([
                "--host", "mosquitto",
                "--port", "1883",
                "--username", "custom",
                "--password", "pass2",
                "--nodes", "2",
                "--scenario", "panas",
                "--interval", "2",
            ])
            self.assertEqual(args_cli.host, "mosquitto")
            self.assertEqual(args_cli.port, 1883)
            self.assertEqual(args_cli.username, "custom")
            self.assertEqual(args_cli.password, "pass2")
            self.assertEqual(args_cli.nodes, 2)
            self.assertEqual(args_cli.scenario, "panas")
            self.assertEqual(args_cli.interval, 2)


if __name__ == "__main__":
    unittest.main()

