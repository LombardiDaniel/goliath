package services

import (
	"context"

	"github.com/LombardiDaniel/gopherbase/models"
	"github.com/LombardiDaniel/gopherbase/oauth"
)

type AuthService interface {
	InitToken(userId uint32, email string, organizationId *string, isAdmin *bool) (string, error)
	ValidateToken(tokenString string) error
	ParseToken(tokenString string) (models.JwtClaims, error)

	InitPasswordResetToken(userId uint32) (string, error)
	ParsePasswordResetToken(tokenString string) (models.JwtPasswordResetClaims, error)

	// LoginOauth logs in the Oauth user, returns bool=true if the user was just created
	// this is to be used in sending welcome email
	LoginOauth(ctx context.Context, oathUser oauth.User) (models.User, bool, error)
}
