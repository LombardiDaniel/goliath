package fiddlers

import (
	"time"

	"github.com/LombardiDaniel/gopherbase/common"
	"github.com/LombardiDaniel/gopherbase/models"
)

func NewOrganization(orgName string, ownerId uint32) (*models.Organization, error) {
	orgId, err := common.GenerateRandomString(5)
	if err != nil {
		return nil, err
	}

	return &models.Organization{
		OrganizationId:   orgId,
		OrganizationName: orgName,
		OwnerUserId:      ownerId,
	}, nil
}

func NewOrganizationInvite(organizationId string, userId uint32, isAdmin bool, otp string) models.OrganizationInvite {
	invExp := time.Now().Add(24 * time.Hour * time.Duration(common.ORG_INVITE_TIMEOUT_DAYS))
	return models.OrganizationInvite{
		OrganizationId: organizationId,
		UserId:         userId,
		IsAdmin:        isAdmin,
		Otp:            &otp,
		Exp:            &invExp,
	}
}
