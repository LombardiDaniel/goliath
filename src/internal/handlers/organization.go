package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LombardiDaniel/goliath/src/internal/dto"
	"github.com/LombardiDaniel/goliath/src/internal/middlewares"
	"github.com/LombardiDaniel/goliath/src/internal/models"
	"github.com/LombardiDaniel/goliath/src/internal/services"
	"github.com/LombardiDaniel/goliath/src/pkg/common"
	"github.com/LombardiDaniel/goliath/src/pkg/constants"
	"github.com/LombardiDaniel/goliath/src/pkg/token"
	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	userService  services.UserService
	emailService services.EmailService
	orgService   services.OrganizationService
}

func NewOrganizationHandler(
	userService services.UserService,
	emailService services.EmailService,
	orgService services.OrganizationService,
) OrganizationHandler {
	return OrganizationHandler{
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
// @Param   payload 	body 		dto.CreateOrganization true "org json"
// @Success 200 		{object} 	dto.Id
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/organizations [POST]
func (c *OrganizationHandler) CreateOrganization(ctx *gin.Context) {
	var createOrg dto.CreateOrganization

	if err := ctx.ShouldBind(&createOrg); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	user, err := token.GetClaimsFromGinCtx[models.JwtClaims](ctx)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	org, err := models.NewOrganization(createOrg.OrganizationName, user.UserId)
	if err != nil {
		slog.Error(fmt.Sprintf("Error while generating organization: %s", err.Error()))
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	err = c.orgService.CreateOrganization(ctx, *org)
	if err != nil {
		if errors.Is(err, constants.ErrDbConflict) {
			ctx.String(http.StatusConflict, "Conflict")
			return
		}
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.JSON(http.StatusOK, dto.Id{Id: org.OrganizationId})
}

// @Summary InviteToOrg
// @Security JWT
// @Tags Organization
// @Description Invite User to Org
// @Consume application/json
// @Accept json
// @Produce plain
// @Param	orgId 		path string true "Organization Id"
// @Param   payload 	body 		dto.CreateOrganizationInvite true "invite json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/organizations/{orgId}/invite [POST]
func (c *OrganizationHandler) InviteToOrg(ctx *gin.Context) {
	var createInv dto.CreateOrganizationInvite

	if err := ctx.ShouldBind(&createInv); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	currUser, err := token.GetClaimsFromGinCtx[models.JwtClaims](ctx)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	otp, err := common.GenerateRandomString(constants.OptLen)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	user, err := c.userService.GetUser(ctx, createInv.UserEmail)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	org, err := c.orgService.GetOrganization(ctx, *currUser.OrganizationId)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	inv := models.NewOrganizationInvite(*currUser.OrganizationId, user.UserId, createInv.Perms, otp)
	err = c.orgService.CreateOrganizationInvite(ctx, inv)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	err = c.emailService.SendOrganizationInvite(user.Email, user.FirstName, otp, org.OrganizationName)
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
func (c *OrganizationHandler) AcceptOrgInvite(ctx *gin.Context) {

	otp := ctx.Query("otp")

	err := c.orgService.ConfirmOrganizationInvite(ctx, otp)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

// @Summary RemoveFromOrg
// @Security JWT
// @Tags Organization
// @Description Removes User from Org
// @Produce plain
// @Param	orgId 		path string true "Organization Id"
// @Param	userId 		path string true "User Id"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/organizations/{orgId}/users/{userId} [DELETE]
func (c *OrganizationHandler) RemoveFromOrg(ctx *gin.Context) {

	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		ctx.String(http.StatusBadRequest, "BadRequest")
		return
	}

	var createInv dto.CreateOrganizationInvite

	if err := ctx.ShouldBind(&createInv); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	currUser, err := token.GetClaimsFromGinCtx[models.JwtClaims](ctx)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	err = c.orgService.RemoveUserFromOrg(ctx, *currUser.OrganizationId, uint32(userId))
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

// @Summary ChangeOwner
// @Security JWT
// @Tags Organization
// @Description Removes User from Org
// @Produce plain
// @Param	orgId 		path string true "Organization Id"
// @Param   payload 	body 		dto.Email true "email json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/organizations/{orgId}/owner [PUT]
func (c *OrganizationHandler) ChangeOwner(ctx *gin.Context) {

	var tgtEmail dto.Email

	if err := ctx.ShouldBind(&tgtEmail); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	currUser, err := token.GetClaimsFromGinCtx[models.JwtClaims](ctx)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	org, err := c.orgService.GetOrganization(ctx, *currUser.OrganizationId)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	if currUser.UserId != org.OwnerUserId {
		ctx.String(http.StatusUnauthorized, "StatusUnauthorized")
		return
	}

	tgtUser, err := c.userService.GetUser(ctx, tgtEmail.Email)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	err = c.orgService.SetOrganizationOwner(ctx, *currUser.OrganizationId, tgtUser.UserId)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func (c *OrganizationHandler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	g := rg.Group("/organizations")

	adminPerms := map[string]models.Permission{
		"admin": models.ReadWritePermission,
		// "PUT:orgInvite": models.ReadWritePermission,
	}

	ownerPerms := map[string]models.Permission{
		"owner": models.ReadWritePermission,
	}

	g.POST("", authMiddleware.AuthorizeUser(), c.CreateOrganization)
	g.POST("/:orgId/invite", authMiddleware.AuthorizeOrganization(adminPerms), c.InviteToOrg)
	g.PUT("/:orgId/owner", authMiddleware.AuthorizeOrganization(ownerPerms), c.ChangeOwner, authMiddleware.Reauthorize())
	g.GET("/accept-invite", c.AcceptOrgInvite)
	g.DELETE("/:orgId/users/:userId", authMiddleware.AuthorizeOrganization(adminPerms), c.RemoveFromOrg)
}
