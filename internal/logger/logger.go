// Package logger provides structured logging with zerolog and slog.
package logger

import (
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

// New creates a zerolog logger with the given level and writer.
func New(level string, output io.Writer) zerolog.Logger {
	writer := output
	if writer == nil {
		var consoleWriter zerolog.ConsoleWriter
		consoleWriter.Out = os.Stderr
		consoleWriter.TimeFormat = time.RFC3339
		writer = &consoleWriter
	}

	zlevel := parseLevel(level)
	return zerolog.New(writer).
		With().
		Timestamp().
		Logger().
		Level(zlevel)
}

// SetupSlog configures the global slog to use zerolog as its backend.
func SetupSlog(zl *zerolog.Logger) {
	var slogOption slogzerolog.Option
	slogOption.Logger = zl
	handler := slogOption.NewZerologHandler()
	slog.SetDefault(slog.New(handler))
}

func parseLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.InfoLevel
	}
}
