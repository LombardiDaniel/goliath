package fiddlers

import "github.com/stripe/stripe-go/v81"

func IsStripeChechouseSessionPaid(cs *stripe.CheckoutSession) bool {
	return cs.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid
}
