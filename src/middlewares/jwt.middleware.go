package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/LombardiDaniel/go-gin-template/common"
	"github.com/LombardiDaniel/go-gin-template/services"
	"github.com/gin-gonic/gin"
)

type AuthMiddlewareJWT struct {
	authService services.AuthService
}

func NewAuthMiddlewareJWT(authService services.AuthService) AuthMiddlewareJWT {
	return AuthMiddlewareJWT{
		authService: authService,
	}
}

// Authorizes the JWT, if it is valid, the attribute `common.GIN_CTX_JWT_CLAIM_KEY_NAME` is set with the `models.JwtClaimsOutput`
// allows use of JWT in cookie
func (m *AuthMiddlewareJWT) AuthorizeJwt() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(common.COOKIE_NAME)
		if err != nil && err != http.ErrNoCookie {
			c.String(http.StatusUnauthorized, "Unauthorized")
			common.ClearAuthCookie(c)
			c.Abort()
			return
		}

		jwtClaims, err := m.authService.ParseToken(tokenStr)
		if err != nil {
			slog.Info(err.Error())
			c.String(http.StatusUnauthorized, "Unauthorized")
			common.ClearAuthCookie(c)
			c.Abort()
			return
		}

		// Renew Cycle:
		expTime := time.Unix(jwtClaims.ExpiresAt, 0)

		expTTL := time.Until(expTime)

		if expTTL > time.Minute*time.Duration(common.JWT_TIMEOUT_SECS/2) {
			slog.Info(fmt.Sprintf("renewing jwt: %s", jwtClaims.Email))
			token, err := m.authService.InitToken(jwtClaims.UserId, jwtClaims.Email, jwtClaims.OrganizationId)
			if err != nil {
				slog.Error(err.Error())
				c.String(http.StatusBadGateway, "BadGateway")
				common.ClearAuthCookie(c)
				c.Abort()
				return
			}

			common.SetAuthCookie(c, token)
		}

		c.Set(common.GIN_CTX_JWT_CLAIM_KEY_NAME, jwtClaims)
		c.Next()
	}
}
