package common

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	secsTimeoutJWT, _        = strconv.Atoi(GetEnvVarDefault("JWT_TIMEOUT_SECS", "43200"))
	ENV               string = os.Getenv("ENVIRONMENT")
	secure            bool   = ENV != ""
	domain            string = ""
)

func SetAuthCookie(ctx *gin.Context, tokenStr string) {
	ctx.Header("Set-Cookie", makeAuthCookie(tokenStr, domain))
}

func ClearAuthCookie(ctx *gin.Context) {
	ctx.Header("Set-Cookie", makeClearAuthCookie(domain))
}

func makeAuthCookie(value string, domain string) string {
	return makeCookie(COOKIE_NAME, value, secsTimeoutJWT, "/", domain, secure, true)
}

func makeClearAuthCookie(domain string) string {
	return makeCookie(COOKIE_NAME, "", secsTimeoutJWT, "/", domain, secure, true)
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
