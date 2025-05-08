package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/LombardiDaniel/goliath/src/internal/domain"
	"github.com/LombardiDaniel/goliath/src/pkg/constants"
	"github.com/LombardiDaniel/goliath/src/pkg/oauth"
	"github.com/LombardiDaniel/goliath/src/pkg/validators"
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

func (s *AuthServiceJwtImpl) InitToken(ctx context.Context, userId uint32, email string, organizationId *string) (string, error) {
	perms, err := s.Permissions(ctx, userId, organizationId)
	if err != nil {
		return "", err
	}

	claims := domain.JwtClaims{
		UserId:         userId,
		Email:          email,
		OrganizationId: organizationId,
		Perms:          perms,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(constants.JwtTimeoutSecs)).Unix(),
			Issuer:    constants.ProjectName + "-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return "", errors.Join(err, errors.New("could not hash jwt"))
	}

	return tokenString, nil
}

func (s *AuthServiceJwtImpl) Permissions(ctx context.Context, userId uint32, organizationId *string) (map[string]domain.Permission, error) {
	q := `
		SELECT
			action_name,
			permission
		FROM organization_user_permissions
		WHERE
			user_id = $1 AND
			organization_id = $2;
	`
	rows, err := s.db.QueryContext(ctx, q, userId, organizationId)
	if err != nil && !errors.Is(validators.FilterSqlPgError(err), constants.ErrNoRows) {
		return nil, err
	}

	actionPerms := make(map[string]domain.Permission)
	for rows.Next() {
		var actionName string
		var perm domain.Permission
		if err := rows.Scan(&actionName, &perm); err != nil {
			return nil, err
		}

		actionPerms[actionName] = perm
	}

	return actionPerms, nil
}

func (s *AuthServiceJwtImpl) ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return s.jwtSecretKey, nil
	})
	if err != nil {
		return errors.Join(err, errors.New("could not parse token"))
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

func (s *AuthServiceJwtImpl) ParseToken(tokenString string) (domain.JwtClaims, error) {
	claims := domain.JwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(s.jwtSecretKey), nil
	})
	if err != nil {
		return claims, errors.Join(err, errors.New("could not parse token to claims"))
	}

	// slog.Debug(fmt.Sprintf("%+v", claims))
	// slog.Debug(fmt.Sprintf("%+v", token.Valid))

	if !token.Valid {
		return claims, errors.New("invalid token")
	}

	return claims, nil
}

func (s *AuthServiceJwtImpl) InitPasswordResetToken(userId uint32) (string, error) {
	claims := domain.JwtPasswordResetClaims{
		UserId:  userId,
		Allowed: true,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * time.Duration(constants.JwtTimeoutSecs)).Unix(),
			Issuer:    constants.ProjectName + "-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func (s *AuthServiceJwtImpl) ParsePasswordResetToken(tokenString string) (domain.JwtPasswordResetClaims, error) {
	claims := domain.JwtPasswordResetClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		return []byte(s.jwtSecretKey), nil
	})
	if err != nil {
		return claims, errors.Join(err, errors.New("could not parse token to claims"))
	}

	// slog.Debug(fmt.Sprintf("%+v", claims))
	// slog.Debug(fmt.Sprintf("%+v", token.Valid))

	if !token.Valid {
		return claims, errors.New("invalid token")
	}

	return claims, nil
}

func (s *AuthServiceJwtImpl) LoginOauth(ctx context.Context, oauthUser oauth.User) (domain.User, bool, error) {
	user := domain.User{}
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
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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
