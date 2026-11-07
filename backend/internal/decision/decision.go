package decision

import "github.com/AristoBudiman/smart-aquarium/backend/internal/core"

// SuhuAcuan adalah batas suhu (°C) untuk menentukan pemanas vs kipas saat Warning/Severe.
const SuhuAcuan = 27.48

// Decide menentukan kondisi aktuator baru berdasarkan status sebelumnya, hasil klasifikasi, dan suhu.
func Decide(prev core.Actuators, c core.Classes, temp float64) core.Actuators {
	out := prev

	switch c.Temp {
	case core.ClassNormal:
		out.Heater = 0
		out.Fan = 0
	case core.ClassCaution:
		// tetap mempertahankan nilai prev.Heater dan prev.Fan
	case core.ClassWarning, core.ClassSevere:
		if temp < SuhuAcuan {
			out.Heater = 1
			out.Fan = 0
		} else {
			out.Heater = 0
			out.Fan = 1
		}
	}

	switch c.Turbidity {
	case core.ClassNormal:
		out.Pump = 0
	case core.ClassCaution:
		// tetap mempertahankan nilai prev.Pump
	case core.ClassWarning, core.ClassSevere:
		out.Pump = 1
	}

	return out
}

