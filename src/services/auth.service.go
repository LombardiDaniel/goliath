package services

import (
	"github.com/LombardiDaniel/go-gin-template/models"
)

type AuthService interface {
	InitToken(userId uint32, email string, organizationId *string, isAdmin *bool) (string, error)
	ValidateToken(tokenString string) error
	ParseToken(tokenString string) (models.JwtClaims, error)

	InitPasswordResetToken(userId uint32) (string, error)
	ParsePasswordResetToken(tokenString string) (models.JwtPasswordResetClaims, error)
}
