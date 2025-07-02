package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/LombardiDaniel/goliath/src/internal/models"
	"github.com/LombardiDaniel/goliath/src/internal/services"
	"github.com/LombardiDaniel/goliath/src/pkg/constants"
	"github.com/LombardiDaniel/goliath/src/pkg/token"

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

// Authorizes the JWT, if it is valid, the attribute `constants.GinCtxJwtClaimKeyName` is set with the `models.JwtClaimsOutput`
// allows use of JWT in cookie
func (m *AuthMiddlewareJwt) AuthorizeUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(constants.JwtCookieName)
		if err != nil && err != http.ErrNoCookie {
			c.String(http.StatusUnauthorized, "Unauthorized")
			token.ClearAuthCookie(c)
			c.Abort()
			return
		}

		jwtClaims, err := m.authService.ParseToken(tokenStr)
		if err != nil {
			slog.Info(err.Error())
			c.String(http.StatusUnauthorized, "Unauthorized")
			token.ClearAuthCookie(c)
			c.Abort()
			return
		}

		// Renew Cycle:
		expTime := time.Unix(jwtClaims.ExpiresAt, 0)

		expTTL := time.Until(expTime)

		if expTTL > time.Minute*time.Duration(constants.JwtTimeoutSecs/2) {
			slog.Info(fmt.Sprintf("renewing jwt: %s", jwtClaims.Email))
			t, err := m.authService.InitToken(c, jwtClaims.UserId, jwtClaims.Email, jwtClaims.OrganizationId)
			if err != nil {
				slog.Error(err.Error())
				c.String(http.StatusBadGateway, "BadGateway")
				token.ClearAuthCookie(c)
				c.Abort()
				return
			}

			token.SetAuthCookie(c, t)
		}

		c.Set(constants.GinCtxJwtClaimKeyName, jwtClaims)
		c.Next()
	}
}

func (m *AuthMiddlewareJwt) AuthorizeOrganization(need map[string]models.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(constants.JwtCookieName)
		if err != nil && err != http.ErrNoCookie {
			c.String(http.StatusUnauthorized, "Unauthorized")
			token.ClearAuthCookie(c)
			c.Abort()
			return
		}

		jwtClaims, err := m.authService.ParseToken(tokenStr)
		if err != nil {
			slog.Info(err.Error())
			c.String(http.StatusUnauthorized, "Unauthorized")
			token.ClearAuthCookie(c)
			c.Abort()
			return
		}

		orgId := c.Param("orgId")
		if jwtClaims.OrganizationId == nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			token.ClearAuthCookie(c)
			c.Abort()
			return
		}

		if orgId != *jwtClaims.OrganizationId {
			c.String(http.StatusUnauthorized, "Unauthorized")
			token.ClearAuthCookie(c)
			c.Abort()
			return
		}

		for action, needPerms := range need {
			if needPerms&jwtClaims.Perms[action] != needPerms { // simple bitwise ops for perms
				c.String(http.StatusUnauthorized, "Unauthorized")
				token.ClearAuthCookie(c)
				c.Abort()
				return
			}
		}

		// Renew Cycle:
		expTime := time.Unix(jwtClaims.ExpiresAt, 0)

		expTTL := time.Until(expTime)

		if expTTL > time.Minute*time.Duration(constants.JwtTimeoutSecs/2) {
			slog.Info(fmt.Sprintf("renewing jwt: %s", jwtClaims.Email))
			t, err := m.authService.InitToken(c, jwtClaims.UserId, jwtClaims.Email, jwtClaims.OrganizationId)
			if err != nil {
				slog.Error(err.Error())
				c.String(http.StatusBadGateway, "BadGateway")
				token.ClearAuthCookie(c)
				c.Abort()
				return
			}

			token.SetAuthCookie(c, t)
		}

		c.Set(constants.GinCtxJwtClaimKeyName, jwtClaims)
		c.Next()
	}
}

func (m *AuthMiddlewareJwt) Reauthorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		jwtClaims, err := token.GetClaimsFromGinCtx[models.JwtClaims](c)
		if err != nil {
			slog.Error(err.Error())
			c.String(http.StatusBadGateway, "BadGateway")
			token.ClearAuthCookie(c)
			c.Abort()
			return
		}

		slog.Info(fmt.Sprintf("renewing jwt: %s", jwtClaims.Email))
		t, err := m.authService.InitToken(c, jwtClaims.UserId, jwtClaims.Email, jwtClaims.OrganizationId)
		if err != nil {
			slog.Error(err.Error())
			c.String(http.StatusBadGateway, "BadGateway")
			token.ClearAuthCookie(c)
			c.Abort()
			return
		}

		token.SetAuthCookie(c, t)

		c.Set(constants.GinCtxJwtClaimKeyName, jwtClaims)
		c.Next()
	}
}
