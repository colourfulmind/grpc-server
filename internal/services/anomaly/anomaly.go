package anomaly

import (
	"math"
)

type Anomaly struct {
	Count     float64
	Sum       float64
	Mean      float64
	Deviation float64
	STD       float64
}

func New() *Anomaly {
	return &Anomaly{
		Count:     0.0,
		Sum:       0.0,
		Mean:      0.0,
		Deviation: 0.0,
		STD:       0.0,
	}
}

func (a *Anomaly) UpdateMeanSTD(frequency float64) bool {
	a.Count += 1
	a.Sum += frequency
	a.Mean = a.Sum / a.Count
	a.Deviation += math.Pow(math.Abs(frequency-a.Mean), 2)
	if a.Count > 1 {
		a.STD = math.Sqrt(a.Deviation / (a.Count - 1))
	} else {
		a.STD = 0
	}

	if a.Count == 100 {
		return true
	}

	return false
}

func (a *Anomaly) IsAnomaly(frequency, k float64) bool {
	low := a.Mean - a.STD*k
	high := a.Mean + a.STD*k
	if !(frequency > low && frequency < high) {
		return true
	}
	return false
}
