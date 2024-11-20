package controllers

import (
	"log/slog"
	"net/http"
	"time"

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
// @Security JWT
// @Tags Organization
// @Description Creates an Organization
// @Consume application/json
// @Accept json
// @Produce plain
// @Param   payload 	body 		schemas.CreateOrganization true "org json"
// @Success 200 		{object} 	schemas.IdString
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

	orgId = orgId[:5] // cut to 5 (size of field)
	org := models.Organization{
		OrganizationId:   orgId,
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

	ctx.JSON(http.StatusOK, schemas.IdString{Id: orgId})
}

// @Summary InviteToOrg
// @Security JWT
// @Tags Organization
// @Description Invite User to Org
// @Consume application/json
// @Accept json
// @Produce plain
// @Param	orgId 		path string true "Organization Id"
// @Param   payload 	body 		schemas.CreateOrganizationInvite true "invite json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/organizations/{orgId}/invite [PUT]
func (c *OrganizationController) InviteToOrg(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	var createInv schemas.CreateOrganizationInvite

	if err := ctx.ShouldBind(&createInv); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	currUser, err := common.GetClaimsFromGinCtx(ctx)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	otp, err := common.GenerateRandomString(common.OTP_LEN)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	user, err := c.userService.GetUser(rCtx, createInv.UserEmail)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	org, err := c.orgService.GetOrganization(rCtx, *currUser.OrganizationId)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	invExp := time.Now().Add(24 * time.Hour * time.Duration(common.ORG_INVITE_TIMEOUT_DAYS))
	err = c.orgService.CreateOrganizationInvite(rCtx, models.OrganizationInvite{
		OrganizationId: *currUser.OrganizationId,
		UserId:         user.UserId,
		IsAdmin:        createInv.IsAdmin,
		Otp:            &otp,
		Exp:            &invExp,
	})
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	err = c.emailService.SendOrganizationInvite(user.FirstName, user.Email, otp, org.OrganizationName)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

// @Summary AcceptOrgInvite
// @Tags Organization
// @Description Accepts the Organization Invite
// @Consume application/json
// @Accept json
// @Produce plain
// @Param   otp 		query 		string true "OneTimePass sent in email"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/organizations/accept-invite [GET]
func (c *OrganizationController) AcceptOrgInvite(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	otp := ctx.Query("otp")

	// currUser, err := common.GetClaimsFromGinCtx(ctx)
	// if err != nil {
	// 	slog.Error(err.Error())
	// 	ctx.String(http.StatusBadGateway, "BadGateway")
	// 	return
	// }

	err := c.orgService.ConfirmOrganizationInvite(rCtx, otp)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func (c *OrganizationController) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	g := rg.Group("/organizations")

	g.PUT("", authMiddleware.AuthorizeUser(), c.CreateOrganization)
	g.PUT("/:orgId/invite", authMiddleware.AuthorizeOrganization(true), c.InviteToOrg)
	g.GET("/accept-invite", c.AcceptOrgInvite)
}
