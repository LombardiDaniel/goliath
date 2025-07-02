package token

import (
	"os"
	"strconv"
	"strings"

	"github.com/LombardiDaniel/goliath/src/pkg/constants"
	"github.com/gin-gonic/gin"
)

var (
	ginMode string = os.Getenv("GIN_MODE")
	secure  bool   = ginMode == "release"
	models  string = getCookiemodels()
)

func getCookiemodels() string {
	cookiemodels := ""
	if secure {
		cookiemodels = strings.SplitN(constants.ApiHostUrl, "://", 2)[1]
		if cookiemodels[len(cookiemodels)-1] == '/' {
			cookiemodels = cookiemodels[0 : len(cookiemodels)-1]
		}
	}

	// slog.Info("Cookie models: " + cookiemodels)

	return cookiemodels
}

func GetClaimsFromGinCtx[T any](ctx *gin.Context) (T, error) {
	claims, ok := ctx.Get(constants.GinCtxJwtClaimKeyName)
	var zero T
	if !ok {
		return zero, constants.ErrAuth
	}

	parsedClaims, ok := claims.(T)
	if !ok {
		return zero, constants.ErrAuth
	}

	return parsedClaims, nil
}

func SetCookieForApp(ctx *gin.Context, cookieName string, value string) {
	ctx.Header(
		"Set-Cookie",
		makeCookie(cookieName, value, constants.JwtTimeoutSecs, "/", models, secure, true),
	)

}

func SetAuthCookie(ctx *gin.Context, token string) {
	ctx.Header(
		"Set-Cookie",
		makeAuthCookie(token, models),
	)
}

func ClearAuthCookie(ctx *gin.Context) {
	ctx.Header(
		"Set-Cookie",
		makeAuthCookie("", models),
	)
}

func makeAuthCookie(value string, models string) string {
	return makeCookie(constants.JwtCookieName, value, constants.JwtTimeoutSecs, "/", models, secure, true)
}

func makeCookie(name string, value string, maxAge int, path string, models string, secure bool, httpOnly bool) string {
	cookieStr := ""

	cookieStr += name + "=" + value + "; "
	cookieStr += "Path" + "=" + path + "; "
	cookieStr += "Max-Age" + "=" + strconv.Itoa(maxAge) + "; "

	if models != "" {
		cookieStr += "models" + "=" + models + "; "
	}

	if httpOnly {
		cookieStr += "HttpOnly; "
	}

	if secure {
		cookieStr += "Secure; "
	}

	cookieStr += "SameSite" + "=Lax;"

	return cookieStr
}
