package services

import (
	"context"

	"github.com/LombardiDaniel/goliath/src/internal/domain"
	"github.com/stripe/stripe-go/v81"
)

// BillingService defines the interface for all billing-related operations.
type BillingService interface {
	// CheckoutURL gets the Stripe Checkout URL to be redirected to in the frontend
	CheckoutURL(ctx context.Context, currencyUnit stripe.Currency, unitAmmount int64, planName string, userId uint32) (string, error)

	// CheckoutSession is the webhook to be used in a daemon
	CheckoutSession(ctx context.Context, sessionId string) (*stripe.CheckoutSession, error)

	// SetCheckoutSessionAsComplete is the webhook to be used in a daemon
	SetCheckoutSessionAsComplete(ctx context.Context, sessionId string) (domain.Payment, error)

	// // Gets the Stripe Client Secret to be used in Embedded Checkout Form
	// 	//
	// 	// Take a look at: https://docs.stripe.com/payments/accept-a-payment?platform=web&ui=stripe-hosted
	// 	GetClientSecret(ctx context.Context, currencyUnit stripe.Currency, unitAmmount int64, planName string) (string, error)
}
