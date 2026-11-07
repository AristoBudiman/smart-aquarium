package api

import (
	"context"
	"log"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
)

// SEMENTARA: implementasi core.Commander sementara sebelum modul MQTT Peran 3 terhubung.
// Implementasi ini hanya mencatat perintah ke log dan mengembalikan nil.
type TempCommander struct{}

var _ core.Commander = (*TempCommander)(nil)

// NewTempCommander membuat instance TempCommander baru.
func NewTempCommander() *TempCommander {
	return &TempCommander{}
}

// SetMode mencatat perubahan mode operasi node akuarium.
func (c *TempCommander) SetMode(_ context.Context, nodeID string, mode core.Mode) error {
	log.Printf("[SEMENTARA] Commander: SetMode node=%s mode=%s", nodeID, mode)
	return nil
}

// SendActuators mencatat perintah pengalihan relay aktuator akuarium.
func (c *TempCommander) SendActuators(_ context.Context, nodeID string, a core.Actuators) error {
	log.Printf("[SEMENTARA] Commander: SendActuators node=%s heater=%d fan=%d pump=%d", nodeID, a.Heater, a.Fan, a.Pump)
	return nil
}
