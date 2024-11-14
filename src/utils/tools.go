package utils

import (
	"log/slog"
	"os"
	"strings"
)

var LOG_LEVEL string = strings.ToUpper(GetEnvVarDefault("LOG_LEVEL", "INFO"))

func InitSlogger() {

	levelsMap := map[string]slog.Level{
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
			Level:     levelsMap[LOG_LEVEL],
		},
	))

	slog.SetDefault(logger)
}

// Gets the `envVarName`, returns defaultVal if envvar is non-existant.
func GetEnvVarDefault(envVarName string, defaultVal string) string {
	envVar := os.Getenv(envVarName)

	if envVar == "" {
		return defaultVal
	}

	return envVar
}

// Removes all occurences of item in slice
func RemoveFrom[T comparable](slice []T, item T) []T {
	var newSlice []T
	for _, v := range slice {
		if v != item {
			newSlice = append(newSlice, v)
		}
	}

	return newSlice
}

func IsSubset(subset []string, superset []string) bool {
	checkMap := make(map[string]bool)
	for _, element := range superset {
		checkMap[element] = true
	}
	for _, value := range subset {
		if !checkMap[value] {
			return false // Return false if an element is not found in the superset
		}
	}
	return true // Return true if all elements are found in the superset
}
