package logger

import (
	"log/slog"
	"os"

	"github.com/chickenzord/sparkyfitness-mcp/internal/config"
)

// InitLogger initializes the global slog logger based on configuration
func InitLogger(cfg *config.Config) {
	// Determine log level
	var level slog.Level
	switch cfg.LogLevel {
	case config.LogLevelDebug:
		level = slog.LevelDebug
	case config.LogLevelInfo:
		level = slog.LevelInfo
	case config.LogLevelWarn:
		level = slog.LevelWarn
	case config.LogLevelError:
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Create handler options
	handlerOpts := &slog.HandlerOptions{
		Level: level,
	}

	// Create handler based on format
	var handler slog.Handler
	switch cfg.LogFormat {
	case config.LogFormatJSON:
		handler = slog.NewJSONHandler(os.Stderr, handlerOpts)
	case config.LogFormatText:
		fallthrough
	default:
		handler = slog.NewTextHandler(os.Stderr, handlerOpts)
	}

	// Set as default logger
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
