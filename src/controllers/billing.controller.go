package controllers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/LombardiDaniel/gopherbase/middlewares"
	"github.com/LombardiDaniel/gopherbase/schemas"
	"github.com/LombardiDaniel/gopherbase/services"
	"github.com/gin-gonic/gin"
)

type BillingController struct {
	billingService services.BillingService
}

func NewBillingController(
	billingService services.BillingService,
) BillingController {
	return BillingController{
		billingService: billingService,
	}
}

// @Summary CheckOutSessionCompletedCallback
// @Tags Billing
// @Description Completes a SessionCompleted
// @Produce plain
// @Param   payload 	body 		schemas.Id true "stripe.session json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/billing/stripe/checkout-session-completed [POST]
func (c *BillingController) CheckOutSessionCompletedCallback(ctx *gin.Context) {
	rCtx := ctx.Request.Context()
	var stripeSessionId schemas.Id

	if err := ctx.ShouldBind(&stripeSessionId); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	slog.Info(stripeSessionId.Id)
	stripeSession, err := c.billingService.GetCheckoutSession(rCtx, stripeSessionId.Id)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	slog.Info(fmt.Sprintf("%+v", stripeSession))
}

func (c *BillingController) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	g := rg.Group("/billing")

	g.POST("/stripe/checkout-session-completed", c.CheckOutSessionCompletedCallback)
}
