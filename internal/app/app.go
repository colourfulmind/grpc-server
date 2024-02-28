package app

import (
	grpcapp "cmd/grpc/main.go/internal/app/grpc"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(log *slog.Logger, gRPCPort int) *App {
	gRPCApp := grpcapp.New(log, gRPCPort)
	return &App{
		GRPCSrv: gRPCApp,
	}
}
