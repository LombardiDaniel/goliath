package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/LombardiDaniel/goliath/src/internal/dto"
	"github.com/LombardiDaniel/goliath/src/internal/middlewares"
	"github.com/LombardiDaniel/goliath/src/internal/models"
	"github.com/LombardiDaniel/goliath/src/internal/services"
	"github.com/LombardiDaniel/goliath/src/pkg/token"
	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v81"
)

type BillingHandler struct {
	billingService services.BillingService
	emailService   services.EmailService
	userService    services.UserService
}

func NewBillingHandler(
	billingService services.BillingService,
	emailService services.EmailService,
	userService services.UserService,
) BillingHandler {
	return BillingHandler{
		billingService: billingService,
		emailService:   emailService,
		userService:    userService,
	}
}

// @Summary CheckoutSessionUrl
// @Security JWT
// @Tags Billing
// @Description Gets the CheckoutSession Url
// @Produce plain
// @Param 	product_id 	path 		string true "product_id"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/billing/stripe/get-checkout-session-url/{product_id} [POST]
func (c *BillingHandler) CheckoutSessionUrl(ctx *gin.Context) {
	prodIdStr := ctx.Param("product_id")
	var val int64 = 300

	fmt.Printf("prodIdStr: %v\n", prodIdStr)

	if prodIdStr != "0" {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	claims, err := token.GetClaimsFromGinCtx[models.JwtClaims](ctx)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	url, err := c.billingService.CheckoutURL(ctx, stripe.CurrencyBRL, val*100, "event", claims.UserId)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.JSON(http.StatusOK, dto.Url{Url: url})
}

// @Summary CheckoutSessionCompletedCallback
// @Security JWT
// @Tags Billing
// @Description Completes a SessionCompleted
// @Produce plain
// @Param   payload 	body 		any true "stripe.Event json"
// @Success 200 		{string} 	OKResponse "OK"
// @Failure 400 		{string} 	ErrorResponse "Bad Request"
// @Failure 409 		{string} 	ErrorResponse "Conflict"
// @Failure 502 		{string} 	ErrorResponse "Bad Gateway"
// @Router /v1/billing/stripe/checkout-session-completed [POST]
func (c *BillingHandler) CheckoutSessionCompletedCallback(ctx *gin.Context) {
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
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	checkoutSession, err := c.billingService.CheckoutSession(ctx, inputCheckoutSession.ID)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	if !isStripeChechouseSessionPaid(checkoutSession) {
		ctx.String(http.StatusUnauthorized, "Unauthorized")
		return
	}

	payment, err := c.billingService.SetCheckoutSessionAsComplete(ctx, checkoutSession.ID)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	user, err := c.userService.GetUserFromId(ctx, payment.UserId)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	err = c.emailService.SendPaymentAccepted(user.Email, user.FirstName, payment)
	if err != nil {
		slog.Error(err.Error())
		ctx.String(http.StatusBadGateway, "BadGateway")
		return
	}

	ctx.String(http.StatusOK, "OK")
}

func (c *BillingHandler) RegisterRoutes(rg *gin.RouterGroup, authMiddleware middlewares.AuthMiddleware) {
	g := rg.Group("/billing")

	g.POST("/stripe/get-checkout-session-url/:product_id", authMiddleware.AuthorizeUser(), c.CheckoutSessionUrl)
	g.POST("/stripe/checkout-session-completed", c.CheckoutSessionCompletedCallback)
}
