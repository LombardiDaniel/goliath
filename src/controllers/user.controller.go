package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/LombardiDaniel/go-gin-template/common"
	"github.com/LombardiDaniel/go-gin-template/middlewares"
	"github.com/LombardiDaniel/go-gin-template/models"
	"github.com/LombardiDaniel/go-gin-template/schemas"
	"github.com/LombardiDaniel/go-gin-template/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	authService  services.AuthService
	userService  services.UserService
	emailService services.EmailService
}

func NewUserController(
	authService services.AuthService,
	userService services.UserService,
	emailService services.EmailService,
) UserController {
	return UserController{
		authService:  authService,
		userService:  userService,
		emailService: emailService,
	}
}

// @Summary CreateUser
// @Tags User
// @Description Creates an User
// @Consume application/json
// @Accept json
// @Produce plain
// @Param   payload 	body 		schemas.CreateUser true "user json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users [PUT]
func (c *UserController) CreateUser(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	var createUser schemas.CreateUser

	if err := ctx.ShouldBind(&createUser); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	hash, err := common.HashPassword(createUser.Password)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while hashing pass '%s': '%s'", createUser.Password, err.Error()))
		ctx.String(http.StatusBadRequest, "BadRequest")
		return
	}

	otp, err := common.GenerateRandomString(common.OTP_LEN)
	if err != nil {
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	unconfirmedUser := models.UnconfirmedUser{
		Email:        createUser.Email,
		Otp:          otp,
		PasswordHash: hash,
		FirstName:    createUser.FirstName,
		LastName:     createUser.LastName,
		DateOfBirth:  createUser.DateOfBirth,
	}

	err = c.userService.CreateUnconfirmedUser(rCtx, unconfirmedUser)
	if err == common.ErrDbConflict {
		ctx.String(http.StatusConflict, "Conflict")
		return
	}

	if err != nil {
		slog.Error(fmt.Sprintf("Error while creating unconfirmedUser '%s': '%s'", unconfirmedUser.Email, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	err = c.emailService.SendEmailConfirmation(unconfirmedUser.Email, unconfirmedUser.FirstName, unconfirmedUser.Otp)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while sending email '%s': '%s'", unconfirmedUser.Email, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

// @Summary ConfirmUser
// @Tags User
// @Description Confirms the User
// @Produce plain
// @Param   otp 		query 		string true "OneTimePass sent in email"
// @Success 301 		{string} 	OKResponse "MovedPermanently"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/confirm [GET]
func (c *UserController) ConfirmUser(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	otp := ctx.Query("otp")

	err := c.userService.ConfirmUser(rCtx, otp)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while confirming user otp='%s': '%s'", otp, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.Header("location", "/")
	ctx.String(http.StatusOK, "OK")
}

// @Summary GetUserOrgs
// @Tags User
// @Security JWT
// @Description Gets Orgs Users Belongs to
// @Produce json
// @Success 200 		{object} 	[]schemas.OrganizationOutput
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/organizations [GET]
func (c *UserController) GetUserOrgs(ctx *gin.Context) {
	rCtx := ctx.Request.Context()

	claims, err := common.GetClaimsFromGinCtx(ctx)
	if err != nil {
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	orgs, err := c.userService.GetUserOrgs(rCtx, claims.UserId)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.JSON(http.StatusOK, orgs)
}

// @Summary InitResetPassword
// @Tags User
// @Description Inits the password reset pipeline
// @Consume application/json
// @Accept json
// @Produce plain
// @Param   payload 	body 		schemas.Email true "email json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/init-reset-password [POST]
func (c *UserController) InitResetPassword(ctx *gin.Context) {
	rCtx := ctx.Request.Context()

	var email schemas.Email

	if err := ctx.ShouldBind(&email); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	otp, err := common.GenerateRandomString(common.OTP_LEN)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.userService.GetUser(rCtx, email.Email)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	err = c.userService.InitPasswordReset(rCtx, user.UserId, otp)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	err = c.emailService.SendPasswordReset(email.Email, user.FirstName, otp)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	ctx.String(http.StatusOK, "OK")
}

// @Summary SetPasswordResetCookie
// @Tags User
// @Description Sets te Passord Reset Cookie
// @Produce plain
// @Param   otp 		query 		string true "OneTimePass sent in email"
// @Success 301 		{string} 	OKResponse "MovedPermanently"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/set-password-reset-cookie [GET]
func (c *UserController) SetPasswordResetCookie(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	otp := ctx.Query("otp")

	reset, err := c.userService.GetPasswordReset(rCtx, otp)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while setting cookie user otp='%s': '%s'", otp, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	tokenStr, err := c.authService.InitPasswordResetToken(reset.UserId)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while confirming user otp='%s': '%s'", otp, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	common.SetCookieForApp(ctx, common.PASSWORD_RESET_TIMEOUT_JWT_COOKIE_NAME, tokenStr)

	ctx.Header("location", "/reset-password") // Sets on browser
	ctx.String(http.StatusOK, "OK")
}

// @Summary ResetPassword
// @Tags User
// @Description Resets the password (auths via special cookie)
// @Consume application/json
// @Accept json
// @Produce plain
// @Param   payload 	body 		schemas.Password true "pw json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/reset-password [POST]
func (c *UserController) ResetPassword(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	var pw schemas.Password

	if err := ctx.ShouldBind(&pw); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	cookieVal, err := ctx.Cookie(common.PASSWORD_RESET_TIMEOUT_JWT_COOKIE_NAME)
	if err != nil {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	claims, err := c.authService.ParsePasswordResetToken(cookieVal)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	err = c.userService.UpdateUserPassword(rCtx, claims.UserId, pw.Password)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	common.SetCookieForApp(ctx, common.PASSWORD_RESET_TIMEOUT_JWT_COOKIE_NAME, "")
	ctx.String(http.StatusOK, "OK")
}

func (c *UserController) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	g := rg.Group("/users")

	g.PUT("", c.CreateUser)
	g.GET("/confirm", c.ConfirmUser)
	g.POST("/init-reset-password", c.InitResetPassword)
	g.GET("/set-password-reset-cookie", c.SetPasswordResetCookie)
	g.POST("/reset-password", c.ResetPassword)
	g.GET("/organizations", authMiddleware.AuthorizeUser(), c.GetUserOrgs)
}
