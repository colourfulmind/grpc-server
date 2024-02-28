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
}

type DataGenerator interface {
	GenerateData() (string, float64)
}

func New() *Data {
	return &Data{}
}

func (d *Data) GenerateData() (string, float64) {
	mean := meanMin + rand.Float64()*(meanMax-meanMin)
	std := sdMin + rand.Float64()*(sdMax-sdMin)
	frequency := rand.New(rand.NewSource(time.Now().UnixNano())).NormFloat64()*std + mean
	return uuid.New().String(), frequency
}
