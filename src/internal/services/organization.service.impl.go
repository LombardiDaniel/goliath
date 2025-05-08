package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/LombardiDaniel/gopherbase/src/internal/domain"
	"github.com/LombardiDaniel/gopherbase/src/pkg/constants"
	"github.com/LombardiDaniel/gopherbase/src/pkg/validators"
)

type OrganizationServicePgImpl struct {
	db *sql.DB
}

func NewOrganizationServicePgImpl(db *sql.DB) OrganizationService {
	return &OrganizationServicePgImpl{
		db: db,
	}
}

func (s *OrganizationServicePgImpl) GetOrganization(ctx context.Context, orgId string) (domain.Organization, error) {
	query := `
		SELECT
			organization_id,
			organization_name,
			billing_plan_id,
			created_at,
			deleted_at,
			owner_user_id
		FROM
			organizations
		WHERE
			organization_id = $1;
	`

	org := domain.Organization{}

	err := s.db.QueryRowContext(ctx, query, orgId).Scan(
		&org.OrganizationId,
		&org.OrganizationName,
		&org.BillingPlanId,
		&org.CreatedAt,
		&org.DeletedAt,
		&org.OwnerUserId,
	)
	return org, errors.Join(err, validators.FilterSqlPgError(err))
}

func (s *OrganizationServicePgImpl) CreateOrganization(ctx context.Context, org domain.Organization) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Join(err, constants.ErrDbTransactionCreate)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
			INSERT INTO organizations (organization_id, organization_name, owner_user_id)
			VALUES ($1, $2, $3);
		`,
		org.OrganizationId,
		org.OrganizationName,
		org.OwnerUserId,
	)
	if err != nil {
		return errors.Join(err, validators.FilterSqlPgError(err))
	}

	_, err = tx.ExecContext(ctx, `
			INSERT INTO organizations_users (organization_id, user_id)
			VALUES ($1, $2);
		`,
		org.OrganizationId,
		org.OwnerUserId,
	)
	if err != nil {
		return errors.Join(err, validators.FilterSqlPgError(err))
	}

	perms := map[string]domain.Permission{
		"admin": domain.AllPermission,
		"owner": domain.AllPermission,
	}
	for action, perm := range perms {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO organization_user_permissions (organization_id, user_id, action_name, permission)
			VALUES ($1, $2, $3, $4);
		`,
			org.OrganizationId,
			org.OwnerUserId,
			action,
			perm,
		)
		if err != nil {
			return errors.Join(err, validators.FilterSqlPgError(err))
		}
	}

	err = tx.Commit()
	return errors.Join(err, validators.FilterSqlPgError(err))
}

func (s *OrganizationServicePgImpl) CreateOrganizationInvite(ctx context.Context, invite domain.OrganizationInvite) error {
	query := `
		INSERT INTO organization_invites (
			organization_id,
			user_id,
			perms_json,
			otp,
			exp
		)
		VALUES ($1, $2, $3, $4, $5);
	`

	permStr, err := json.Marshal(invite.Perms)
	if err != nil {
		return errors.Join(err, errors.New("could not marshal invite perms"))
	}
	_, err = s.db.ExecContext(ctx, query,
		invite.OrganizationId,
		invite.UserId,
		permStr,
		invite.Otp,
		invite.Exp,
	)
	return errors.Join(err, validators.FilterSqlPgError(err))
}

func (s *OrganizationServicePgImpl) ConfirmOrganizationInvite(ctx context.Context, otp string) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Join(err, constants.ErrDbTransactionCreate)
	}
	defer tx.Rollback()

	var inv domain.OrganizationInvite
	var permsString string
	err = tx.QueryRowContext(ctx, `
		SELECT
			organization_id,
			user_id,
			perms_json,
			otp,
			exp
		FROM
			organization_invites
		WHERE
			otp = $1 AND
			exp > NOW();
	`, otp).Scan(
		&inv.OrganizationId,
		&inv.UserId,
		&permsString,
		&inv.Otp,
		&inv.Exp,
	)
	if err != nil {
		return errors.Join(err, validators.FilterSqlPgError(err))
	}

	err = json.Unmarshal([]byte(permsString), &inv.Perms)
	if err != nil {
		return errors.Join(err, errors.New("could not unmarshal perms to json"))
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO organizations_users (organization_id, user_id)
		VALUES ($1, $2);
	`,
		inv.OrganizationId,
		inv.UserId,
	)
	if err != nil {
		return errors.Join(err, validators.FilterSqlPgError(err))
	}

	for action, perm := range inv.Perms {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO organization_user_permissions (organization_id, user_id, action_name, permission)
			VALUES ($1, $2, $3, $4);
		`,
			inv.OrganizationId,
			inv.UserId,
			action,
			perm,
		)
		if err != nil {
			return errors.Join(err, validators.FilterSqlPgError(err))
		}
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM organization_invites
		WHERE otp = $1;
	`, otp)
	if err != nil {
		return errors.Join(err, validators.FilterSqlPgError(err))
	}

	return tx.Commit()
}

func (s *OrganizationServicePgImpl) RemoveUserFromOrg(ctx context.Context, orgId string, userId uint32) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Join(err, constants.ErrDbTransactionCreate)
	}
	defer tx.Rollback()

	var isOwner bool
	err = s.db.QueryRowContext(ctx, `
		SELECT owner_user_id = $1
		FROM organizations
		WHERE organization_id = $2;
	`,
		userId,
		orgId,
	).Scan(&isOwner)
	if err != nil {
		return errors.Join(err, validators.FilterSqlPgError(err))
	}

	if isOwner {
		return errors.Join(constants.ErrDbConflict, errors.New("cannot remove owner of organization"))
	}

	_, err = s.db.ExecContext(ctx, `
		DELETE FROM organizations_users
		WHERE organization_id = $1 AND user_id = $2;
	`,
		orgId,
		userId,
	)
	if err != nil {
		return errors.Join(err, validators.FilterSqlPgError(err))
	}

	return tx.Commit()
}

func (s *OrganizationServicePgImpl) SetOrganizationOwner(ctx context.Context, orgId string, userId uint32) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return errors.Join(err, constants.ErrDbTransactionCreate)
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
		UPDATE organizations
		SET owner_user_id = $1
		WHERE organization_id = $2;
	`,
		userId,
		orgId,
	)
	if err != nil {
		return errors.Join(err, validators.FilterSqlPgError(err))
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE organization_user_permissions
		SET user_id = $1
		WHERE
			action_name = 'owner' AND
			organization_id = $2;
	`,
		userId,
		orgId,
	)
	if err != nil {
		return errors.Join(err, validators.FilterSqlPgError(err))
	}

	return tx.Commit()
}

func (s *OrganizationServicePgImpl) DeleteExpiredOrgInvites() error {
	_, err := s.db.Exec(`
		DELETE FROM organization_invites
    	WHERE exp < NOW();
	`)
	return errors.Join(err, validators.FilterSqlPgError(err))
}

func (s *OrganizationServicePgImpl) SetPerms(ctx context.Context, action string, userId uint32, perms domain.Permission) error {
	_, err := s.db.Exec(`
        INSERT INTO organization_user_permissions (organization_id, user_id, action_name, permission)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (action_name, organization_id, user_id)
        DO UPDATE SET
            permission = EXCLUDED.permission;
    `)
	return errors.Join(err, validators.FilterSqlPgError(err))
}
