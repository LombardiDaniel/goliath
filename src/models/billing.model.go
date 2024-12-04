package models

import "time"

type Billing struct {
	OrderId          uint32     `json:"orderId"`
	UserId           uint32     `json:"userId"`
	Ammount          float32    `json:"ammount"`
	PaymentStatus    string     `json:"paymentStatus"`
	PaymentSessionId string     `json:"paymentSessionId"`
	CreatedAt        time.Time  `json:"createdAt"`
	CompletedAt      *time.Time `json:"completedAt"`
}
