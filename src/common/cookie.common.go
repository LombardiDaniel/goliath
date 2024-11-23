package common

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

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
		cookieDomain = strings.Split(API_HOST_URL, "://")[1]
		if cookieDomain[len(cookieDomain)-1] == '/' {
			cookieDomain = cookieDomain[0 : len(cookieDomain)-1]
		}
	}

	slog.Info("Cookie Domain: " + cookieDomain)

	return cookieDomain
}

// func SetCookieForApp(ctx *gin.Context, cookieName string, value string) {
// 	for _, domain := range cookieDomains {
// 		ctx.Header(
// 			"Set-Cookie",
// 			makeCookie(cookieName, value, JWT_TIMEOUT_SECS, "/", domain, secure, true),
// 		)
// 	}
// }

func SetCookieForApp(ctx *gin.Context, cookieName string, value string) {
	ctx.Header(
		"Set-Cookie",
		makeCookie(cookieName, value, JWT_TIMEOUT_SECS, "/", domain, secure, true),
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
	return makeCookie(JWT_COOKIE_NAME, value, JWT_TIMEOUT_SECS, "/", domain, secure, true)
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
