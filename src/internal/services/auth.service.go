package services

import (
	"context"

	"github.com/LombardiDaniel/goliath/src/internal/domain"
	"github.com/LombardiDaniel/goliath/src/pkg/oauth"
)

// AuthService defines the interface for authentication-related operations.
// It provides methods for creating and validating JWTs, parsing tokens,
// handling password reset tokens, and managing OAuth user logins.
type AuthService interface {
	// InitToken generates a new JWT for a user.
	InitToken(ctx context.Context, userId uint32, email string, organizationId *string) (string, error)

	// Permissions retrieves the permissions for user in organization.
	Permissions(ctx context.Context, userId uint32, organizationId *string) (map[string]domain.Permission, error)

	// ValidateToken checks the validity of a given JWT.
	ValidateToken(tokenString string) error

	// ParseToken extracts claims from a given JWT.
	ParseToken(tokenString string) (domain.JwtClaims, error)

	// InitPasswordResetToken generates a special JWT for password reset purposes.
	InitPasswordResetToken(userId uint32) (string, error)

	// ParsePasswordResetToken extracts claims from a password reset JWT.
	ParsePasswordResetToken(tokenString string) (domain.JwtPasswordResetClaims, error)

	// LoginOauth logs in an OAuth user and determines if the user was newly created.
	LoginOauth(ctx context.Context, oathUser oauth.User) (domain.User, bool, error)
}
