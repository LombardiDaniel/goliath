package services

import (
	"context"
	"database/sql"

	"github.com/LombardiDaniel/gopherbase/common"
	"github.com/LombardiDaniel/gopherbase/models"
)

type OrganizationServicePgImpl struct {
	db *sql.DB
}

func NewOrganizationServicePgImpl(db *sql.DB) OrganizationService {
	return &OrganizationServicePgImpl{
		db: db,
	}
}

func (s *OrganizationServicePgImpl) GetOrganization(ctx context.Context, orgId string) (models.Organization, error) {
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
	// query := `
	// 	SELECT
	// 		o.organization_id,
	// 		o.organization_name,
	// 		o.billing_plan_id,
	// 		o.created_at,
	// 		o.deleted_at,
	// 		u.email
	// 	FROM
	// 		organizations o
	// 	INNER JOIN
	// 		organizations_users ou ON o.organization_id = ou.organization_id
	// 	INNER JOIN
	// 		users u ON ou.user_id = u.user_id
	// 	WHERE
	// 		ou.user_id = $1;
	// `

	org := models.Organization{}

	err := s.db.QueryRowContext(ctx, query, orgId).Scan(
		&org.OrganizationId,
		&org.OrganizationName,
		&org.BillingPlanId,
		&org.CreatedAt,
		&org.DeletedAt,
		&org.OwnerUserId,
	)

	return org, err
}

func (s *OrganizationServicePgImpl) CreateOrganization(ctx context.Context, org models.Organization) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
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
		return err
	}

	_, err = tx.ExecContext(ctx, `
			INSERT INTO organizations_users (organization_id, user_id, is_admin)
			VALUES ($1, $2, true);
		`,
		org.OrganizationId,
		org.OwnerUserId,
	)
	if err != nil {
		return err
	}

	err = tx.Commit()
	return common.FilterSqlPgError(err)
}

func (s *OrganizationServicePgImpl) CreateOrganizationInvite(ctx context.Context, invite models.OrganizationInvite) error {
	query := `
		INSERT INTO organization_invites (
			organization_id,
			user_id,
			is_admin,
			otp,
			exp
		)
		VALUES ($1, $2, $3, $4, $5);
	`

	_, err := s.db.ExecContext(ctx, query,
		invite.OrganizationId,
		invite.UserId,
		invite.IsAdmin,
		invite.Otp,
		invite.Exp,
	)
	return common.FilterSqlPgError(err)
}

func (s *OrganizationServicePgImpl) ConfirmOrganizationInvite(ctx context.Context, otp string) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var inv models.OrganizationInvite
	err = tx.QueryRowContext(ctx, `
		SELECT
			organization_id,
			user_id,
			is_admin,
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
		&inv.IsAdmin,
		&inv.Otp,
		&inv.Exp,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO organizations_users (organization_id, user_id, is_admin)
		VALUES ($1, $2, $3);
	`,
		inv.OrganizationId,
		inv.UserId,
		inv.IsAdmin,
	)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		DELETE FROM organization_invites
		WHERE otp = $1;
	`, otp)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *OrganizationServicePgImpl) RemoveUserFromOrg(ctx context.Context, orgId string, userId uint32) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
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
		return err
	}

	if isOwner {
		return common.ErrDbConflict
	}

	_, err = s.db.ExecContext(ctx, `
		DELETE FROM organizations_users
		WHERE organization_id = $1 AND user_id = $2;
	`,
		orgId,
		userId,
	)
	if err != nil {
		return err
	}

	return common.FilterSqlPgError(tx.Commit())
}

func (s *OrganizationServicePgImpl) SetOrganizationOwner(ctx context.Context, orgId string, userId uint32) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
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
		return err
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE organizations_users
		SET is_admin = true
		WHERE organization_id = $1 AND user_id = $2;
	`,
		orgId,
		userId,
	)
	if err != nil {
		return err
	}

	return common.FilterSqlPgError(tx.Commit())
}

func (s *OrganizationServicePgImpl) DeleteExpiredOrgInvites() error {
	_, err := s.db.Exec(`
		DELETE FROM organization_invites
    	WHERE exp < NOW();
	`)
	return err
}
