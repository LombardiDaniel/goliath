package services

import (
	"context"

	"github.com/LombardiDaniel/gopherbase/models"
	"github.com/LombardiDaniel/gopherbase/schemas"
)

type UserService interface {
	CreateUser(ctx context.Context, user models.User) error
	CreateUnconfirmedUser(ctx context.Context, unconfirmedUser models.UnconfirmedUser) error
	ConfirmUser(ctx context.Context, otp string) error
	GetUser(ctx context.Context, email string) (models.User, error)
	GetUserFromId(ctx context.Context, id uint32) (models.User, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserOrgs(ctx context.Context, userId uint32) ([]schemas.OrganizationOutput, error)
	InitPasswordReset(ctx context.Context, userId uint32, otp string) error
	GetPasswordReset(ctx context.Context, otp string) (models.PasswordReset, error)
	UpdateUserPassword(ctx context.Context, userId uint32, pw string) error

	EditUser(ctx context.Context, userId uint32, user schemas.EditUser) error
	SetAvatarUrl(ctx context.Context, userId uint32, url string) error

	DeleteExpiredPwResets() error
}
