package handlers

import (
	"fmt"
	"net/http"

	"log/slog"

	"github.com/LombardiDaniel/goliath/src/internal/domain"
	"github.com/LombardiDaniel/goliath/src/internal/dto"
	"github.com/LombardiDaniel/goliath/src/internal/middlewares"
	"github.com/LombardiDaniel/goliath/src/internal/services"
	"github.com/LombardiDaniel/goliath/src/pkg/constants"
	"github.com/LombardiDaniel/goliath/src/pkg/oauth"
	"github.com/LombardiDaniel/goliath/src/pkg/token"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService        services.AuthService
	userService        services.UserService
	emailService       services.EmailService
	oauthProvidersMap  map[string]oauth.Provider
	oauthProvidersUrls map[string]string
}

func NewAuthHandler(
	authService services.AuthService,
	userService services.UserService,
	emailService services.EmailService,
	oauthProvidersMap map[string]oauth.Provider,
) AuthHandler {
	oauthProvidersUrls := make(map[string]string)
	for k, v := range oauthProvidersMap {
		oauthProvidersUrls[k] = v.GetAuthUrl()
	}

	return AuthHandler{
		authService:        authService,
		userService:        userService,
		emailService:       emailService,
		oauthProvidersMap:  oauthProvidersMap,
		oauthProvidersUrls: oauthProvidersUrls,
	}
}

// @Summary Login
// @Tags Auth
// @Description Authenticates a user and provides a Token to Authorize API calls
// @Consume multipart/form-data
// @Produce json
// @Param email formData string true "User credentials"
// @Param password formData string true "User credentials"
// @Success 200 {object} domain.JwtClaimsOutput
// @Failure 400 string BadRequest
// @Failure 401 string Unauthorized
// @Failure 502 string BadGateway
// @Router /v1/auth/login [POST]
func (c *AuthHandler) Login(ctx *gin.Context) {
	var loginForm dto.LoginForm
	if err := ctx.ShouldBind(&loginForm); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.userService.GetUser(ctx, loginForm.Email)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while retrieving User user '%s': '%s'", loginForm.Email, err.Error()))
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	if !token.CheckPasswordHash(loginForm.Password, user.PasswordHash) {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	slog.Info(fmt.Sprintf("user login: %s", user.Email))

	t, err := c.authService.InitToken(
		ctx,
		user.UserId,
		user.Email,
		nil,
	)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while generating token for user '%s': '%s'", loginForm.Email, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	token.SetAuthCookie(ctx, t)

	claims, err := c.authService.ParseToken(t)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while parsing token for user '%s': '%s'", loginForm.Email, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.JSON(http.StatusOK, claims)
}

// @Summary Validate JWT
// @Security JWT
// @Tags Auth
// @Description Mock Endpoint to test validation of JSON Web Token (JWT) in Headers or Cookie
// @Consume application/json
// @Produce json
// @Success 200 {object} domain.JwtClaimsOutput
// @Failure 400 string BadRequest
// @Failure 401 string Unauthorized
// @Failure 502 string BadGateway
// @Router /v1/auth/validate [GET]
func (c *AuthHandler) Validate(ctx *gin.Context) {
	userClaimsRaw, ok := ctx.Get(constants.GinCtxJwtClaimKeyName)
	if !ok {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	userClaims, ok := userClaimsRaw.(domain.JwtClaims)
	if !ok {
		ctx.String(http.StatusBadGateway, "BadGateay")
		return
	}

	ctx.JSON(http.StatusOK, userClaims)
}

// @Summary Logout
// @Tags Auth
// @Description Removes the cookie
// @Success 200 string OK
// @Failure 400 string BadRequest
// @Failure 401 string Unauthorized
// @Failure 502 string BadGateway
// @Router /v1/auth/logout [POST]
func (c *AuthHandler) Logout(ctx *gin.Context) {
	token.ClearAuthCookie(ctx)
	ctx.String(http.StatusOK, "OK")
}

// @Summary SetOrg
// @Tags Auth
// @Security JWT
// @Description Sets the current User Org on JWT
// @Produce json
// @Param orgId path string true "orgId"
// @Success 200 		{object} 	domain.JwtClaimsOutput
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/auth/set-organization/{orgId} [POST]
func (c *AuthHandler) SetOrg(ctx *gin.Context) {
	orgId := ctx.Param("orgId")

	claims, err := token.GetClaimsFromGinCtx[domain.JwtClaims](ctx)
	if err != nil {
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	orgs, err := c.userService.GetUserOrgs(ctx, claims.UserId)
	if err != nil {
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	var claimsOrg *dto.OrganizationOutput = nil
	for _, org := range orgs {
		if org.OrganizationId == orgId {
			claimsOrg = &org
		}
	}

	if claimsOrg == nil {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	t, err := c.authService.InitToken(ctx, claims.UserId, claims.Email, &claimsOrg.OrganizationId)
	if err != nil {
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	parsedClaims, err := c.authService.ParseToken(t)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while parsing token for user '%s': '%s'", claims.Email, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	token.SetAuthCookie(ctx, t)
	ctx.JSON(http.StatusOK, parsedClaims)
}

// @Summary GetOauthProviders
// @Tags Auth
// @Description Gets OauthProviders and their URLs
// @Produce json
// @Success 200 		{object} 	map[string]string
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/auth/providers [GET]
func (c *AuthHandler) GetOauthProviders(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.oauthProvidersUrls)
}

// @Summary OauthCallback
// @Tags Auth
// @Description Oauth Provider Callbacks
// @Produce json
// @Param 	provider 	path 		string true "provider name"
// @Param   code 		query 		string true "code"
// @Success 302 		{string} 	OKResponse "StatusFound"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/auth/{provider}/callback [GET]
func (c *AuthHandler) OauthCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	provider, ok := c.oauthProvidersMap[ctx.Param("provider")]
	if !ok {
		ctx.String(http.StatusNotFound, "NotFound")
		return
	}

	oauthUser, err := provider.Auth(ctx, code)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	user, inserted, err := c.authService.LoginOauth(ctx, *oauthUser)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}
	if inserted {
		err = c.emailService.SendAccountCreated(user.Email, user.FirstName)
		if err != nil {
			slog.Error(err.Error())
			ctx.String(http.StatusBadGateway, "BadGateway")
			return
		}
	}

	t, err := c.authService.InitToken(ctx, user.UserId, user.Email, nil)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	token.SetAuthCookie(ctx, t)

	// ctx.Header("location", "/")
	ctx.Header("location", constants.AppHostUrl)
	ctx.String(http.StatusFound, "Found")
}

// Register Routes, needs jwtService use on authentication middleware
func (c *AuthHandler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	g := rg.Group("/auth")

	g.POST("/login", c.Login)
	g.POST("/logout", c.Logout)
	g.POST("/set-organization/:orgId", authMiddleware.AuthorizeUser(), c.SetOrg)
	g.GET("/validate", authMiddleware.AuthorizeUser(), c.Validate)

	// Oauth
	g.GET("/providers", c.GetOauthProviders)
	g.GET("/:provider/callback", c.OauthCallback)
}
