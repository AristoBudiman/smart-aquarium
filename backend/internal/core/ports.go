// Package core berisi tipe dan interface bersama Peran 3 dan Peran 4.
// Perubahan file ini sama dengan perubahan kontrak.
package core

import (
	"context"
	"time"
)

type Class string

const (
	ClassNormal  Class = "NORMAL"
	ClassCaution Class = "CAUTION"
	ClassWarning Class = "WARNING"
	ClassSevere  Class = "SEVERE"
)

func (c Class) Rank() int {
	switch c {
	case ClassSevere:
		return 3
	case ClassWarning:
		return 2
	case ClassCaution:
		return 1
	default:
		return 0
	}
}

// Worst mengembalikan kelas paling parah.
func Worst(classes ...Class) Class {
	hasil := ClassNormal
	for _, c := range classes {
		if c.Rank() > hasil.Rank() {
			hasil = c
		}
	}
	return hasil
}

type Mode string

const (
	ModeAuto   Mode = "auto"
	ModeManual Mode = "manual"
)

type Source string

const (
	SourceAuto     Source = "auto"
	SourceManual   Source = "manual"
	SourceFailsafe Source = "failsafe"
)

type Actuators struct {
	Heater int `json:"heater"`
	Fan    int `json:"fan"`
	Pump   int `json:"pump"`
}

type ActuatorState struct {
	Actuators
	Failsafe bool `json:"failsafe"`
}

type Telemetry struct {
	NodeID string    `json:"node_id"`
	TS     time.Time `json:"ts"`
	Temp   float64   `json:"temp"`
	NTU    float64   `json:"ntu"`
}

type Classes struct {
	Temp      Class `json:"class_temp"`
	Turbidity Class `json:"class_turbidity"`
	Overall   Class `json:"class_overall"`
	Status    Class `json:"status"`
}

// NodeSnapshot adalah kondisi terakhir yang dibutuhkan logika keputusan.
type NodeSnapshot struct {
	NodeID     string
	Mode       Mode
	Online     bool
	HasState   bool
	State      ActuatorState
	LastStatus Class
}

type Alert struct {
	ID      int64     `json:"id"`
	NodeID  string    `json:"node_id"`
	TS      time.Time `json:"ts"`
	Level   Class     `json:"level"`
	Message string    `json:"message"`
}

// Store diimplementasikan Peran 4 di internal/db dan dipakai Peran 3.
type Store interface {
	EnsureNode(ctx context.Context, nodeID string) (created bool, err error)
	Snapshot(ctx context.Context, nodeID string) (NodeSnapshot, error)
	SaveTelemetry(ctx context.Context, t Telemetry) error
	SavePrediction(ctx context.Context, nodeID string, ts time.Time, c Classes) error
	SaveState(ctx context.Context, nodeID string, s ActuatorState) error
	SetOnline(ctx context.Context, nodeID string, online bool) error
	LogActuators(ctx context.Context, nodeID string, a Actuators, source Source) error
	SaveAlert(ctx context.Context, a Alert) (id int64, err error)
	MarkAlertSent(ctx context.Context, id int64) error
	TelegramChatIDs(ctx context.Context, nodeID string) ([]string, error)
}

// Commander diimplementasikan Peran 3 di internal/mqtt dan dipakai Peran 4.
type Commander interface {
	SetMode(ctx context.Context, nodeID string, mode Mode) error
	SendActuators(ctx context.Context, nodeID string, a Actuators) error
}

// Broadcaster diimplementasikan Peran 4 di internal/ws dan dipakai Peran 3.
type Broadcaster interface {
	Publish(event any)
}

// Event WebSocket. Bentuk JSON-nya ditetapkan di docs/kontrak.md.
type TelemetryEvent struct {
	Type   string    `json:"type"` // "telemetry"
	NodeID string    `json:"node_id"`
	TS     time.Time `json:"ts"`
	Temp   float64   `json:"temp"`
	NTU    float64   `json:"ntu"`
	Classes
}

type StateEvent struct {
	Type   string `json:"type"` // "state"
	NodeID string `json:"node_id"`
	ActuatorState
}

type NodeStatusEvent struct {
	Type   string `json:"type"` // "node_status"
	NodeID string `json:"node_id"`
	Online bool   `json:"online"`
}

type ModeEvent struct {
	Type   string `json:"type"` // "mode"
	NodeID string `json:"node_id"`
	Mode   Mode   `json:"mode"`
}

type AlertEvent struct {
	Type    string `json:"type"` // "alert"
	NodeID  string `json:"node_id"`
	Level   Class  `json:"level"`
	Message string `json:"message"`
}
