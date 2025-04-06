package common

import (
	"log/slog"
	"net/url"
	"os"
	"strings"
	"text/template"
)

// LOG_LEVEL defines the logging level for the application, defaulting to "INFO".
var LOG_LEVEL string = strings.ToUpper(GetEnvVarDefault("LOG_LEVEL", "INFO"))

// InitSlogger initializes the global logger with the specified log level.
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

// SplitName splits a full name into first and last name. If no last name exists, it returns an empty string.
func SplitName(fullName string) (string, string) {
	names := strings.SplitN(fullName, " ", 2)

	if len(names) == 0 {
		return "", ""
	}

	firstName := names[0]

	if len(names) == 1 {
		return firstName, ""
	}

	return firstName, names[1]
}

// LoadHTMLTemplate loads an HTML template from the specified file path. Panics if the file cannot be loaded.
func LoadHTMLTemplate(templatePath string) *template.Template {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	return t
}

// GetEnvVarDefault retrieves the value of the specified environment variable. Returns a default value if the variable is not set.
func GetEnvVarDefault(envVarName string, defaultVal string) string {
	envVar, ok := os.LookupEnv(envVarName)

	if !ok {
		return defaultVal
	}

	return envVar
}

// RemoveFrom removes all occurrences of a specified item from a slice.
func RemoveFrom[T comparable](slice []T, item T) []T {
	var newSlice []T
	for _, v := range slice {
		if v != item {
			newSlice = append(newSlice, v)
		}
	}

	return newSlice
}

// IsSubset checks if all elements of the subset are present in the superset.
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

// ExtractHostFromUrl extracts the hostname from a given URL string. Returns an error if the URL is invalid.
func ExtractHostFromUrl(rawUrl string) (string, error) {
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	host := parsedURL.Hostname()

	return host, nil
}

// UrlIsSecure checks if a given URL uses the HTTPS scheme. Returns an error if the URL is invalid.
func UrlIsSecure(rawUrl string) (bool, error) {
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return false, err
	}

	return parsedURL.Scheme == "https", nil
}

// Batch collects items from a channel into a slice of the specified size. Stops if the channel is closed or empty.
// Does NOT block if the channel is empty.
func Batch[T any](ch <-chan T, size uint32) []T {
	batch := make([]T, 0, size)

	for range size {
		select {
		case item, ok := <-ch:
			if !ok {
				// Channel is closed, return the batch collected so far
				return batch
			}
			batch = append(batch, item)
		default:
			// Channel has no data ready — return immediately
			return batch
		}
	}

	return batch
}
