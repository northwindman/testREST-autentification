package main

import (
	"github.com/northwindman/testREST-autentification/internal/config"
	slogpretty "github.com/northwindman/testREST-autentification/internal/lib/logger"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting url-shortener",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)

	// TODO: init storage

	// TODO: init app

}

func setupLogger(env string) *slog.Logger {

	var log = &slog.Logger{}

	switch env {
	case envLocal:
		log = setupPrettySlog() // use only in local mode
	case envDev:
		log = slog.New(slog.NewTextHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo}),
		)
	case envProd:
		log = slog.New(slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{
				Level: slog.LevelInfo}),
		)
	}

	return log

}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
