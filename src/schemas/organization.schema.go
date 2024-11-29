package schemas

type OrganizationOutput struct {
	OrganizationId   string `json:"organizationId" binding:"required"`
	OrganizationName string `json:"organizationName" binding:"required"`
	IsAdmin          bool   `json:"isAdmin" binding:"required"`
	IsOwner          bool   `json:"isOwner" binding:"required"`
}

type CreateOrganization struct {
	OrganizationName string `json:"organizationName" binding:"required"`
}

type CreateOrganizationInvite struct {
	// UserId  uint32 `json:"userId" binding:"required"`
	UserEmail string `json:"userEmail" binding:"required"`
	IsAdmin   bool   `json:"isAdmin" biding:"required"`
}
