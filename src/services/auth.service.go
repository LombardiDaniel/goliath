package services

import (
	"context"

	"github.com/LombardiDaniel/gopherbase/models"
	"github.com/LombardiDaniel/gopherbase/oauth"
)

type AuthService interface {
	// Creates a new JWT
	InitToken(userId uint32, email string, organizationId *string, isAdmin *bool) (string, error)

	// Validates JWT, returns error if it is not valid
	ValidateToken(tokenString string) error

	// Parses the JWT to a Claims struct
	ParseToken(tokenString string) (models.JwtClaims, error)

	// Creates a special password-reset-JWT
	InitPasswordResetToken(userId uint32) (string, error)

	// Parses the special password-reset-JWT to its claims Struct
	ParsePasswordResetToken(tokenString string) (models.JwtPasswordResetClaims, error)

	// LoginOauth logs in the Oauth user, returns bool=true if the user was just created
	// this is to be used in sending welcome email
	LoginOauth(ctx context.Context, oathUser oauth.User) (models.User, bool, error)
}
