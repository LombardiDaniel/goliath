package dto

import "github.com/LombardiDaniel/goliath/src/internal/models"

type OrganizationOutput struct {
	OrganizationId   string                       `json:"organizationId" binding:"required"`
	OrganizationName string                       `json:"organizationName" binding:"required"`
	Perms            map[string]models.Permission `json:"perms" binding:"required"`
	IsOwner          bool                         `json:"isOwner" binding:"required"`
}

type CreateOrganization struct {
	OrganizationName string `json:"organizationName" binding:"required"`
}

type CreateOrganizationInvite struct {
	UserEmail string                       `json:"userEmail" binding:"required"`
	Perms     map[string]models.Permission `json:"perms" binding:"required"`
}
