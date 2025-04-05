package services

import (
	"context"

	"github.com/LombardiDaniel/gopherbase/models"
	"github.com/LombardiDaniel/gopherbase/schemas"
)

// UserService defines the interface for user-related operations.
// It provides methods for managing users, handling user confirmations,
// password resets, and user profile updates.
type UserService interface {
	// CreateUser creates a new user.
	CreateUser(ctx context.Context, user models.User) error

	// CreateUnconfirmedUser creates a new unconfirmed user.
	CreateUnconfirmedUser(ctx context.Context, unconfirmedUser models.UnconfirmedUser) error

	// ConfirmUser confirms a user using a one-time password (OTP).
	ConfirmUser(ctx context.Context, otp string) error

	// GetUser retrieves a user by their email address.
	GetUser(ctx context.Context, email string) (models.User, error)

	// GetUserFromId retrieves a user by their ID.
	GetUserFromId(ctx context.Context, id uint32) (models.User, error)

	// GetUsers retrieves all users.
	GetUsers(ctx context.Context) ([]models.User, error)

	// GetUserOrgs retrieves the organizations a user belongs to.
	GetUserOrgs(ctx context.Context, userId uint32) ([]schemas.OrganizationOutput, error)

	// InitPasswordReset initializes a password reset for a user.
	InitPasswordReset(ctx context.Context, userId uint32, otp string) error

	// GetPasswordReset retrieves a password reset request by its OTP.
	GetPasswordReset(ctx context.Context, otp string) (models.PasswordReset, error)

	// UpdateUserPassword updates a user's password.
	UpdateUserPassword(ctx context.Context, userId uint32, pw string) error

	// EditUser updates a user's profile information.
	EditUser(ctx context.Context, userId uint32, user schemas.EditUser) error

	// SetAvatarUrl sets the avatar URL for a user.
	SetAvatarUrl(ctx context.Context, userId uint32, url string) error

	// DeleteExpiredPwResets deletes all expired password reset requests.
	DeleteExpiredPwResets() error
}
