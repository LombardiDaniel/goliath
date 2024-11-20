package controllers

import (
	"fmt"
	"net/http"

	"log/slog"

	"github.com/LombardiDaniel/go-gin-template/common"
	"github.com/LombardiDaniel/go-gin-template/middlewares"
	"github.com/LombardiDaniel/go-gin-template/models"
	"github.com/LombardiDaniel/go-gin-template/schemas"
	"github.com/LombardiDaniel/go-gin-template/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService  services.AuthService
	userService  services.UserService
	emailService services.EmailService
}

func NewAuthController(
	authService services.AuthService,
	userService services.UserService,
	emailService services.EmailService,
) AuthController {
	return AuthController{
		authService:  authService,
		userService:  userService,
		emailService: emailService,
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

// Register Routes, needs jwtService use on authentication middleware
func (c *AuthController) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {

	g := rg.Group("/auth")

	g.POST("/login", c.Login)
	g.POST("/logout", c.Logout)
	g.POST("/set-organization/:orgId", authMiddleware.AuthorizeUser(), c.SetOrg)
	g.GET("/validate", authMiddleware.AuthorizeUser(), c.Validate)
}
