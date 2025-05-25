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
	domain  string = getCookieDomain()
)

func getCookieDomain() string {
	cookieDomain := ""
	if secure {
		cookieDomain = strings.SplitN(constants.ApiHostUrl, "://", 2)[1]
		if cookieDomain[len(cookieDomain)-1] == '/' {
			cookieDomain = cookieDomain[0 : len(cookieDomain)-1]
		}
	}

	// slog.Info("Cookie Domain: " + cookieDomain)

	return cookieDomain
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
		makeCookie(cookieName, value, constants.JwtTimeoutSecs, "/", domain, secure, true),
	)

}

func SetAuthCookie(ctx *gin.Context, token string) {
	ctx.Header(
		"Set-Cookie",
		makeAuthCookie(token, domain),
	)
}

func ClearAuthCookie(ctx *gin.Context) {
	ctx.Header(
		"Set-Cookie",
		makeAuthCookie("", domain),
	)
}

func makeAuthCookie(value string, domain string) string {
	return makeCookie(constants.JwtCookieName, value, constants.JwtTimeoutSecs, "/", domain, secure, true)
}

func makeCookie(name string, value string, maxAge int, path string, domain string, secure bool, httpOnly bool) string {
	cookieStr := ""

	cookieStr += name + "=" + value + "; "
	cookieStr += "Path" + "=" + path + "; "
	cookieStr += "Max-Age" + "=" + strconv.Itoa(maxAge) + "; "

	if domain != "" {
		cookieStr += "Domain" + "=" + domain + "; "
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
