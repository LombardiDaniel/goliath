package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/LombardiDaniel/gopherbase/common"
	"github.com/LombardiDaniel/gopherbase/models"
	"github.com/LombardiDaniel/gopherbase/oauth"
	"github.com/golang-jwt/jwt"
)

type AuthServiceJwtImpl struct {
	jwtSecretKey string
	db           *sql.DB
}

func NewAuthServiceJwtImpl(jwtSecretKey string, db *sql.DB) AuthService {
	return &AuthServiceJwtImpl{
		jwtSecretKey: jwtSecretKey,
		db:           db,
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

func (s *AuthServiceJwtImpl) InitPasswordResetToken(userId uint32) (string, error) {
	claims := models.JwtPasswordResetClaims{
		UserId:  userId,
		Allowed: true,
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
func (s *AuthServiceJwtImpl) ParsePasswordResetToken(tokenString string) (models.JwtPasswordResetClaims, error) {
	claims := models.JwtPasswordResetClaims{}
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

func (s *AuthServiceJwtImpl) LoginOauth(ctx context.Context, oauthUser oauth.User) (models.User, bool, error) {
	user := models.User{}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return user, false, err
	}

	defer tx.Rollback()

	// check if user exists on curr email
	// also creates oauth_users entry for this provider
	err = tx.QueryRowContext(ctx, `
		SELECT
			user_id,
			email,
			password_hash,
			first_name,
			last_name,
			date_of_birth,
			avatar_url,
			created_at,
			updated_at,
			is_active
		FROM users WHERE email = $1;
	`, oauthUser.Email).Scan(
		&user.UserId,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err != nil && err != sql.ErrNoRows {
		return user, false, err
	}

	if err == nil {
		_, err = tx.ExecContext(ctx, `
				INSERT INTO oauth_users (email, user_id, oauth_provider)
				VALUES ($1, $2, $3)
				ON CONFLICT (email, oauth_provider) DO NOTHING;
			`, oauthUser.Email, user.UserId, oauthUser.Provider,
		)
		if err != nil {
			return user, false, err
		}
		return user, false, tx.Commit()
	}

	// here error is sql.ErrNoRows
	err = tx.QueryRowContext(ctx, `
			INSERT INTO users 
				(email, password_hash, first_name, last_name, avatar_url)
			VALUES
				($1, $2, $3, $4, $5)
			RETURNING 
				user_id,
				email,
				password_hash,
				first_name,
				last_name,
				date_of_birth,
				avatar_url,
				created_at,
				updated_at,
				is_active;
		`,
		oauthUser.Email,
		"oauth",
		oauthUser.FirstName,
		oauthUser.LastName,
		oauthUser.PictureUrl,
	).Scan(
		&user.UserId,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err != nil {
		return user, false, err
	}

	_, err = tx.ExecContext(ctx, `
			INSERT INTO oauth_users (email, user_id, oauth_provider)
			VALUES ($1, $2, $3)
			ON CONFLICT (email, oauth_provider) DO NOTHING;
		`, oauthUser.Email, user.UserId, oauthUser.Provider,
	)
	if err != nil {
		return user, false, err
	}
	return user, false, tx.Commit()
}
