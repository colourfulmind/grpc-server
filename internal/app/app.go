package app

import (
	grpcapp "cmd/grpc/main.go/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, gRPCPort int, StoragePath string, TokenTTL time.Duration) *App {
	gRPCApp := grpcapp.New(log, gRPCPort)
	return &App{
		GRPCSrv: gRPCApp,
	}
}
