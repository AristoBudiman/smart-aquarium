package ingest

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
)

// SEMENTARA: implementasi Store palsu untuk unit test Peran 3.
type fakeStore struct {
	mu          sync.Mutex
	nodes       map[string]*core.NodeSnapshot
	telemetries []core.Telemetry
	predictions []core.Classes
	actLogs     []actuatorLogEntry
	alerts      []core.Alert

	ensureErr    error
	saveAlertErr error
	onEnsureNode func(nodeID string)
}

type actuatorLogEntry struct {
	nodeID string
	act    core.Actuators
	source core.Source
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		nodes: make(map[string]*core.NodeSnapshot),
	}
}

func (f *fakeStore) EnsureNode(_ context.Context, nodeID string) (bool, error) {
	f.mu.Lock()
	hook := f.onEnsureNode
	err := f.ensureErr
	f.mu.Unlock()

	if hook != nil {
		hook(nodeID)
	}
	if err != nil {
		return false, err
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	if _, ok := f.nodes[nodeID]; ok {
		return false, nil
	}
	f.nodes[nodeID] = &core.NodeSnapshot{
		NodeID:     nodeID,
		Mode:       core.ModeManual,
		LastStatus: core.ClassNormal,
	}
	return true, nil
}

func (f *fakeStore) Snapshot(_ context.Context, nodeID string) (core.NodeSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	n, ok := f.nodes[nodeID]
	if !ok {
		return core.NodeSnapshot{NodeID: nodeID, Mode: core.ModeManual, LastStatus: core.ClassNormal}, nil
	}
	return *n, nil
}

func (f *fakeStore) SaveTelemetry(_ context.Context, t core.Telemetry) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.telemetries = append(f.telemetries, t)
	return nil
}

func (f *fakeStore) SavePrediction(_ context.Context, nodeID string, _ time.Time, c core.Classes) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.predictions = append(f.predictions, c)
	if n, ok := f.nodes[nodeID]; ok {
		n.LastStatus = c.Status
	}
	return nil
}

func (f *fakeStore) SaveState(_ context.Context, nodeID string, s core.ActuatorState) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if n, ok := f.nodes[nodeID]; ok {
		n.HasState = true
		n.State = s
	}
	return nil
}

func (f *fakeStore) SetOnline(_ context.Context, nodeID string, online bool) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if n, ok := f.nodes[nodeID]; ok {
		n.Online = online
	}
	return nil
}

func (f *fakeStore) LogActuators(_ context.Context, nodeID string, a core.Actuators, source core.Source) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.actLogs = append(f.actLogs, actuatorLogEntry{nodeID: nodeID, act: a, source: source})
	return nil
}

func (f *fakeStore) SaveAlert(_ context.Context, a core.Alert) (int64, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.saveAlertErr != nil {
		return 0, f.saveAlertErr
	}
	a.ID = int64(len(f.alerts) + 1)
	f.alerts = append(f.alerts, a)
	return a.ID, nil
}

func (f *fakeStore) MarkAlertSent(_ context.Context, _ int64) error { return nil }
func (f *fakeStore) TelegramChatIDs(_ context.Context) ([]string, error) {
	return nil, nil
}

// SEMENTARA: implementasi Commander, Broadcaster, dan Notifier palsu untuk unit test Peran 3.
type fakeCommander struct {
	mu    sync.Mutex
	modes []core.Mode
	cmds  []core.Actuators
}

func (c *fakeCommander) SetMode(_ context.Context, _ string, mode core.Mode) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.modes = append(c.modes, mode)
	return nil
}

func (c *fakeCommander) SendActuators(_ context.Context, _ string, a core.Actuators) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cmds = append(c.cmds, a)
	return nil
}

type fakeBroadcaster struct {
	mu     sync.Mutex
	events []any
}

func (b *fakeBroadcaster) Publish(event any) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.events = append(b.events, event)
}

type fakeNotifier struct {
	mu     sync.Mutex
	alerts []core.Alert
}

func (n *fakeNotifier) Notify(_ context.Context, alert core.Alert) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.alerts = append(n.alerts, alert)
}

func TestTelemetryNewNodeAndStatusEscalation(t *testing.T) {
	store := newFakeStore()
	cmd := &fakeCommander{}
	bc := &fakeBroadcaster{}
	notif := &fakeNotifier{}
	svc := New(store, cmd, bc, notif)
	ctx := context.Background()

	// Kirim telemetry dengan suhu 33 °C (WARNING) pada node baru aq-test01
	svc.HandleMessage(ctx, "aquarium/aq-test01/telemetry", []byte(`{"temp":33,"ntu":10}`))

	// 1. Node baru harus mengirim config manual retained
	if len(cmd.modes) != 1 || cmd.modes[0] != core.ModeManual {
		t.Fatalf("SetMode = %+v, ingin [manual]", cmd.modes)
	}

	// 2. Telemetry dan prediksi harus tersimpan dengan status WARNING
	if len(store.telemetries) != 1 {
		t.Fatalf("jumlah telemetry = %d, ingin 1", len(store.telemetries))
	}
	if len(store.predictions) != 1 || store.predictions[0].Status != core.ClassWarning {
		t.Fatalf("prediksi = %+v, ingin Status WARNING", store.predictions)
	}

	// 3. Karena status naik dari NORMAL -> WARNING, alert dibuat dan notifikasi dikirim
	if len(store.alerts) != 1 || store.alerts[0].Level != core.ClassWarning {
		t.Fatalf("alerts = %+v, ingin 1 alert WARNING", store.alerts)
	}
	if len(notif.alerts) != 1 {
		t.Fatalf("notifier alerts = %d, ingin 1", len(notif.alerts))
	}

	// 4. TelemetryEvent dan AlertEvent harus dikirim ke Broadcaster
	if len(bc.events) != 2 {
		t.Fatalf("jumlah event = %d, ingin 2 (TelemetryEvent + AlertEvent)", len(bc.events))
	}

	// 5. Kirim telemetry kedua dengan status yang sama (WARNING): tidak boleh membuat alert duplikat
	svc.HandleMessage(ctx, "aquarium/aq-test01/telemetry", []byte(`{"temp":33,"ntu":10}`))
	if len(store.alerts) != 1 {
		t.Fatalf("alert duplikat dibuat saat status tetap WARNING: %+v", store.alerts)
	}
}

func TestTelemetryAutoModeDecision(t *testing.T) {
	store := newFakeStore()
	store.nodes["aq-auto01"] = &core.NodeSnapshot{
		NodeID:     "aq-auto01",
		Mode:       core.ModeAuto,
		HasState:   false,
		LastStatus: core.ClassNormal,
	}
	cmd := &fakeCommander{}
	bc := &fakeBroadcaster{}
	svc := New(store, cmd, bc, nil)
	ctx := context.Background()

	// HasState false: harus mengirim cmd meskipun hasil Decide sama dengan prev nol, atau saat berbeda
	svc.HandleMessage(ctx, "aquarium/aq-auto01/telemetry", []byte(`{"temp":27.3,"ntu":10}`))
	if len(cmd.cmds) != 1 || len(store.actLogs) != 1 || store.actLogs[0].source != core.SourceAuto {
		t.Fatalf("cmd/log tidak terpanggil saat HasState false: cmds=%+v logs=%+v", cmd.cmds, store.actLogs)
	}

	// Set HasState true dengan state sesuai hasil terakhir (0,0,0), lalu kirim kondisi Normal lagi -> tidak kirim cmd
	store.nodes["aq-auto01"].HasState = true
	store.nodes["aq-auto01"].State = core.ActuatorState{Actuators: core.Actuators{Heater: 0, Fan: 0, Pump: 0}}
	svc.HandleMessage(ctx, "aquarium/aq-auto01/telemetry", []byte(`{"temp":27.3,"ntu":10}`))
	if len(cmd.cmds) != 1 {
		t.Fatalf("cmd seharusnya tidak dikirim jika state sudah sama: %+v", cmd.cmds)
	}

	// Kirim suhu 33 °C (Warning >= 27.48) -> Fan harus menyala (0,1,0)
	svc.HandleMessage(ctx, "aquarium/aq-auto01/telemetry", []byte(`{"temp":33,"ntu":10}`))
	if len(cmd.cmds) != 2 || cmd.cmds[1] != (core.Actuators{Heater: 0, Fan: 1, Pump: 0}) {
		t.Fatalf("cmd ke-2 = %+v, ingin {Heater:0 Fan:1 Pump:0}", cmd.cmds)
	}
}

func TestTelemetryOutOfRangeTempAlertOncePerCondition(t *testing.T) {
	store := newFakeStore()
	cmd := &fakeCommander{}
	bc := &fakeBroadcaster{}
	svc := New(store, cmd, bc, nil)
	ctx := context.Background()

	// Suhu 60 °C (di luar 0-50 °C): ditolak dan buat alert WARNING sekali
	svc.HandleMessage(ctx, "aquarium/aq-temp01/telemetry", []byte(`{"temp":60,"ntu":10}`))
	if len(store.telemetries) != 0 {
		t.Fatal("telemetry di luar rentang seharusnya ditolak")
	}
	if len(store.alerts) != 1 || store.alerts[0].Level != core.ClassWarning {
		t.Fatalf("alert tidak dibuat untuk suhu di luar rentang: %+v", store.alerts)
	}

	// Suhu -5 °C berturut-turut: masih dalam kondisi tidak valid, jangan buat alert baru
	svc.HandleMessage(ctx, "aquarium/aq-temp01/telemetry", []byte(`{"temp":-5,"ntu":10}`))
	if len(store.alerts) != 1 {
		t.Fatalf("alert seharusnya hanya dibuat sekali per perubahan kondisi: %+v", store.alerts)
	}

	// Suhu kembali normal (27 °C): kondisi tidak valid direset
	svc.HandleMessage(ctx, "aquarium/aq-temp01/telemetry", []byte(`{"temp":27,"ntu":10}`))
	if len(store.telemetries) != 1 {
		t.Fatal("telemetry valid seharusnya disimpan")
	}

	// Suhu kembali di luar rentang (55 °C): buat alert lagi karena kondisi berubah lagi
	svc.HandleMessage(ctx, "aquarium/aq-temp01/telemetry", []byte(`{"temp":55,"ntu":10}`))
	if len(store.alerts) != 2 {
		t.Fatalf("alert ke-2 seharusnya dibuat setelah kondisi kembali tidak valid: %+v", store.alerts)
	}
}

func TestOutOfRangeTempEnsureNodeAndSaveAlertFailureRetries(t *testing.T) {
	store := newFakeStore()
	svc := New(store, nil, nil, nil)
	ctx := context.Background()

	// 1. EnsureNode gagal -> tidak boleh memanggil SaveAlert dan pesan berikutnya harus tetap mencoba
	store.ensureErr = errors.New("db mati")
	svc.HandleMessage(ctx, "aquarium/aq-err01/telemetry", []byte(`{"temp":60,"ntu":10}`))
	if len(store.alerts) != 0 {
		t.Fatalf("SaveAlert seharusnya tidak dipanggil saat EnsureNode gagal: %+v", store.alerts)
	}

	// 2. EnsureNode pulih, tapi SaveAlert sempat gagal -> pesan berikutnya tetap mencoba membuat alert
	store.ensureErr = nil
	store.saveAlertErr = errors.New("insert alert gagal")
	svc.HandleMessage(ctx, "aquarium/aq-err01/telemetry", []byte(`{"temp":60,"ntu":10}`))
	if len(store.alerts) != 0 {
		t.Fatalf("alerts seharusnya masih kosong saat SaveAlert gagal: %+v", store.alerts)
	}

	// 3. SaveAlert pulih -> pesan berikutnya berhasil membuat alert
	store.saveAlertErr = nil
	svc.HandleMessage(ctx, "aquarium/aq-err01/telemetry", []byte(`{"temp":60,"ntu":10}`))
	if len(store.alerts) != 1 {
		t.Fatalf("alert seharusnya dibuat pada percobaan ulang: %+v", store.alerts)
	}
}

func TestPerNodeSerializationAndCrossNodeParallelism(t *testing.T) {
	store := newFakeStore()
	svc := New(store, nil, nil, nil)
	ctx := context.Background()

	// Bagian A: Dua pesan dari node yang sama tidak diproses bersamaan
	masukPertama := make(chan struct{})
	lepasPertama := make(chan struct{})
	var aktifNodeSama int32
	var tumpangTindih int32

	store.onEnsureNode = func(nodeID string) {
		if nodeID != "aq-sama" {
			return
		}
		if atomic.AddInt32(&aktifNodeSama, 1) > 1 {
			atomic.StoreInt32(&tumpangTindih, 1)
		}
		select {
		case masukPertama <- struct{}{}:
			<-lepasPertama
		default:
		}
		atomic.AddInt32(&aktifNodeSama, -1)
	}

	var wgSama sync.WaitGroup
	wgSama.Add(2)
	go func() {
		defer wgSama.Done()
		svc.HandleMessage(ctx, "aquarium/aq-sama/state", []byte(`{"heater":0,"fan":1,"pump":0,"failsafe":false}`))
	}()
	<-masukPertama

	selesaiKedua := make(chan struct{})
	go func() {
		defer wgSama.Done()
		svc.HandleMessage(ctx, "aquarium/aq-sama/telemetry", []byte(`{"temp":27,"ntu":10}`))
		close(selesaiKedua)
	}()

	select {
	case <-selesaiKedua:
		t.Fatal("pesan kedua dari node yang sama selesai sebelum pesan pertama dilepas")
	case <-time.After(30 * time.Millisecond):
	}

	close(lepasPertama)
	wgSama.Wait()
	if atomic.LoadInt32(&tumpangTindih) != 0 {
		t.Fatal("dua pesan dari node yang sama diproses bersamaan")
	}

	// Bagian B: Dua pesan dari node berbeda tetap bisa diproses paralel
	tibaNodeBerbeda := make(chan string, 2)
	lepasNodeBerbeda := make(chan struct{})
	store.onEnsureNode = func(nodeID string) {
		if nodeID == "aq-par1" || nodeID == "aq-par2" {
			tibaNodeBerbeda <- nodeID
			<-lepasNodeBerbeda
		}
	}

	var wgBeda sync.WaitGroup
	wgBeda.Add(2)
	go func() {
		defer wgBeda.Done()
		svc.HandleMessage(ctx, "aquarium/aq-par1/status", []byte("online"))
	}()
	go func() {
		defer wgBeda.Done()
		svc.HandleMessage(ctx, "aquarium/aq-par2/status", []byte("online"))
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-tibaNodeBerbeda:
		case <-time.After(500 * time.Millisecond):
			t.Fatal("node berbeda tidak berjalan secara paralel")
		}
	}
	close(lepasNodeBerbeda)
	wgBeda.Wait()
}

func TestStateAndStatusTopics(t *testing.T) {
	store := newFakeStore()
	bc := &fakeBroadcaster{}
	svc := New(store, nil, bc, nil)
	ctx := context.Background()

	// State normal (failsafe false)
	svc.HandleMessage(ctx, "aquarium/aq-st01/state", []byte(`{"heater":0,"fan":1,"pump":0,"failsafe":false}`))
	if len(store.actLogs) != 0 {
		t.Fatalf("tidak boleh log failsafe saat failsafe false: %+v", store.actLogs)
	}

	// State berubah menjadi failsafe true -> harus LogActuators dengan source failsafe
	svc.HandleMessage(ctx, "aquarium/aq-st01/state", []byte(`{"heater":0,"fan":0,"pump":1,"failsafe":true}`))
	if len(store.actLogs) != 1 || store.actLogs[0].source != core.SourceFailsafe {
		t.Fatalf("log failsafe tidak tercatat: %+v", store.actLogs)
	}

	// State tetap failsafe true -> tidak log ulang
	svc.HandleMessage(ctx, "aquarium/aq-st01/state", []byte(`{"heater":0,"fan":0,"pump":1,"failsafe":true}`))
	if len(store.actLogs) != 1 {
		t.Fatalf("log failsafe tercatat ganda: %+v", store.actLogs)
	}

	// Status online / offline
	svc.HandleMessage(ctx, "aquarium/aq-st01/status", []byte("online"))
	snap, _ := store.Snapshot(ctx, "aq-st01")
	if !snap.Online {
		t.Fatal("node seharusnya online")
	}
	svc.HandleMessage(ctx, "aquarium/aq-st01/status", []byte("offline"))
	snap, _ = store.Snapshot(ctx, "aq-st01")
	if snap.Online {
		t.Fatal("node seharusnya offline")
	}
}
