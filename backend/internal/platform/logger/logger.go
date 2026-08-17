package logger

import (
	"log/slog"
	"os"
)

func New(
	serviceName string,
	environment string,
) *slog.Logger {
	options := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: environment == "local",
	}

	var handler slog.Handler

	if environment == "production" {
		handler = slog.NewJSONHandler(
			os.Stdout,
			options,
		)
	} else {
		handler = slog.NewTextHandler(
			os.Stdout,
			options,
		)
	}

	return slog.New(handler).With(
		"service", serviceName,
		"environment", environment,
	)
}
