package app

import (
	httpApp "auth-service/internal/infrastructure/http"
	"auth-service/internal/infrastructure/storage"
	"auth-service/internal/usecase"
	"auth-service/pkg/oauth2/github"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	httpApp *httpApp.Handler
}

func New(log *slog.Logger, connStr string, tokenTTL time.Duration, gh *github.Service) *App {
	strg, err := storage.New(connStr)
	if err != nil {
		panic(err)
	}
	authService := usecase.New(strg, strg, strg, log, tokenTTL)

	r := mux.NewRouter()

	httpApp := httpApp.New(authService, gh, r, log)

	return &App{
		httpApp: httpApp,
	}
}

func (app *App) Run(host, port string) {
	app.httpApp.SetupRoutes()
	if err := http.ListenAndServe(host+":"+port, app.httpApp.Router()); err != nil {
		panic(err)
	}
}
