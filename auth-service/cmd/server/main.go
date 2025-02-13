package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"
	"auth-service/pkg/jwt"
	"log/slog"
	"os"
)

func main() {
	cfg := config.LoadConfig()

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	log.Info("Starting server")

	jwt.Init(cfg.JWTSecret)

	app := app.New(log, cfg.DBDSN, cfg.TokenTTL)

	app.Run(cfg.AppHost, cfg.AppPort)
}
