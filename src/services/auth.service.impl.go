package services

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/LombardiDaniel/go-gin-template/common"
	"github.com/LombardiDaniel/go-gin-template/models"
	"github.com/golang-jwt/jwt"
)

type AuthServiceJwtImpl struct {
	jwtSecretKey string
}

func NewAuthServiceJwtImpl(jwtSecretKey string) AuthService {
	return &AuthServiceJwtImpl{
		jwtSecretKey: jwtSecretKey,
	}
}

func (s *AuthServiceJwtImpl) InitToken(userId uint32, email string, organizationId *string, isAdmin *bool) (string, error) {
	claims := models.JwtClaims{
		UserId:         userId,
		Email:          email,
		OrganizationId: organizationId,
		IsAdmin:        isAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(common.JWT_TIMEOUT_SECS)).Unix(),
			Issuer:    common.PROJECT_NAME + "-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthServiceJwtImpl) ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

func (s *AuthServiceJwtImpl) ParseToken(tokenString string) (models.JwtClaims, error) {
	claims := models.JwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecretKey), nil
	})

	if err != nil {
		return claims, err
	}

	slog.Debug(fmt.Sprintf("%+v", claims))
	slog.Debug(fmt.Sprintf("%+v", token.Valid))

	if !token.Valid {
		return claims, errors.New("invalid token")
	}

	return claims, nil
}
