package domain

import (
	"time"

	"github.com/LombardiDaniel/goliath/src/pkg/common"
	"github.com/LombardiDaniel/goliath/src/pkg/constants"
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

func NewOrganization(orgName string, ownerId uint32) (*Organization, error) {
	orgId, err := common.GenerateRandomString(5)
	if err != nil {
		return nil, err
	}

	return &Organization{
		OrganizationId:   orgId,
		OrganizationName: orgName,
		OwnerUserId:      ownerId,
	}, nil
}

func NewOrganizationInvite(organizationId string, userId uint32, perms map[string]Permission, otp string) OrganizationInvite {
	invExp := time.Now().Add(24 * time.Hour * time.Duration(constants.OrgInviteTimeoutDays))
	return OrganizationInvite{
		OrganizationId: organizationId,
		UserId:         userId,
		Perms:          perms,
		Otp:            &otp,
		Exp:            &invExp,
	}
}
