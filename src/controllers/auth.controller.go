package controllers

import (
	"fmt"
	"net/http"

	"log/slog"

	"github.com/LombardiDaniel/go-gin-template/common"
	"github.com/LombardiDaniel/go-gin-template/middlewares"
	"github.com/LombardiDaniel/go-gin-template/models"
	"github.com/LombardiDaniel/go-gin-template/oauth"
	"github.com/LombardiDaniel/go-gin-template/schemas"
	"github.com/LombardiDaniel/go-gin-template/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService        services.AuthService
	userService        services.UserService
	emailService       services.EmailService
	oauthProvidersMap  map[string]oauth.Provider
	oauthProvidersUrls map[string]string
}

func NewAuthController(
	authService services.AuthService,
	userService services.UserService,
	emailService services.EmailService,
	oauthProvidersMap map[string]oauth.Provider,
) AuthController {
	oauthProvidersUrls := make(map[string]string)
	for k, v := range oauthProvidersMap {
		oauthProvidersUrls[k] = v.GetAuthUrl()
	}

	return AuthController{
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
// @Success 200 {object} models.JwtClaimsOutput
// @Failure 400 string BadRequest
// @Failure 401 string Unauthorized
// @Failure 502 string BadGateway
// @Router /v1/auth/login [POST]
func (c *AuthController) Login(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	var loginForm schemas.LoginForm
	if err := ctx.ShouldBind(&loginForm); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.userService.GetUser(rCtx, loginForm.Email)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while retrieving User user '%s': '%s'", loginForm.Email, err.Error()))
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	if !common.CheckPasswordHash(loginForm.Password, user.PasswordHash) {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	slog.Info(fmt.Sprintf("user login: %s", user.Email))

	token, err := c.authService.InitToken(
		user.UserId,
		user.Email,
		nil,
		nil,
	)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while generating token for user '%s': '%s'", loginForm.Email, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	common.SetAuthCookie(ctx, token)

	claims, err := c.authService.ParseToken(token)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while parsing token for user '%s': '%s'", loginForm.Email, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.JSON(http.StatusOK, claims)
}

// @Security JWT
// @Tags Auth
// @Summary Validate JWT
// @Description Mock Endpoint to test validation of JSON Web Token (JWT) in Headers or Cookie
// @Consume application/json
// @Produce json
// @Success 200 {object} models.JwtClaimsOutput
// @Failure 400 string BadRequest
// @Failure 401 string Unauthorized
// @Failure 502 string BadGateway
// @Router /v1/auth/validate [GET]
func (c *AuthController) Validate(ctx *gin.Context) {
	userClaimsRaw, ok := ctx.Get(common.GIN_CTX_JWT_CLAIM_KEY_NAME)
	if !ok {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	userClaims, ok := userClaimsRaw.(models.JwtClaims)
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
func (c *AuthController) Logout(ctx *gin.Context) {
	common.ClearAuthCookie(ctx)
	ctx.String(http.StatusOK, "OK")
}

// @Summary SetOrg
// @Tags Auth
// @Security JWT
// @Description Sets the current User Org on JWT
// @Produce json
// @Param orgId path string true "orgId"
// @Success 200 		{object} 	models.JwtClaimsOutput
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/auth/set-organization/{orgId} [POST]
func (c *AuthController) SetOrg(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	orgId := ctx.Param("orgId")

	claims, err := common.GetClaimsFromGinCtx(ctx)
	if err != nil {
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	orgs, err := c.userService.GetUserOrgs(rCtx, claims.UserId)
	if err != nil {
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	var claimsOrg *schemas.OrganizationOutput = nil
	for _, org := range orgs {
		if org.OrganizationId == orgId {
			claimsOrg = &org
		}
	}

	if claimsOrg == nil {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	token, err := c.authService.InitToken(claims.UserId, claims.Email, &claimsOrg.OrganizationId, &claimsOrg.IsAdmin)
	if err != nil {
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	parsedClaims, err := c.authService.ParseToken(token)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while parsing token for user '%s': '%s'", claims.Email, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	common.SetAuthCookie(ctx, token)
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
func (c *AuthController) GetOauthProviders(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.oauthProvidersUrls)
}

// @Summary GetOauthProviders
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
func (c *AuthController) OauthCallback(ctx *gin.Context) {
	rCtx := ctx.Request.Context()

	code := ctx.Query("code")
	provider, ok := c.oauthProvidersMap[ctx.Param("provider")]
	if !ok {
		ctx.String(http.StatusNotFound, "NotFound")
		return
	}

	oauthUser, err := provider.Auth(rCtx, code)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	user, inserted, err := c.userService.LoginOauth(ctx, *oauthUser)
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

	token, err := c.authService.InitToken(user.UserId, user.Email, nil, nil)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	common.SetAuthCookie(ctx, token)

	ctx.Header("location", "/")
	ctx.String(http.StatusFound, "Found")
}

// Register Routes, needs jwtService use on authentication middleware
func (c *AuthController) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {

	g := rg.Group("/auth")

	g.POST("/login", c.Login)
	g.POST("/logout", c.Logout)
	g.POST("/set-organization/:orgId", authMiddleware.AuthorizeUser(), c.SetOrg)
	g.GET("/validate", authMiddleware.AuthorizeUser(), c.Validate)

	// Oauth
	g.GET("/providers", c.GetOauthProviders)
	g.GET("/:provider/callback", c.OauthCallback)
}
