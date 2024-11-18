package services

import (
	"context"
	"database/sql"

	"github.com/LombardiDaniel/go-gin-template/common"
	"github.com/LombardiDaniel/go-gin-template/models"
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
		VALUES ($1, $2, $3, $4, $5)
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
			SELECT user_id FROM users WHERE email = $1
		`,
		unconfirmedUser.Email,
	).Scan(&userId)

	if err != sql.ErrNoRows {
		if err == nil {
			return common.ErrDbConflict
		}
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
			SELECT * FROM unconfirmed_users WHERE otp = $1
		`, otp).Scan(
		&unconfirmedUser.Email,
		&unconfirmedUser.Otp,
		&unconfirmedUser.PasswordHash,
		&unconfirmedUser.FirstName,
		&unconfirmedUser.LastName,
		&unconfirmedUser.DateOfBirth,
	)
	if err != nil {
		return err
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
		return err
	}

	_, err = tx.ExecContext(ctx, `
			DELETE FROM unconfirmed_users WHERE otp = $1;
		`,
		unconfirmedUser.Otp,
	)
	if err != nil {
		return err
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
		&user.LastLogin,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *UserServicePgImpl) GetUsers(ctx context.Context) ([]models.User, error) {
	return nil, nil
}

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
		return orgs, err
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
