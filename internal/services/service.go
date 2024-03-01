package services

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const (
	meanMin float64 = -10
	meanMax float64 = 10
	sdMin   float64 = 0.3
	sdMax   float64 = 1.5
)

type Data struct {
	SessionID string
	Mean      float64
	STD       float64
}

func New() *Data {
	mean := meanMin + rand.Float64()*(meanMax-meanMin)
	std := sdMin + rand.Float64()*(sdMax-sdMin)
	return &Data{
		SessionID: uuid.New().String(),
		Mean:      mean,
		STD:       std,
	}
}

func (d *Data) GenerateFrequency() float64 {
	return rand.New(rand.NewSource(time.Now().UnixNano())).NormFloat64()*d.STD + d.Mean
}
