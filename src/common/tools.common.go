package common

import (
	"log/slog"
	"net/url"
	"os"
	"strings"
	"text/template"
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

func LoadHTMLTemplate(templatePath string) *template.Template {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	return t
}

// Gets the `envVarName`, returns defaultVal if envvar is non-existant.
func GetEnvVarDefault(envVarName string, defaultVal string) string {
	envVar, ok := os.LookupEnv(envVarName)

	if !ok {
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

func ExtractHostFromUrl(rawUrl string) (string, error) {
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	host := parsedURL.Hostname()

	return host, nil
}

func UrlIsSecure(rawUrl string) (bool, error) {
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return false, err
	}

	return parsedURL.Scheme == "https", nil
}
