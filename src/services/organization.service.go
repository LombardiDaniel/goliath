package services

import (
	"context"

	"github.com/LombardiDaniel/gopherbase/models"
)

type OrganizationService interface {
	GetOrganization(ctx context.Context, orgId string) (models.Organization, error)
	CreateOrganization(ctx context.Context, org models.Organization) error
	CreateOrganizationInvite(ctx context.Context, invite models.OrganizationInvite) error
	ConfirmOrganizationInvite(ctx context.Context, otp string) error
	RemoveUserFromOrg(ctx context.Context, orgId string, userId uint32) error
	SetOrganizationOwner(ctx context.Context, orgId string, userId uint32) error
	DeleteExpiredOrgInvites() error
}
