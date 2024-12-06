package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/LombardiDaniel/gopherbase/common"
	"github.com/LombardiDaniel/gopherbase/fiddlers"
	"github.com/LombardiDaniel/gopherbase/services"
	"github.com/gin-gonic/gin"
)

type AuthMiddlewareJwt struct {
	authService services.AuthService
}

func NewAuthMiddlewareJwt(authService services.AuthService) AuthMiddleware {
	return &AuthMiddlewareJwt{
		authService: authService,
	}
}

// Authorizes the JWT, if it is valid, the attribute `common.GIN_CTX_JWT_CLAIM_KEY_NAME` is set with the `models.JwtClaimsOutput`
// allows use of JWT in cookie
func (m *AuthMiddlewareJwt) AuthorizeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(common.JWT_COOKIE_NAME)
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
			token, err := m.authService.InitToken(jwtClaims.UserId, jwtClaims.Email, jwtClaims.OrganizationId, jwtClaims.IsAdmin)
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

func (m *AuthMiddlewareJwt) AuthorizeOrganization(needAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(common.JWT_COOKIE_NAME)
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

		orgId := c.Param("orgId")
		if jwtClaims.OrganizationId == nil || jwtClaims.IsAdmin == nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			common.ClearAuthCookie(c)
			c.Abort()
			return
		}

		if orgId != *jwtClaims.OrganizationId {
			c.String(http.StatusUnauthorized, "Unauthorized")
			common.ClearAuthCookie(c)
			c.Abort()
			return
		}

		if needAdmin && !*jwtClaims.IsAdmin {
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
			token, err := m.authService.InitToken(jwtClaims.UserId, jwtClaims.Email, jwtClaims.OrganizationId, jwtClaims.IsAdmin)
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

func (m *AuthMiddlewareJwt) Reauthorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtClaims, err := fiddlers.GetClaimsFromGinCtx(c)
		if err != nil {
			slog.Error(err.Error())
			c.String(http.StatusBadGateway, "BadGateway")
			common.ClearAuthCookie(c)
			c.Abort()
			return
		}

		slog.Info(fmt.Sprintf("renewing jwt: %s", jwtClaims.Email))
		token, err := m.authService.InitToken(jwtClaims.UserId, jwtClaims.Email, jwtClaims.OrganizationId, jwtClaims.IsAdmin)
		if err != nil {
			slog.Error(err.Error())
			c.String(http.StatusBadGateway, "BadGateway")
			common.ClearAuthCookie(c)
			c.Abort()
			return
		}

		common.SetAuthCookie(c, token)

		c.Set(common.GIN_CTX_JWT_CLAIM_KEY_NAME, jwtClaims)
		c.Next()
	}
}
