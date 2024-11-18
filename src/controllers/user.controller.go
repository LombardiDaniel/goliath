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
	userService  services.UserService
	emailService services.EmailService
}

func NewUserController(
	userService services.UserService,
	emailService services.EmailService,
) UserController {
	return UserController{
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

	otp, err := common.GenerateRandomString(128)
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

	err = c.emailService.SendAccountConfirmation(unconfirmedUser.FirstName+" "+unconfirmedUser.LastName, unconfirmedUser.Email, unconfirmedUser.Otp)
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
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users/confirm [POST]
func (c *UserController) ConfirmUser(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	otp := ctx.Query("otp")

	err := c.userService.ConfirmUser(rCtx, otp)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while confirming user otp='%s': '%s'", otp, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func (c *UserController) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddlewareJWT) {
	r := rg.Group("/users")

	r.PUT("", c.CreateUser)
	r.POST("/confirm", c.ConfirmUser)
}
