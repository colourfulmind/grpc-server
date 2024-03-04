package main

import (
	"context"
	grpcclient "grpc/internal/clients/receiver/grpc"
	"grpc/internal/config"
	"grpc/pkg/logger/logsetup"
	"grpc/pkg/logger/sl"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	const op = "grpc.New"

	// initialize config
	cfg := config.MustLoad()

	// initialize logger
	log := logsetup.SetupLogger(cfg.Env)
	log.Info("starting client", slog.Any("config", cfg))

	// connect to server
	cc, err := grpcclient.NewConnection(
		context.Background(),
		log,
		cfg.Clients.Receiver.Address,
		cfg.Clients.Receiver.RetriesCount,
		cfg.Clients.Receiver.Timeout,
	)
	if err != nil {
		log.Error("failed to connect to server", op, sl.Err(err))
		os.Exit(1)
	}
	defer cc.Close()

	// initialize client
	client, err := grpcclient.New(cc, log, cfg)
	if err != nil {
		os.Exit(1)
	}

	// start client
	go client.Receiver()

	// shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sgl := <-stop
	log.Info("stopping client", slog.String("signal", sgl.String()))
	log.Info("client stopped")

}
