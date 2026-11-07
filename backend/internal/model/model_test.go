package model

import (
	"encoding/json"
	"os"
	"testing"
)

type titik struct {
	Temp      float64 `json:"temp"`
	NTU       float64 `json:"ntu"`
	Suhu      int     `json:"suhu"`
	Kekeruhan int     `json:"kekeruhan"`
	Overall   int     `json:"overall"`
}

func TestSamaDenganPython(t *testing.T) {
	data, err := os.ReadFile("testdata.json")
	if err != nil {
		t.Fatal(err)
	}
	var daftar []titik
	if err := json.Unmarshal(data, &daftar); err != nil {
		t.Fatal(err)
	}
	for _, x := range daftar {
		h := Predict(x.Temp, x.NTU)
		if h.Temp != kelas[x.Suhu] || h.Turbidity != kelas[x.Kekeruhan] || h.Overall != kelas[x.Overall] {
			t.Fatalf("berbeda pada suhu=%v ntu=%v: %+v", x.Temp, x.NTU, h)
		}
	}
}
