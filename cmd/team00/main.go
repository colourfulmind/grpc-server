package main

import (
	"grpc/internal/app"
	"grpc/internal/config"
	"grpc/internal/lib/logger/handlers/slogpretty"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

func main() {
	// initialize config
	cfg := config.MustLoad()

	// initialize logger
	log := SetupLogger(cfg.Env)
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

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case EnvLocal:
		//log = slog.New(
		//	slog.NewTextHandler(
		//		os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug},
		//	),
		//)
		log = SetupPrettySlog()
	case EnvDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case EnvProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}

	return log
}

func SetupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
