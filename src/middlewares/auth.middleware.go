package middlewares

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/LombardiDaniel/go-gin-template/services"
	"github.com/LombardiDaniel/go-gin-template/utils"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService services.AuthService
}

func NewAuthMiddleware(authService services.AuthService) AuthMiddleware {
	return AuthMiddleware{
		authService: authService,
	}
}

func (m *AuthMiddleware) Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// authCookie, err := ctx.Cookie(auth_cookie_token_name)
		// if err != nil {
		// 	ctx.String(http.StatusUnauthorized, "Unauthorized")
		// 	ctx.Abort()
		// 	return
		// }
		// would still need to auth

		apiKey := ctx.GetHeader("Authorization")
		if apiKey == "" {
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		tokenStrs := strings.SplitN(apiKey, " ", 2)
		if len(tokenStrs) != 2 {
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		err := m.authService.Authenticate(ctx.Request.Context(), tokenStrs[1])
		if err != nil {
			if err != utils.ErrAuth {
				slog.Error(err.Error())
			}
			ctx.String(http.StatusUnauthorized, "Unauthorized")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
