package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/LombardiDaniel/goliath/src/pkg/common"
)

func init() {
	InitSlogger()
}

// LOG_LEVEL defines the logging level for the application, defaulting to "INFO".
var LOG_LEVEL string = strings.ToUpper(common.GetEnvVarDefault("LOG_LEVEL", "INFO"))

// InitSlogger initializes the global logger with the specified log level.
func InitSlogger() {
	lvls := map[string]slog.Level{
		"DEBUG":   slog.LevelDebug,
		"INFO":    slog.LevelInfo,
		"WARN":    slog.LevelWarn,
		"WARNING": slog.LevelWarn,
		"ERROR":   slog.LevelError,
	}

	logger := slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     lvls[LOG_LEVEL],
		},
	))

	slog.SetDefault(logger)
}
