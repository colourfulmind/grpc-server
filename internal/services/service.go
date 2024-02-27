package services

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
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
	log *slog.Logger
}

type DataGenerator interface {
	GenerateData(ctx context.Context) (string, float64)
}

func New(log *slog.Logger) *Data {
	return &Data{
		log: log,
	}
}

func (d *Data) GenerateData(ctx context.Context) (string, float64) {
	const op = "service.GenerateData"
	log := d.log.With(slog.String("op", op))
	log.Info("start generating data for response")

	mean := meanMin + rand.Float64()*(meanMax-meanMin)
	std := sdMin + rand.Float64()*(sdMax-sdMin)
	frequency := rand.New(rand.NewSource(time.Now().UnixNano())).NormFloat64()*std + mean

	log.Info("data generated successfully")
	return uuid.New().String(), frequency
}