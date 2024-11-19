package services

import (
	"context"
	"database/sql"

	"github.com/LombardiDaniel/go-gin-template/common"
	"github.com/LombardiDaniel/go-gin-template/models"
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
			deleted_at
		FROM
			organizations
		WHERE
			organization_id = $1;
	`

	org := models.Organization{}

	err := s.db.QueryRowContext(ctx, query, orgId).Scan(
		&org.OrganizationId,
		&org.OrganizationName,
		&org.BillingPlanId,
		&org.CreatedAt,
		&org.DeletedAt,
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
				invite_otp,
				invite_exp
			)
			VALUES ($1, $2, $3, $4, $5);
	`

	_, err := s.db.ExecContext(ctx, query,
		invite.OrganizationId,
		invite.UserId,
		invite.IsAdmin,
		invite.InviteOtp,
		invite.InviteExp,
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
		SELECT * FROM organization_invites
		WHERE otp = $1 AND exp > NOW();
	`, otp).Scan(
		&inv.OrganizationId,
		&inv.UserId,
		&inv.IsAdmin,
		&inv.InviteOtp,
		&inv.InviteExp,
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
