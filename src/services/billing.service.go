package services

import (
	"context"

	"github.com/LombardiDaniel/gopherbase/models"
	"github.com/stripe/stripe-go/v81"
)

type BillingService interface {
	// Gets the Stripe Checkout URL to be redirected to in the frontend
	CreatePayment(ctx context.Context, currencyUnit stripe.Currency, unitAmmount int64, planName string, userId uint32) (string, error)

	// Webhook to be used in a daemon
	GetCheckoutSession(ctx context.Context, sessionId string) (*stripe.CheckoutSession, error)

	// Webhook to be used in a daemon
	SetCheckoutSessionAsComplete(ctx context.Context, sessionId string) (models.Payment, error)

	// // Gets the Stripe Client Secret to be used in Embedded Checkout Form
	// 	//
	// 	// Take a look at: https://docs.stripe.com/payments/accept-a-payment?platform=web&ui=stripe-hosted
	// 	GetClientSecret(ctx context.Context, currencyUnit stripe.Currency, unitAmmount int64, planName string) (string, error)
}
