package models

import (
	"time"
)

type Organization struct {
	OrganizationId   string     `json:"organizationId" binding:"required,min=1"`
	OrganizationName string     `json:"organizationName" binding:"required,min=1"`
	BillingPlanId    *uint32    `json:"billingPlanId" binding:"required"`
	CreatedAt        time.Time  `json:"createdAt" binding:"required"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty"`
	OwnerUserId      uint32     `json:"ownerUserId,omitempty"`
}

type FrontendConfig struct {
	OrganizationId string `json:"organizationId"`
	PrimaryColor   string `json:"primaryColor"`
	SecondaryColor string `json:"secondaryColor"`
}

type OrganizationInvite struct {
	OrganizationId string     `json:"organizationId" binding:"required,min=1"`
	UserId         uint32     `json:"userId" binding:"required"`
	IsAdmin        bool       `json:"isAdmin" biding:"required"`
	Otp            *string    `json:"otp,omitempty"`
	Exp            *time.Time `json:"exp,omitempty"`
}
