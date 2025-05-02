package services

import (
	"context"

	"github.com/LombardiDaniel/gopherbase/models"
)

// OrganizationService defines the interface for organization-related operations.
// It provides methods for managing organizations, handling organization invites,
// and performing administrative tasks such as setting owners and removing users.
type OrganizationService interface {
	// GetOrganization retrieves an organization by its ID.
	GetOrganization(ctx context.Context, orgId string) (models.Organization, error)

	// CreateOrganization creates a new organization.
	CreateOrganization(ctx context.Context, org models.Organization) error

	// CreateOrganizationInvite creates an invitation for a user to join an organization.
	CreateOrganizationInvite(ctx context.Context, invite models.OrganizationInvite) error

	// ConfirmOrganizationInvite confirms an organization invite using a one-time password (OTP).
	ConfirmOrganizationInvite(ctx context.Context, otp string) error

	// RemoveUserFromOrg removes a user from an organization by their user ID.
	RemoveUserFromOrg(ctx context.Context, orgId string, userId uint32) error

	// SetOrganizationOwner sets a user as the owner of an organization.
	SetOrganizationOwner(ctx context.Context, orgId string, userId uint32) error

	// DeleteExpiredOrgInvites deletes all expired organization invites.
	DeleteExpiredOrgInvites() error

	// SetPerms changes the permission of a user
	SetPerms(ctx context.Context, action string, userId uint32, perms models.Permission) error
}
