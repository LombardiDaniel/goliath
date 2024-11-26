package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/LombardiDaniel/go-gin-template/common"
	"github.com/LombardiDaniel/go-gin-template/models"
	"github.com/LombardiDaniel/go-gin-template/oauth"
	"github.com/LombardiDaniel/go-gin-template/schemas"
)

type UserServicePgImpl struct {
	db *sql.DB
}

func NewUserServicePgImpl(db *sql.DB) UserService {
	return &UserServicePgImpl{
		db: db,
	}
}

func (s *UserServicePgImpl) CreateUser(ctx context.Context, user models.User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth)
		VALUES ($1, $2, $3, $4, $5);
	`

	err := s.db.QueryRowContext(ctx, query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
	).Err()

	if err != nil {
		return common.FilterSqlPgError(err)
	}

	return nil
}

func (s *UserServicePgImpl) CreateUnconfirmedUser(ctx context.Context, unconfirmedUser models.UnconfirmedUser) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var userId uint32
	err = tx.QueryRowContext(ctx, `
			SELECT user_id
			FROM users WHERE email = $1;
		`,
		unconfirmedUser.Email,
	).Scan(&userId)
	if err != nil && err != sql.ErrNoRows {
		return common.FilterSqlPgError(err)
	}

	err = tx.QueryRowContext(ctx, `
			INSERT INTO unconfirmed_users (email, otp, password_hash, first_name, last_name, date_of_birth)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (email) DO UPDATE
			SET 
				otp = EXCLUDED.otp,
				password_hash = EXCLUDED.password_hash,
				first_name = EXCLUDED.first_name,
				last_name = EXCLUDED.last_name,
				date_of_birth = EXCLUDED.date_of_birth;
		`,
		unconfirmedUser.Email,
		unconfirmedUser.Otp,
		unconfirmedUser.PasswordHash,
		unconfirmedUser.FirstName,
		unconfirmedUser.LastName,
		unconfirmedUser.DateOfBirth,
	).Err()

	if err != nil {
		return common.FilterSqlPgError(err)
	}

	return tx.Commit()
}

func (s *UserServicePgImpl) ConfirmUser(ctx context.Context, otp string) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	unconfirmedUser := models.UnconfirmedUser{}
	err = s.db.QueryRowContext(ctx, `
			SELECT
				email,
				otp,
				password_hash,
				first_name,
				last_name,
				date_of_birth
			FROM
				unconfirmed_users WHERE otp = $1
		`, otp).Scan(
		&unconfirmedUser.Email,
		&unconfirmedUser.Otp,
		&unconfirmedUser.PasswordHash,
		&unconfirmedUser.FirstName,
		&unconfirmedUser.LastName,
		&unconfirmedUser.DateOfBirth,
	)
	if err != nil {
		return common.FilterSqlPgError(err)
	}

	_, err = tx.ExecContext(ctx, `
			INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth)
			VALUES ($1, $2, $3, $4, $5);
		`,
		unconfirmedUser.Email,
		unconfirmedUser.PasswordHash,
		unconfirmedUser.FirstName,
		unconfirmedUser.LastName,
		unconfirmedUser.DateOfBirth,
	)
	if err != nil {
		return common.FilterSqlPgError(err)
	}

	_, err = tx.ExecContext(ctx, `
			DELETE FROM unconfirmed_users WHERE otp = $1;
		`,
		unconfirmedUser.Otp,
	)
	if err != nil {
		return common.FilterSqlPgError(err)
	}

	return tx.Commit()
}

func (s *UserServicePgImpl) GetUser(ctx context.Context, email string) (models.User, error) {
	query := `
		SELECT
			user_id,
			email,
			password_hash,
			first_name,
			last_name,
			date_of_birth,
			created_at,
			updated_at,
			is_active
		FROM users WHERE email = $1
	`

	user := models.User{}

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.UserId,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		// &user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err != nil {
		return user, common.FilterSqlPgError(err)
	}

	return user, nil
}

func (s *UserServicePgImpl) GetUserFromId(ctx context.Context, id uint32) (models.User, error) {
	query := `
		SELECT
			user_id,
			email,
			password_hash,
			first_name,
			last_name,
			date_of_birth,
			created_at,
			updated_at,
			is_active
		FROM users WHERE user_id = $1
	`

	user := models.User{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.UserId,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		// &user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err != nil {
		return user, common.FilterSqlPgError(err)
	}

	return user, nil
}

// func (s *UserServicePgImpl) GetUsers(ctx context.Context) ([]models.User, error) {
// 	return nil, errors.New("not implemented")
// }

func (s *UserServicePgImpl) GetUserOrgs(ctx context.Context, userId uint32) ([]schemas.OrganizationOutput, error) {
	query := `
		SELECT DISTINCT
			o.organization_id,
			o.organization_name,
			ou.is_admin
		FROM
			organizations o
		INNER JOIN
			organizations_users ou ON o.organization_id = ou.organization_id
		WHERE ou.user_id = $1;
	`

	orgs := []schemas.OrganizationOutput{}

	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return orgs, common.FilterSqlPgError(err)
	}
	defer rows.Close()

	for rows.Next() {
		newOrg := schemas.OrganizationOutput{}
		err := rows.Scan(&newOrg.OrganizationId, &newOrg.OrganizationName, &newOrg.IsAdmin)
		if err != nil {
			return orgs, err
		}
		orgs = append(orgs, newOrg)
	}

	return orgs, nil
}

func (s *UserServicePgImpl) InitPasswordReset(ctx context.Context, userId uint32, otp string) error {
	query := `
		INSERT INTO password_resets (user_id, otp, exp)
		VALUES ($1, $2, $3);
	`

	_, err := s.db.ExecContext(ctx, query,
		userId,
		otp,
		time.Now().Add(24*time.Hour*time.Duration(common.PASSWORD_RESET_TIMEOUT_DAYS)),
	)

	return common.FilterSqlPgError(err)
}

func (s *UserServicePgImpl) GetPasswordReset(ctx context.Context, otp string) (models.PasswordReset, error) {
	query := `
		SELECT user_id, otp, exp
		FROM password_resets
		WHERE otp = $1 AND exp > NOW();
	`
	var passReset models.PasswordReset

	err := s.db.QueryRowContext(ctx, query, otp).Scan(
		&passReset.UserId,
		&passReset.Otp,
		&passReset.Exp,
	)

	return passReset, common.FilterSqlPgError(err)
}

func (s *UserServicePgImpl) UpdateUserPassword(ctx context.Context, userId uint32, pw string) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	pwHash, err := common.HashPassword(pw)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE users
		SET password_hash = $1
		WHERE user_id = $2;
	`, pwHash, userId)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM password_resets
		WHERE user_id = $1;
	`, userId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *UserServicePgImpl) LoginOauth(ctx context.Context, oauthUser oauth.User) (models.User, bool, error) {
	user := models.User{}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return user, false, err
	}

	defer tx.Rollback()

	// check if user exists on curr email
	// if not, create and aso create oauth_users entry

	err = tx.QueryRowContext(ctx, `
		SELECT
			user_id,
			email,
			password_hash,
			first_name,
			last_name,
			date_of_birth,
			created_at,
			updated_at,
			is_active
		FROM users WHERE email = $1
	`, oauthUser.Email).Scan(
		&user.UserId,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		// &user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err != sql.ErrNoRows && err != nil {
		return user, false, err
	}

	if err == nil {
		return user, false, err
	}

	// here error is sql.ErrNoRows
	err = tx.QueryRowContext(ctx, `
			INSERT INTO users 
				(email, password_hash, first_name, last_name)
			VALUES
				($1, $2, $3, $4)
			RETURNING *;
		`,
		oauthUser.Email,
		"oauth",
		oauthUser.FirstName,
		oauthUser.LastName,
	).Scan(
		&user.UserId,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		// &user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err != nil {
		return user, false, err
	}

	_, err = tx.ExecContext(ctx, `
			INSERT INTO oauth_users
				(email, user_id, oauth_provider)
			VALUES
				($1, $2, $3);
		`,
		user.Email,
		user.UserId,
		oauthUser.Provider,
	)
	if err != nil {
		return user, false, err
	}

	return user, true, tx.Commit()
}

func (s *UserServicePgImpl) EditUser(ctx context.Context, userId uint32, user schemas.EditUser) error {

	pwHash, err := common.HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, `
		UPDATE users
		SET 
			password_hash = $1,
			first_name = $2,
    		last_name = $3,
			date_of_birth = $4
		WHERE user_id = $5;
	`,
		pwHash,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		userId,
	)

	return err
}
