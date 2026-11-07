package model

import "github.com/AristoBudiman/smart-aquarium/backend/internal/core"

var kelas = [...]core.Class{core.ClassNormal, core.ClassCaution, core.ClassWarning, core.ClassSevere}

func indeksTerbesar(v []float64) int {
	terbesar := 0
	for i := range v {
		if v[i] > v[terbesar] {
			terbesar = i
		}
	}
	return terbesar
}

// Predict menjalankan ketiga model Random Forest.
func Predict(temp, ntu float64) core.Classes {
	suhu := kelas[indeksTerbesar(ScoreSuhu([]float64{temp}))]
	kekeruhan := kelas[indeksTerbesar(ScoreKekeruhan([]float64{ntu}))]
	overall := kelas[indeksTerbesar(ScoreOverall([]float64{temp, ntu}))]
	return core.Classes{
		Temp:      suhu,
		Turbidity: kekeruhan,
		Overall:   overall,
		Status:    core.Worst(suhu, kekeruhan, overall),
	}
}
