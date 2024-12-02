package services

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/LombardiDaniel/gopherbase/common"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

type BillingServiceStripeImpl struct {
	db            *sql.DB
	appSuccessUrl string
	appCancelUrl  string
}

// https://www.youtube.com/watch?v=M4aCgy67f243
// https://www.youtube.com/playlist?list=PLy1nL-pvL2M5eqpSBR9KL7K0lcnWo0V0a
// https://docs.stripe.com/payments/accept-a-payment?platform=web&ui=stripe-hosted
// https://www.youtube.com/watch?v=ePmEVBu8w6Y

func NewBillingService(db *sql.DB, stripeApiKey string) BillingService {
	successUrl, err := url.JoinPath(common.APP_HOST_URL, "/billing/success")
	if err != nil {
		panic(err)
	}

	cancelUrl, err := url.JoinPath(common.APP_HOST_URL, "/billing/cancel")
	if err != nil {
		panic(err)
	}

	stripe.Key = stripeApiKey

	return &BillingServiceStripeImpl{
		db:            db,
		appSuccessUrl: successUrl,
		appCancelUrl:  cancelUrl,
	}
}

func (s *BillingServiceStripeImpl) GetCheckoutUrl(ctx context.Context, currencyUnit stripe.Currency, unitAmmount int64, planName string) (string, error) {
	params := &stripe.CheckoutSessionParams{
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(currencyUnit)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(planName),
					},
					UnitAmount: stripe.Int64(unitAmmount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: &s.appSuccessUrl,
		CancelURL:  &s.appCancelUrl,
	}

	checkout, err := session.New(params)
	if err != nil {
		return "", err
	}
	return checkout.URL, nil
}

func (s *BillingServiceStripeImpl) GetClientSecret(ctx context.Context, currencyUnit stripe.Currency, unitAmmount int64, planName string) (string, error) {
	panic("not impl")
	params := &stripe.CheckoutSessionParams{
		Mode:   stripe.String(string(stripe.CheckoutSessionModePayment)),
		UIMode: stripe.String("embedded"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(currencyUnit)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(planName),
					},
					UnitAmount: stripe.Int64(unitAmmount),
				},
				Quantity: stripe.Int64(1),
			},
		},
		ReturnURL: stripe.String("https://example.com/checkout/return?session_id={CHECKOUT_SESSION_ID}"),
	}

	checkout, err := session.New(params)
	if err != nil {
		return "", err
	}
	return checkout.URL, nil
}

func (s *BillingServiceStripeImpl) GetCheckoutSession(ctx context.Context, sessionId string) (*stripe.CheckoutSession, error) {
	checkout, err := session.Get(sessionId, nil)
	return checkout, err
}
