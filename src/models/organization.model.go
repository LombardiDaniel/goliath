package models

import (
	"time"
)

// Organization represents an organization in the system.
type Organization struct {
	OrganizationId   string     `json:"organizationId" binding:"required,min=1"`
	OrganizationName string     `json:"organizationName" binding:"required,min=1"`
	BillingPlanId    *uint32    `json:"billingPlanId" binding:"required"`
	CreatedAt        time.Time  `json:"createdAt" binding:"required"`
	DeletedAt        *time.Time `json:"deletedAt,omitempty"`
	OwnerUserId      uint32     `json:"ownerUserId,omitempty"`
}

// FrontendConfig represents the frontend configuration for an organization.
type FrontendConfig struct {
	OrganizationId string `json:"organizationId"`
	PrimaryColor   string `json:"primaryColor"`
	SecondaryColor string `json:"secondaryColor"`
}

// OrganizationInvite represents an invitation to join an organization.
type OrganizationInvite struct {
	OrganizationId string                `json:"organizationId" binding:"required,min=1"`
	UserId         uint32                `json:"userId" binding:"required"`
	Perms          map[string]Permission `json:"perms" binding:"required"`
	Otp            *string               `json:"otp,omitempty"`
	Exp            *time.Time            `json:"exp,omitempty"`
}
