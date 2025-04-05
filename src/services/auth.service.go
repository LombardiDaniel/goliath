package services

import (
	"context"

	"github.com/LombardiDaniel/gopherbase/models"
	"github.com/LombardiDaniel/gopherbase/oauth"
)

// AuthService defines the interface for authentication-related operations.
// It provides methods for creating and validating JWTs, parsing tokens,
// handling password reset tokens, and managing OAuth user logins.
type AuthService interface {
	// InitToken generates a new JWT for a user.
	InitToken(userId uint32, email string, organizationId *string, isAdmin *bool) (string, error)

	// ValidateToken checks the validity of a given JWT.
	ValidateToken(tokenString string) error

	// ParseToken extracts claims from a given JWT.
	ParseToken(tokenString string) (models.JwtClaims, error)

	// InitPasswordResetToken generates a special JWT for password reset purposes.
	InitPasswordResetToken(userId uint32) (string, error)

	// ParsePasswordResetToken extracts claims from a password reset JWT.
	ParsePasswordResetToken(tokenString string) (models.JwtPasswordResetClaims, error)

	// LoginOauth logs in an OAuth user and determines if the user was newly created.
	LoginOauth(ctx context.Context, oathUser oauth.User) (models.User, bool, error)
}
