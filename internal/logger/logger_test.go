package logger_test

import (
	"bytes"
	"testing"

	"github.com/omarluq/hetzner-blackbsd/internal/logger"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("creates logger with custom writer", func(t *testing.T) {
		t.Parallel()

		var buf bytes.Buffer
		zl := logger.New("info", &buf)

		zl.Info().Msg("test message")

		assert.Contains(t, buf.String(), "test message")
	})

	t.Run("creates logger with nil writer uses stderr console", func(t *testing.T) {
		t.Parallel()

		zl := logger.New("debug", nil)

		assert.Equal(t, zerolog.DebugLevel, zl.GetLevel())
	})
}

func TestParseLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected zerolog.Level
	}{
		{
			name:     "trace level",
			input:    "trace",
			expected: zerolog.TraceLevel,
		},
		{
			name:     "debug level",
			input:    "debug",
			expected: zerolog.DebugLevel,
		},
		{
			name:     "info level",
			input:    "info",
			expected: zerolog.InfoLevel,
		},
		{
			name:     "warn level",
			input:    "warn",
			expected: zerolog.WarnLevel,
		},
		{
			name:     "error level",
			input:    "error",
			expected: zerolog.ErrorLevel,
		},
		{
			name:     "unknown defaults to info",
			input:    "unknown",
			expected: zerolog.InfoLevel,
		},
		{
			name:     "empty defaults to info",
			input:    "",
			expected: zerolog.InfoLevel,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, logger.ParseLevel(tt.input))
		})
	}
}

func TestSetupSlog(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	zl := logger.New("info", &buf)
	logger.SetupSlog(&zl)
}
