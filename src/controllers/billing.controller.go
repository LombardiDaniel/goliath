package controllers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/LombardiDaniel/gopherbase/common"
	"github.com/LombardiDaniel/gopherbase/middlewares"
	"github.com/LombardiDaniel/gopherbase/schemas"
	"github.com/LombardiDaniel/gopherbase/services"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
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

// @Summary GetCheckoutSessionUrl
// @Tags Billing
// @Description Gets the CheckoutSession Url
// @Produce plain
// @Param 	value 	path 		string true "value"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/billing/stripe/get-checkout-session-url/{value} [POST]
func (c *BillingController) GetCheckoutSessionUrl(ctx *gin.Context) {
	valStr := ctx.Param("value")
	if valStr != "300" {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	claims, err := common.GetClaimsFromGinCtx(ctx)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	url, err := c.billingService.CreateOrder(ctx, stripe.CurrencyBRL, 300*100, "event", claims.UserId)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.JSON(http.StatusOK, schemas.Url{Url: url})
}

// @Summary CheckOutSessionCompletedCallback
// @Tags Billing
// @Description Completes a SessionCompleted
// @Produce plain
// @Param   payload 	body 		any true "stripe.Event json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/billing/stripe/checkout-session-completed [POST]
func (c *BillingController) CheckOutSessionCompletedCallback(ctx *gin.Context) {
	var stripeEvent stripe.Event

	if err := ctx.ShouldBind(&stripeEvent); err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}

	var inputCheckoutSession stripe.CheckoutSession

	switch stripeEvent.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		err := json.Unmarshal(stripeEvent.Data.Raw, &inputCheckoutSession)
		if err != nil {
			slog.Error(err.Error())
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
	default:
		slog.Warn(fmt.Sprintf("Unhandled event type: %s", string(stripeEvent.Type)))
	}

	checkoutSession, err := c.billingService.GetCheckoutSession(ctx, inputCheckoutSession.ID)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	slog.Info(fmt.Sprintf("%+v", checkoutSession))
}

func (c *BillingController) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	g := rg.Group("/billing")

	g.POST("/stripe/get-checkout-session-url/:value", authMiddleware.AuthorizeUser(), c.GetCheckoutSessionUrl)
	g.POST("/stripe/checkout-session-completed", c.CheckOutSessionCompletedCallback)
}
