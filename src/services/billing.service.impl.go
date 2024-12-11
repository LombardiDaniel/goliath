package services

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/LombardiDaniel/gopherbase/common"
	"github.com/LombardiDaniel/gopherbase/models"
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

func (s *BillingServiceStripeImpl) CreatePayment(ctx context.Context, currencyUnit stripe.Currency, unitAmmount int64, planName string, userId uint32) (string, error) {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	var paymentId string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO payments
			(user_id, unit_ammount, unit_currency)
		VALUES
			($1, $2, LOWER($3))
		RETURNING payment_id;
		`,
		userId,
		unitAmmount,
		string(currencyUnit),
	).Scan(&paymentId)
	if err != nil {
		return "", err
	}

	params := &stripe.CheckoutSessionParams{
		ClientReferenceID: stripe.String(paymentId),
		Mode:              stripe.String(string(stripe.CheckoutSessionModePayment)),
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

	_, err = tx.ExecContext(ctx, `
		UPDATE payments
		SET
			stripe_checkout_session_id = $1,
			completed_at = NOW()
		WHERE payment_id = $2;
		`,
		checkout.ID,
		paymentId,
	)
	if err != nil {
		return "", err
	}

	return checkout.URL, tx.Commit()
}

func (s *BillingServiceStripeImpl) GetCheckoutSession(ctx context.Context, sessionId string) (*stripe.CheckoutSession, error) {
	checkout, err := session.Get(sessionId, nil)
	return checkout, err
}

func (s *BillingServiceStripeImpl) SetCheckoutSessionAsComplete(ctx context.Context, sessionId string) (models.Payment, error) {
	var p models.Payment
	err := s.db.QueryRowContext(ctx, `
		UPDATE payments
		SET payment_status = 'complete'
		WHERE stripe_checkout_session_id = $1
		RETURNING *;
		`,
		sessionId,
	).Scan(
		&p.PaymentId,
		&p.UserId,
		&p.UnitAmmount,
		&p.UnitCurrency,
		&p.PaymentStatus,
		&p.StripeCheckoutSessionId,
		&p.CreatedAt,
		&p.CompletedAt,
	)

	return p, err
}

// func (s *BillingServiceStripeImpl) GetClientSecret(ctx context.Context, currencyUnit stripe.Currency, unitAmmount int64, planName string) (string, error) {
// 	panic("not impl")
// 	params := &stripe.CheckoutSessionParams{
// 		Mode:   stripe.String(string(stripe.CheckoutSessionModePayment)),
// 		UIMode: stripe.String("embedded"),
// 		LineItems: []*stripe.CheckoutSessionLineItemParams{
// 			{
// 				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
// 					Currency: stripe.String(string(currencyUnit)),
// 					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
// 						Name: stripe.String(planName),
// 					},
// 					UnitAmount: stripe.Int64(unitAmmount),
// 				},
// 				Quantity: stripe.Int64(1),
// 			},
// 		},
// 		ReturnURL: stripe.String("https://example.com/checkout/return?session_id={CHECKOUT_SESSION_ID}"),
// 	}

// 	checkout, err := session.New(params)
// 	if err != nil {
// 		return "", err
// 	}
// 	return checkout.URL, nil
// }
