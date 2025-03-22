package main

import (
	"auth-service/internal/app"
	"auth-service/internal/config"
	"auth-service/pkg/jwt"
	"auth-service/pkg/oauth2/github"
	"auth-service/pkg/oauth2/google"
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
	google := google.NewGoogleService(cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.GoogleRedirectURL, log)

	log.Info("Starting server")

	fmt.Println(cfg.GithubClientID)

	jwt.Init(cfg.JWTSecret)

	app := app.New(log, DBDSN, cfg.TokenTTL, githhub, google)

	app.Run(cfg.AppHost, cfg.AppPort)
}
