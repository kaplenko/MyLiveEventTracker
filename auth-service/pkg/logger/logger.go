package logger

import (
	"os"

	"log/slog"
)

var lg *slog.Logger

func InitLogger() {
	lg = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func Info(msg string, args ...any) {
	lg.Info(msg, args...)
}

func Error(msg string, args ...any) {
	lg.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	lg.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	lg.Warn(msg, args...)
}
