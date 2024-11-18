package middlewares

// import (
// 	"log/slog"
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"ticktr.ai/api-service/services"
// 	"ticktr.ai/api-service/common"
// )

// type AuthMiddleware struct {
// 	authService services.AuthService
// }

// func NewAuthMiddleware(authService services.AuthService) AuthMiddleware {
// 	return AuthMiddleware{
// 		authService: authService,
// 	}
// }

// func (m *AuthMiddleware) Authorize(needAdmin bool, allowApi bool) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		// First check for JWT in cookies (faster), then try ApiKey
// 		oID := ctx.Param("organizationID")

// 		hasValidJWT := false
// 		hasValidApiKey := false

// 		// JWT
// 		JWT, err := ctx.Cookie(common.JWT_COOKIE_NAME)
// 		if err == nil {
// 			jwtClaims, err := m.authService.AuthenticateJWT(oID, needAdmin, JWT)
// 			if err == nil {
// 				hasValidJWT = true
// 				ctx.Set(common.CTX_CLAIM_KEY_NAME, jwtClaims)
// 			}
// 		}

// 		// ApiKey
// 		apiKey := ctx.GetHeader("Authorization")
// 		if apiKey != "" {
// 			err = m.authService.AuthenticateAPI(oID, apiKey)
// 			if err == nil {
// 				hasValidApiKey = true
// 				ctx.Set(common.CTX_BPK_KEY_NAME, oID)
// 			}
// 		}

// 		if !hasValidJWT && !hasValidApiKey {
// 			slog.Debug(err.Error())
// 			ctx.String(http.StatusUnauthorized, "Unauthorized")
// 			ctx.Abort()
// 			return
// 		}

// 		ctx.Next()
// 	}
// }

// func (m *AuthMiddleware) AuthorizeNoOrganization() gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		authHeader := ctx.GetHeader("Authorization")

// 		authHeaderSplit := strings.Split(authHeader, " ")
// 		if len(authHeaderSplit) != 2 {
// 			ctx.String(http.StatusUnauthorized, "Unauthorized")
// 			ctx.Abort()
// 			return
// 		}

// 		authToken := authHeaderSplit[1]

// 		jwtClaims, err := m.authService.AuthenticateNoOrganizationJWT(authToken)
// 		if err != nil {
// 			ctx.String(http.StatusUnauthorized, "Unauthorized")
// 			ctx.Abort()
// 			return
// 		}

// 		if jwtClaims.Email != ctx.Param("email") {
// 			ctx.String(http.StatusUnauthorized, "Unauthorized")
// 			ctx.Abort()
// 			return
// 		}

// 		ctx.Set(common.CTX_CLAIM_KEY_NAME, jwtClaims)

// 		ctx.Next()
// 	}
// }
