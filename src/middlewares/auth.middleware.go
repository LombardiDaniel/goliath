package middlewares

// import (
// 	"log/slog"
// 	"net/http"

// 	"github.com/LombardiDaniel/go-gin-template/common"
// 	"github.com/LombardiDaniel/go-gin-template/services"
// 	"github.com/LombardiDaniel/go-gin-template/common"
// 	"github.com/gin-gonic/gin"
// )

// type AuthMiddlewareJWT struct {
// 	authService services.AuthService
// }

// func NewAuthMiddlewareJWT(authService services.AuthService) AuthMiddlewareJWT {
// 	return AuthMiddlewareJWT{
// 		authService: authService,
// 	}
// }

// // Authorizes the JWT, if it is valid, the attribute `common.GIN_CTX_JWT_CLAIM_KEY_NAME` is set with the `models.JwtClaimsOutput`
// // allows use of JWT in cookie
// func (m *AuthMiddlewareJWT) AuthorizeJwt() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tokenStr, err := c.Cookie(common.COOKIE_NAME)
// 		if err != nil && err != http.ErrNoCookie {
// 			c.String(http.StatusUnauthorized, "Unauthorized")
// 			common.ClearAuthCookie(c)
// 			c.Abort()
// 			return
// 		}

// 		jwtClaims, err := m.authService.ParseToken(tokenStr)
// 		if err != nil {
// 			slog.Info(err.Error())
// 			c.String(http.StatusUnauthorized, "Unauthorized")
// 			common.ClearAuthCookie(c)
// 			c.Abort()
// 			return
// 		}

// 		c.Set(common.GIN_CTX_JWT_CLAIM_KEY_NAME, jwtClaims)
// 		c.Next()
// 	}
// }
