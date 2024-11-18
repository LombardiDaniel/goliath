package services

import (
	"github.com/LombardiDaniel/go-gin-template/models"
)

type AuthService interface {
	// InitToken(user models.User) (string, error)
	InitToken(userId uint32, email string, organizationId *string) (string, error)
	ValidateToken(tokenString string) error
	ParseToken(tokenString string) (models.JwtClaims, error)
}
