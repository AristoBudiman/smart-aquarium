package ingest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/decision"
	"github.com/AristoBudiman/smart-aquarium/backend/internal/model"
)

const (
	MinTemp = 0.0
	MaxTemp = 50.0
	MinNTU  = 0.0
)

// Notifier mengirim notifikasi saat alert baru dibuat.
type Notifier interface {
	Notify(ctx context.Context, alert core.Alert)
}

// Service menangani pesan MQTT masuk untuk telemetry, state, dan status node.
type Service struct {
	store       core.Store
	cmd         core.Commander
	broadcaster core.Broadcaster
	notifier    Notifier

	mu          sync.Mutex
	tempInvalid map[string]bool
}

// New membuat instance Service ingest baru.
func New(store core.Store, cmd core.Commander, broadcaster core.Broadcaster, notifier Notifier) *Service {
	return &Service{
		store:       store,
		cmd:         cmd,
		broadcaster: broadcaster,
		notifier:    notifier,
		tempInvalid: make(map[string]bool),
	}
}

// HandleMessage memproses pesan MQTT berdasarkan topiknya:
// aquarium/{node_id}/telemetry, aquarium/{node_id}/state, atau aquarium/{node_id}/status.
func (s *Service) HandleMessage(ctx context.Context, topic string, payload []byte) {
	nodeID, kind, ok := parseTopic(topic)
	if !ok {
		return
	}

	switch kind {
	case "telemetry":
		s.handleTelemetry(ctx, nodeID, payload)
	case "state":
		s.handleState(ctx, nodeID, payload)
	case "status":
		s.handleStatus(ctx, nodeID, payload)
	}
}

func parseTopic(topic string) (nodeID, kind string, ok bool) {
	parts := strings.Split(topic, "/")
	if len(parts) != 3 || parts[0] != "aquarium" || parts[1] == "" || parts[2] == "" {
		return "", "", false
	}
	return parts[1], parts[2], true
}

type rawTelemetry struct {
	Temp *float64 `json:"temp"`
	NTU  *float64 `json:"ntu"`
}

func (s *Service) handleTelemetry(ctx context.Context, nodeID string, payload []byte) {
	// 1. Baca node_id dari topik dan validasi JSON.
	var raw rawTelemetry
	if err := json.Unmarshal(payload, &raw); err != nil || raw.Temp == nil || raw.NTU == nil {
		log.Printf("ingest: payload telemetry %s tidak valid", nodeID)
		return
	}
	temp := *raw.Temp
	ntu := *raw.NTU
	if math.IsNaN(temp) || math.IsInf(temp, 0) || math.IsNaN(ntu) || math.IsInf(ntu, 0) {
		log.Printf("ingest: nilai NaN/Inf pada telemetry %s", nodeID)
		return
	}
	if ntu < MinNTU {
		log.Printf("ingest: nilai ntu negatif (%v) pada %s ditolak", ntu, nodeID)
		return
	}

	// 2. Tolak suhu di luar rentang kontrak (0-50 °C) dan buat alert sekali per perubahan kondisi.
	if temp < MinTemp || temp > MaxTemp {
		if s.markTempInvalid(nodeID, true) {
			s.ensureNode(ctx, nodeID)
			msg := fmt.Sprintf("Suhu %s °C di luar rentang valid (WARNING)", formatDesimalKoma(temp, 1))
			s.saveAndEmitAlert(ctx, nodeID, time.Now().UTC(), core.ClassWarning, msg)
		}
		return
	}
	s.markTempInvalid(nodeID, false)

	// 3. Panggil EnsureNode. Kalau node baru, kirim config manual retained lewat Commander.SetMode.
	if !s.ensureNode(ctx, nodeID) {
		return
	}

	// 4. Ambil Snapshot, karena LastStatus harus dibaca sebelum prediksi baru disimpan.
	snap, err := s.store.Snapshot(ctx, nodeID)
	if err != nil {
		log.Printf("ingest: gagal ambil snapshot %s: %v", nodeID, err)
		return
	}

	now := time.Now().UTC()

	// 5. Simpan telemetry.
	if err := s.store.SaveTelemetry(ctx, core.Telemetry{
		NodeID: nodeID,
		TS:     now,
		Temp:   temp,
		NTU:    ntu,
	}); err != nil {
		log.Printf("ingest: gagal simpan telemetry %s: %v", nodeID, err)
		return
	}

	// 6. Jalankan model.Predict, lalu simpan prediksi.
	classes := model.Predict(temp, ntu)
	if err := s.store.SavePrediction(ctx, nodeID, now, classes); err != nil {
		log.Printf("ingest: gagal simpan prediksi %s: %v", nodeID, err)
		return
	}

	// 7. Kalau mode Auto, jalankan Decide. Pakai prev nol jika HasState false.
	// 8. Kirim cmd dan LogActuators dengan source auto kalau hasilnya berbeda, atau kalau HasState false.
	var autoAction string
	if snap.Mode == core.ModeAuto {
		var prev core.Actuators
		if snap.HasState {
			prev = snap.State.Actuators
		}
		next := decision.Decide(prev, classes, temp)
		if !snap.HasState || next != prev {
			if s.cmd != nil {
				if err := s.cmd.SendActuators(ctx, nodeID, next); err != nil {
					log.Printf("ingest: gagal kirim cmd %s: %v", nodeID, err)
				}
			}
			if err := s.store.LogActuators(ctx, nodeID, next, core.SourceAuto); err != nil {
				log.Printf("ingest: gagal log aktuator %s: %v", nodeID, err)
			}
			autoAction = ringkasanAksiAuto(prev, next, snap.HasState)
		}
	}

	// 9. Kalau status utama naik dan bukan NORMAL, simpan alert dan kirim notifikasi.
	var alertEvent *core.AlertEvent
	if classes.Status != core.ClassNormal && classes.Status.Rank() > snap.LastStatus.Rank() {
		msg := formatPesanStatus(nodeID, temp, ntu, classes, autoAction)
		ev := s.saveAndNotifyAlert(ctx, nodeID, now, classes.Status, msg)
		alertEvent = &ev
	}

	// 10. Kirim TelemetryEvent dan, bila ada, AlertEvent.
	if s.broadcaster != nil {
		s.broadcaster.Publish(core.TelemetryEvent{
			Type:    "telemetry",
			NodeID:  nodeID,
			TS:      now,
			Temp:    temp,
			NTU:     ntu,
			Classes: classes,
		})
		if alertEvent != nil {
			s.broadcaster.Publish(*alertEvent)
		}
	}
}

type rawState struct {
	Heater   *int  `json:"heater"`
	Fan      *int  `json:"fan"`
	Pump     *int  `json:"pump"`
	Failsafe *bool `json:"failsafe"`
}

func (s *Service) handleState(ctx context.Context, nodeID string, payload []byte) {
	var raw rawState
	if err := json.Unmarshal(payload, &raw); err != nil ||
		raw.Heater == nil || raw.Fan == nil || raw.Pump == nil || raw.Failsafe == nil {
		log.Printf("ingest: payload state %s tidak valid", nodeID)
		return
	}
	if (*raw.Heater != 0 && *raw.Heater != 1) ||
		(*raw.Fan != 0 && *raw.Fan != 1) ||
		(*raw.Pump != 0 && *raw.Pump != 1) {
		log.Printf("ingest: nilai aktuator state %s di luar {0,1}", nodeID)
		return
	}

	if !s.ensureNode(ctx, nodeID) {
		return
	}

	snap, err := s.store.Snapshot(ctx, nodeID)
	if err != nil {
		log.Printf("ingest: gagal ambil snapshot %s: %v", nodeID, err)
		return
	}

	st := core.ActuatorState{
		Actuators: core.Actuators{
			Heater: *raw.Heater,
			Fan:    *raw.Fan,
			Pump:   *raw.Pump,
		},
		Failsafe: *raw.Failsafe,
	}

	if err := s.store.SaveState(ctx, nodeID, st); err != nil {
		log.Printf("ingest: gagal simpan state %s: %v", nodeID, err)
		return
	}

	// Saat failsafe berubah menjadi true, catat LogActuators dengan source failsafe.
	prevFailsafe := snap.HasState && snap.State.Failsafe
	if st.Failsafe && !prevFailsafe {
		if err := s.store.LogActuators(ctx, nodeID, st.Actuators, core.SourceFailsafe); err != nil {
			log.Printf("ingest: gagal log failsafe %s: %v", nodeID, err)
		}
	}

	if s.broadcaster != nil {
		s.broadcaster.Publish(core.StateEvent{
			Type:          "state",
			NodeID:        nodeID,
			ActuatorState: st,
		})
	}
}

func (s *Service) handleStatus(ctx context.Context, nodeID string, payload []byte) {
	text := strings.TrimSpace(string(payload))
	var online bool
	switch text {
	case "online":
		online = true
	case "offline":
		online = false
	default:
		log.Printf("ingest: payload status %s tidak dikenal: %q", nodeID, text)
		return
	}

	if !s.ensureNode(ctx, nodeID) {
		return
	}

	if err := s.store.SetOnline(ctx, nodeID, online); err != nil {
		log.Printf("ingest: gagal set online %s: %v", nodeID, err)
		return
	}

	if s.broadcaster != nil {
		s.broadcaster.Publish(core.NodeStatusEvent{
			Type:   "node_status",
			NodeID: nodeID,
			Online: online,
		})
	}
}

func (s *Service) ensureNode(ctx context.Context, nodeID string) bool {
	created, err := s.store.EnsureNode(ctx, nodeID)
	if err != nil {
		log.Printf("ingest: gagal EnsureNode %s: %v", nodeID, err)
		return false
	}
	if created && s.cmd != nil {
		if err := s.cmd.SetMode(ctx, nodeID, core.ModeManual); err != nil {
			log.Printf("ingest: gagal kirim config awal %s: %v", nodeID, err)
		}
	}
	return true
}

// markTempInvalid memperbarui status suhu tidak valid per node dan mengembalikan true
// hanya saat transisi dari valid -> tidak valid.
func (s *Service) markTempInvalid(nodeID string, invalid bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	wasInvalid := s.tempInvalid[nodeID]
	s.tempInvalid[nodeID] = invalid
	return invalid && !wasInvalid
}

func (s *Service) saveAndNotifyAlert(ctx context.Context, nodeID string, ts time.Time, level core.Class, msg string) core.AlertEvent {
	alert := core.Alert{
		NodeID:  nodeID,
		TS:      ts,
		Level:   level,
		Message: msg,
	}
	id, err := s.store.SaveAlert(ctx, alert)
	if err != nil {
		log.Printf("ingest: gagal simpan alert %s: %v", nodeID, err)
	} else {
		alert.ID = id
		if s.notifier != nil {
			s.notifier.Notify(ctx, alert)
		}
	}
	return core.AlertEvent{
		Type:    "alert",
		NodeID:  nodeID,
		Level:   level,
		Message: msg,
	}
}

func (s *Service) saveAndEmitAlert(ctx context.Context, nodeID string, ts time.Time, level core.Class, msg string) {
	ev := s.saveAndNotifyAlert(ctx, nodeID, ts, level, msg)
	if s.broadcaster != nil {
		s.broadcaster.Publish(ev)
	}
}

func formatDesimalKoma(v float64, presisi int) string {
	s := fmt.Sprintf("%.*f", presisi, v)
	return strings.ReplaceAll(s, ".", ",")
}

func formatPesanStatus(nodeID string, temp, ntu float64, c core.Classes, autoAction string) string {
	msg := fmt.Sprintf(
		"Akuarium %s: suhu %s °C (%s), kekeruhan %s NTU (%s).",
		nodeID,
		formatDesimalKoma(temp, 1),
		c.Temp,
		formatDesimalKoma(ntu, 0),
		c.Turbidity,
	)
	if autoAction != "" {
		msg += " " + autoAction
	}
	return msg
}

func ringkasanAksiAuto(prev, next core.Actuators, hasState bool) string {
	var aksi []string
	if (!hasState || prev.Heater != next.Heater) && next.Heater == 1 {
		aksi = append(aksi, "Pemanas dinyalakan otomatis.")
	}
	if (!hasState || prev.Fan != next.Fan) && next.Fan == 1 {
		aksi = append(aksi, "Kipas dinyalakan otomatis.")
	}
	if (!hasState || prev.Pump != next.Pump) && next.Pump == 1 {
		aksi = append(aksi, "Pompa dinyalakan otomatis.")
	}
	return strings.Join(aksi, " ")
}

