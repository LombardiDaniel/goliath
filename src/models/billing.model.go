package models

import "time"

type Payment struct {
	PaymentId               string     `json:"paymentId"`
	UserId                  uint32     `json:"userId"`
	UnitAmmount             uint32     `json:"unitAmmount"`
	UnitCurrency            string     `json:"unitCurrency"`
	PaymentStatus           string     `json:"paymentStatus"`
	StripeCheckoutSessionId string     `json:"stripeCheckoutSessionId"`
	CreatedAt               time.Time  `json:"createdAt"`
	CompletedAt             *time.Time `json:"completedAt"`
}
