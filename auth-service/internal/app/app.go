package app

import (
	httpApp "auth-service/internal/infrastructure/http"
	"auth-service/internal/infrastructure/storage"
	"auth-service/internal/usecase"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	httpApp *httpApp.Handler
}

func New(log *slog.Logger, connStr string, tokenTTL time.Duration) *App {
	storage, err := storage.New(connStr)
	if err != nil {
		panic(err)
	}
	authService := usecase.New(storage, storage, log, tokenTTL)

	r := mux.NewRouter()

	httpApp := httpApp.New(authService, r, log)

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
