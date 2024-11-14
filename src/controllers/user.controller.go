package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/LombardiDaniel/go-gin-template/middlewares"
	"github.com/LombardiDaniel/go-gin-template/models"
	"github.com/LombardiDaniel/go-gin-template/services"
	"github.com/LombardiDaniel/go-gin-template/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(
	userService services.UserService,
) UserController {
	return UserController{
		userService: userService,
	}
}

// @Summary CreateUser
// @Tags User
// @Description Creates an User
// @Consume application/json
// @Accept json
// @Produce plain
// @Param   payload 	body 		models.User true "user json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/users [PUT]
func (c *UserController) CreateUser(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	var user models.User

	if err := ctx.ShouldBind(&user); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while hashing pass '%s': '%s'", user.Password, err.Error()))
		ctx.String(http.StatusBadRequest, "BadRequest")
		return
	}

	user.Password = hash

	err = c.userService.CreateUser(rCtx, user)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while creating user '%s': '%s'", user.Username, err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func (c *UserController) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	r := rg.Group("/users")

	r.PUT("/", authMiddleware.Authorize(), c.CreateUser)
	// r.GET("/", authMiddleware.Authorize(), c.GetUsers)
}
