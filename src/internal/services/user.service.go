package services

import (
	"context"

	"github.com/LombardiDaniel/goliath/src/internal/domain"
	"github.com/LombardiDaniel/goliath/src/internal/dto"
)

// UserService defines the interface for user-related operations.
// It provides methods for managing users, handling user confirmations,
// password resets, and user profile updates.
type UserService interface {
	// CreateUser creates a new user.
	CreateUser(ctx context.Context, user domain.User) error

	// CreateUnconfirmedUser creates a new unconfirmed user.
	CreateUnconfirmedUser(ctx context.Context, unconfirmedUser domain.UnconfirmedUser) error

	// ConfirmUser confirms a user using a one-time password (OTP).
	ConfirmUser(ctx context.Context, otp string) error

	// GetUser retrieves a user by their email address.
	GetUser(ctx context.Context, email string) (domain.User, error)

	// GetUserFromId retrieves a user by their ID.
	GetUserFromId(ctx context.Context, id uint32) (domain.User, error)

	// GetUsers retrieves all users.
	GetUsers(ctx context.Context) ([]domain.User, error)

	// GetUserOrgs retrieves the organizations a user belongs to.
	GetUserOrgs(ctx context.Context, userId uint32) ([]dto.OrganizationOutput, error)

	// InitPasswordReset initializes a password reset for a user.
	InitPasswordReset(ctx context.Context, userId uint32, otp string) error

	// GetPasswordReset retrieves a password reset request by its OTP.
	GetPasswordReset(ctx context.Context, otp string) (domain.PasswordReset, error)

	// UpdateUserPassword updates a user's password.
	UpdateUserPassword(ctx context.Context, userId uint32, pw string) error

	// EditUser updates a user's profile information.
	EditUser(ctx context.Context, userId uint32, user dto.EditUser) error

	// SetAvatarUrl sets the avatar URL for a user.
	SetAvatarUrl(ctx context.Context, userId uint32, url string) error

	// DeleteExpiredPwResets deletes all expired password reset requests.
	DeleteExpiredPwResets() error
}
