package controllers

import (
	"log/slog"
	"net/http"

	"github.com/LombardiDaniel/go-gin-template/common"
	"github.com/LombardiDaniel/go-gin-template/middlewares"
	"github.com/LombardiDaniel/go-gin-template/models"
	"github.com/LombardiDaniel/go-gin-template/schemas"
	"github.com/LombardiDaniel/go-gin-template/services"
	"github.com/gin-gonic/gin"
)

type OrganizationController struct {
	userService  services.UserService
	emailService services.EmailService
	orgService   services.OrganizationService
}

func NewOrganizationController(
	userService services.UserService,
	emailService services.EmailService,
	orgService services.OrganizationService,
) OrganizationController {
	return OrganizationController{
		userService:  userService,
		emailService: emailService,
		orgService:   orgService,
	}
}

// @Summary CreateOrganization
// @Tags Organization
// @Description Creates an Organization
// @Consume application/json
// @Accept json
// @Produce plain
// @Param   payload 	body 		schemas.CreateOrganization true "org json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/organizations [PUT]
func (c *OrganizationController) CreateOrganization(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	var createOrg schemas.CreateOrganization

	if err := ctx.ShouldBind(&createOrg); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := common.GetClaimsFromGinCtx(ctx)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	orgId, err := common.GenerateRandomString(10)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}
	org := models.Organization{
		OrganizationId:   orgId[:5],
		OrganizationName: createOrg.OrganizationName,
		OwnerUserId:      &user.UserId,
	}

	err = c.orgService.CreateOrganization(rCtx, org)
	if err != nil {
		if err == common.ErrDbConflict {
			ctx.String(http.StatusConflict, "Conflict")
			return
		}
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func (c *OrganizationController) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	r := rg.Group("/organizations")

	r.PUT("", authMiddleware.Authorize(), c.CreateOrganization)
}
