package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/LombardiDaniel/gopherbase/common"
	"github.com/LombardiDaniel/gopherbase/models"
	"github.com/LombardiDaniel/gopherbase/schemas"
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

	var count uint32 = 0
	err = tx.QueryRowContext(ctx, `
			SELECT COUNT(user_id)
			FROM users WHERE email = $1;
		`,
		unconfirmedUser.Email,
	).Scan(&count)
	if err != nil {
		return common.FilterSqlPgError(err)
	}

	if count != 0 {
		return common.ErrDbConflict
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
			avatar_url,
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
		&user.AvatarUrl,
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
			avatar_url,
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
		&user.AvatarUrl,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.IsActive,
	)
	if err != nil {
		return user, common.FilterSqlPgError(err)
	}

	return user, nil
}

func (s *UserServicePgImpl) GetUsers(ctx context.Context) ([]models.User, error) {
	query := `
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
		FROM users;
	`

	users := []models.User{}

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return users, common.FilterSqlPgError(err)
	}
	defer rows.Close()

	for rows.Next() {
		u := models.User{}
		err := rows.Scan(
			&u.UserId,
			&u.Email,
			&u.PasswordHash,
			&u.FirstName,
			&u.LastName,
			&u.DateOfBirth,
			&u.AvatarUrl,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.IsActive,
		)
		if err != nil {
			return users, err
		}
		users = append(users, u)
	}

	return users, nil
}

func (s *UserServicePgImpl) GetUserOrgs(ctx context.Context, userId uint32) ([]schemas.OrganizationOutput, error) {
	query := `
		SELECT DISTINCT
			o.organization_id,
			o.organization_name,
			ou.is_admin,
			o.owner_user_id = ou.user_id
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
		err := rows.Scan(&newOrg.OrganizationId, &newOrg.OrganizationName, &newOrg.IsAdmin, &newOrg.IsOwner)
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

func (s *UserServicePgImpl) EditUser(ctx context.Context, userId uint32, user schemas.EditUser) error {

	_, err := s.db.ExecContext(ctx, `
			UPDATE users
			SET 
				first_name = $1,
				last_name = $2,
				date_of_birth = $3
			WHERE user_id = $4;
		`,
		user.FirstName,
		user.LastName,
		user.DateOfBirth,
		userId,
	)

	return err
}

func (s *UserServicePgImpl) DeleteExpiredPwResets() error {
	_, err := s.db.Exec(`
		DELETE FROM password_resets
    	WHERE exp < NOW();
	`)
	return err
}

func (s *UserServicePgImpl) SetAvatarUrl(ctx context.Context, userId uint32, url string) error {
	_, err := s.db.ExecContext(ctx, `
			UPDATE users
			SET 
				avatar_url = $1
			WHERE user_id = $2;
		`,
		url,
		userId,
	)

	return err
}
