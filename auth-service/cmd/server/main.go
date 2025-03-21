package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"
	"auth-service/pkg/jwt"
	"fmt"
	"log/slog"
	"os"
)

func main() {
	cfg := config.LoadConfig()

	DBDSN := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	log.Info("Starting server")

	jwt.Init(cfg.JWTSecret)

	app := app.New(log, DBDSN, cfg.TokenTTL)

	app.Run(cfg.AppHost, cfg.AppPort)
}
