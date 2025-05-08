package handlers

import (
	"github.com/stripe/stripe-go/v81"
)

func isStripeChechouseSessionPaid(cs *stripe.CheckoutSession) bool {
	return cs.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid
}
