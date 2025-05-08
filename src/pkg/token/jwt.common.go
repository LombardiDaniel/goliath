package token

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/LombardiDaniel/goliath/src/pkg/constants"
	"github.com/gin-gonic/gin"
)

func GetJwtHeaderOrCookie(c *gin.Context) (string, error) {
	const BEARER_SCHEMA = "Bearer "
	authHeader := c.GetHeader("Authorization")

	tokenHeaderStr := authHeader[len(BEARER_SCHEMA):]
	tokenCookieStr, err := c.Cookie(constants.JwtCookieName)
	if err != nil && err != http.ErrNoCookie {
		slog.Error(err.Error())
	}

	if tokenHeaderStr == "" && tokenCookieStr == "" {
		return "", errors.New("ErrNoAuth")
	}

	if tokenCookieStr == "" {
		return tokenHeaderStr, nil
	}

	return tokenCookieStr, nil
}
