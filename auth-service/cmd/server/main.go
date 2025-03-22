package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"
	"auth-service/pkg/jwt"
	"auth-service/pkg/oauth2/github"
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

	githhub := github.NewGithubService(cfg.GithubClientID, cfg.GithubClientSecret, cfg.GithubRedirectURL, log)

	log.Info("Starting server")

	fmt.Println(cfg.GithubClientID)
	log.Info("GithubClientSecret: %V", cfg.GithubClientSecret)
	log.Info("GithubRedirectURL: %V", cfg.GithubRedirectURL)

	jwt.Init(cfg.JWTSecret)

	app := app.New(log, DBDSN, cfg.TokenTTL, githhub)

	app.Run(cfg.AppHost, cfg.AppPort)
}
