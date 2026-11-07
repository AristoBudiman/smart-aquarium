package decision

import (
	"testing"

	"github.com/AristoBudiman/smart-aquarium/backend/internal/core"
)

func TestDecide(t *testing.T) {
	kasus := []struct {
		nama string
		prev core.Actuators
		c    core.Classes
		temp float64
		want core.Actuators
	}{
		{
			nama: "Suhu Normal mematikan pemanas dan kipas, Kekeruhan Normal mematikan pompa",
			prev: core.Actuators{Heater: 1, Fan: 1, Pump: 1},
			c:    core.Classes{Temp: core.ClassNormal, Turbidity: core.ClassNormal},
			temp: 27.0,
			want: core.Actuators{Heater: 0, Fan: 0, Pump: 0},
		},
		{
			nama: "Suhu Caution mempertahankan pemanas dan kipas (ON/OFF), Kekeruhan Caution mempertahankan pompa (ON)",
			prev: core.Actuators{Heater: 1, Fan: 0, Pump: 1},
			c:    core.Classes{Temp: core.ClassCaution, Turbidity: core.ClassCaution},
			temp: 25.0,
			want: core.Actuators{Heater: 1, Fan: 0, Pump: 1},
		},
		{
			nama: "Suhu Caution mempertahankan pemanas dan kipas (OFF/ON), Kekeruhan Caution mempertahankan pompa (OFF)",
			prev: core.Actuators{Heater: 0, Fan: 1, Pump: 0},
			c:    core.Classes{Temp: core.ClassCaution, Turbidity: core.ClassCaution},
			temp: 30.0,
			want: core.Actuators{Heater: 0, Fan: 1, Pump: 0},
		},
		{
			nama: "Suhu Warning di bawah 27.48 menyalakan pemanas dan mematikan kipas",
			prev: core.Actuators{Heater: 0, Fan: 1, Pump: 0},
			c:    core.Classes{Temp: core.ClassWarning, Turbidity: core.ClassNormal},
			temp: 22.0,
			want: core.Actuators{Heater: 1, Fan: 0, Pump: 0},
		},
		{
			nama: "Suhu Severe di bawah 27.48 menyalakan pemanas dan mematikan kipas",
			prev: core.Actuators{Heater: 0, Fan: 0, Pump: 0},
			c:    core.Classes{Temp: core.ClassSevere, Turbidity: core.ClassNormal},
			temp: 18.5,
			want: core.Actuators{Heater: 1, Fan: 0, Pump: 0},
		},
		{
			nama: "Suhu Warning tepat 27.48 mematikan pemanas dan menyalakan kipas",
			prev: core.Actuators{Heater: 1, Fan: 0, Pump: 0},
			c:    core.Classes{Temp: core.ClassWarning, Turbidity: core.ClassNormal},
			temp: 27.48,
			want: core.Actuators{Heater: 0, Fan: 1, Pump: 0},
		},
		{
			nama: "Suhu Severe tepat 27.48 mematikan pemanas dan menyalakan kipas",
			prev: core.Actuators{Heater: 1, Fan: 0, Pump: 0},
			c:    core.Classes{Temp: core.ClassSevere, Turbidity: core.ClassNormal},
			temp: 27.48,
			want: core.Actuators{Heater: 0, Fan: 1, Pump: 0},
		},
		{
			nama: "Suhu Warning di atas 27.48 mematikan pemanas dan menyalakan kipas",
			prev: core.Actuators{Heater: 1, Fan: 0, Pump: 0},
			c:    core.Classes{Temp: core.ClassWarning, Turbidity: core.ClassNormal},
			temp: 33.0,
			want: core.Actuators{Heater: 0, Fan: 1, Pump: 0},
		},
		{
			nama: "Suhu Severe di atas 27.48 mematikan pemanas dan menyalakan kipas",
			prev: core.Actuators{Heater: 0, Fan: 0, Pump: 0},
			c:    core.Classes{Temp: core.ClassSevere, Turbidity: core.ClassNormal},
			temp: 36.0,
			want: core.Actuators{Heater: 0, Fan: 1, Pump: 0},
		},
		{
			nama: "Kekeruhan Warning menyalakan pompa",
			prev: core.Actuators{Heater: 0, Fan: 0, Pump: 0},
			c:    core.Classes{Temp: core.ClassNormal, Turbidity: core.ClassWarning},
			temp: 27.0,
			want: core.Actuators{Heater: 0, Fan: 0, Pump: 1},
		},
		{
			nama: "Kekeruhan Severe menyalakan pompa",
			prev: core.Actuators{Heater: 0, Fan: 0, Pump: 0},
			c:    core.Classes{Temp: core.ClassNormal, Turbidity: core.ClassSevere},
			temp: 27.0,
			want: core.Actuators{Heater: 0, Fan: 0, Pump: 1},
		},
	}

	for _, tc := range kasus {
		t.Run(tc.nama, func(t *testing.T) {
			got := Decide(tc.prev, tc.c, tc.temp)
			if got != tc.want {
				t.Fatalf("Decide(%+v, %+v, %v) = %+v, ingin %+v", tc.prev, tc.c, tc.temp, got, tc.want)
			}
		})
	}
}

