package common

import (
	"crypto/rand"
	"math/big"
	"net/url"
	"os"
	"strings"
)

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

// GetEnvVarDefault retrieves the value of the specified environment variable.
// Returns a default value if the variable is not set.
func GetEnvVarDefault(name string, def string) string {
	envVar, ok := os.LookupEnv(name)
	if !ok {
		return def
	}

	return envVar
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

// GenerateRandomString generates a random string of size n
func GenerateRandomString(n int) (string, error) {
	const choices = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(choices))))
		if err != nil {
			return "", err
		}
		b[i] = choices[num.Int64()]
	}
	return string(b), nil
}
