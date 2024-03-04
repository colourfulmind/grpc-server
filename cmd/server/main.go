package main

import (
	"grpc/internal/app"
	"grpc/internal/config"
	"grpc/pkg/logger/logsetup"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// initialize config
	cfg := config.MustLoad()

	// initialize logger
	log := logsetup.SetupLogger(cfg.Env)
	log.Info("starting application", slog.Any("config", cfg))

	// initialize app
	application := app.New(log, cfg.GRPC.Port)
	go application.GRPCSrv.MustRun()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sgl := <-stop
	log.Info("stopping application", slog.String("signal", sgl.String()))
	application.GRPCSrv.Stop()
	log.Info("application stopped")
}
