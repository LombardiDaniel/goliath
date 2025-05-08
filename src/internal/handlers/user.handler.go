package handlers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LombardiDaniel/gopherbase/src/internal/domain"
	"github.com/LombardiDaniel/gopherbase/src/internal/dto"
	"github.com/LombardiDaniel/gopherbase/src/internal/middlewares"
	"github.com/LombardiDaniel/gopherbase/src/internal/services"
	"github.com/LombardiDaniel/gopherbase/src/pkg/common"
	"github.com/LombardiDaniel/gopherbase/src/pkg/constants"
	"github.com/LombardiDaniel/gopherbase/src/pkg/storage"
	"github.com/LombardiDaniel/gopherbase/src/pkg/token"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	authService  services.AuthService
	userService  services.UserService
	emailService services.EmailService
	objService   services.ObjectService
}

func NewUserHandler(
	authService services.AuthService,
	userService services.UserService,
	emailService services.EmailService,
	objService services.ObjectService,
) UserHandler {
	return UserHandler{
		authService:  authService,
		userService:  userService,
		emailService: emailService,
		objService:   objService,
	}
}

// @Summary CreateUser
// @Tags User
// @Description Creates an User
// @Consume application/json
// @Accept json
// @Produce plain
// @Param   payload 	body 		dto.CreateUser true "user json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users [PUT]
func (c *UserHandler) CreateUser(ctx *gin.Context) {
	var createUser dto.CreateUser

	if err := ctx.ShouldBind(&createUser); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	unconfirmedUser, err := domain.NewUnconfirmedUser(createUser.Email, createUser.Password, createUser.FirstName, createUser.LastName, createUser.DateOfBirth)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while generating unconfirmed user '%s': '%s'", createUser.Email, err.Error()))
		ctx.String(http.StatusBadRequest, "BadRequest")
		return
	}

	err = c.userService.CreateUnconfirmedUser(ctx, *unconfirmedUser)
	if errors.Is(err, constants.ErrDbConflict) {
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
// @Success 302 		{string} 	OKResponse "StatusFound"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/confirm [GET]
func (c *UserHandler) ConfirmUser(ctx *gin.Context) {
	otp := ctx.Query("otp")

	err := c.userService.ConfirmUser(ctx, otp)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while confirming user otp='%s': '%s'", otp, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	// ctx.Header("location", "/")
	ctx.Header("location", constants.AppHostUrl)
	ctx.String(http.StatusFound, "Found")
}

// @Summary GetUserOrgs
// @Tags User
// @Security JWT
// @Description Gets Orgs Users Belongs to
// @Produce json
// @Success 200 		{object} 	[]dto.OrganizationOutput
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/organizations [GET]
func (c *UserHandler) GetUserOrgs(ctx *gin.Context) {
	claims, err := token.GetClaimsFromGinCtx[domain.JwtClaims](ctx)
	if err != nil {
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	orgs, err := c.userService.GetUserOrgs(ctx, claims.UserId)
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
// @Param   payload 	body 		dto.Email true "email json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/init-reset-password [POST]
func (c *UserHandler) InitResetPassword(ctx *gin.Context) {
	var email dto.Email

	if err := ctx.ShouldBind(&email); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	otp, err := common.GenerateRandomString(constants.OptLen)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.userService.GetUser(ctx, email.Email)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	err = c.userService.InitPasswordReset(ctx, user.UserId, otp)
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
// @Success 302 		{string} 	OKResponse "StatusFound"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/set-password-reset-cookie [GET]
func (c *UserHandler) SetPasswordResetCookie(ctx *gin.Context) {
	otp := ctx.Query("otp")

	reset, err := c.userService.GetPasswordReset(ctx, otp)
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

	token.SetCookieForApp(ctx, constants.PasswordResetTimeoutJwtCookieName, tokenStr)

	ctx.Header("location", "/reset-password")
	ctx.String(http.StatusFound, "Found")
}

// @Summary ResetPassword
// @Tags User
// @Description Resets the password (auths via special cookie)
// @Consume application/json
// @Accept json
// @Produce plain
// @Param   payload 	body 		dto.Password true "pw json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/reset-password [POST]
func (c *UserHandler) ResetPassword(ctx *gin.Context) {
	var pw dto.Password

	if err := ctx.ShouldBind(&pw); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	cookieVal, err := ctx.Cookie(constants.PasswordResetTimeoutJwtCookieName)
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

	err = c.userService.UpdateUserPassword(ctx, claims.UserId, pw.Password)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	token.SetCookieForApp(ctx, constants.PasswordResetTimeoutJwtCookieName, "")
	ctx.String(http.StatusOK, "OK")
}

// @Summary EditUser
// @Tags User
// @Description Edits User Info
// @Consume application/json
// @Accept json
// @Produce plain
// @Param   payload 	body 		dto.EditUser true "editUser json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/edit [POST]
func (c *UserHandler) EditUser(ctx *gin.Context) {
	var editUser dto.EditUser

	if err := ctx.ShouldBind(&editUser); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	claims, err := token.GetClaimsFromGinCtx[domain.JwtClaims](ctx)
	if err != nil {
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	err = c.userService.EditUser(ctx, claims.UserId, editUser)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

// @Summary SetPicture
// @Tags User
// @Description Sets User Picture
// @Consume application/json
// @Accept json
// @Produce plain
// @Param   payload 	body 		dto.UloadPicture true "picture json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/profile-picture [PUT]
func (c *UserHandler) SetPicture(ctx *gin.Context) {
	var pic dto.UloadPicture

	if err := ctx.ShouldBind(&pic); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	claims, err := token.GetClaimsFromGinCtx[domain.JwtClaims](ctx)
	if err != nil {
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	picBytes, err := base64.StdEncoding.DecodeString(pic.Content)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	if len(picBytes) > 100*1024 { // if > 100k
		ctx.String(http.StatusBadGateway, "ImgTooLarge")
		return
	}

	imgFmt, err := common.GetImageFormat(picBytes)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	if imgFmt != "png" && imgFmt != "jpeg" {
		ctx.String(http.StatusUnsupportedMediaType, "UnsupportedMediaType")
		return
	}

	objPath := storage.GetPublicPath(storage.UserAvatars, strconv.Itoa(int(claims.UserId)))
	err = c.objService.Upload(ctx, constants.S3Bucket, objPath, int64(len(picBytes)), bytes.NewReader(picBytes))
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	objUrl, err := storage.GetFullObjUrl(objPath)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	err = c.userService.SetAvatarUrl(ctx, claims.UserId, objUrl)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func (c *UserHandler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	g := rg.Group("/users")

	g.PUT("", c.CreateUser)
	g.GET("/confirm", c.ConfirmUser)
	g.POST("/init-reset-password", c.InitResetPassword)
	g.GET("/set-password-reset-cookie", c.SetPasswordResetCookie)
	g.POST("/reset-password", c.ResetPassword)
	g.GET("/organizations", authMiddleware.AuthorizeUser(), c.GetUserOrgs)
	g.POST("/edit", authMiddleware.AuthorizeUser(), c.EditUser)
	g.PUT("/profile-picture", authMiddleware.AuthorizeUser(), c.SetPicture)
}
