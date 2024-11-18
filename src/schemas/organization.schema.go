package schemas

type OrganizationOutput struct {
	OrganizationId   string `json:"organizationId" binding:"required"`
	OrganizationName string `json:"organizationName" binding:"required"`
	IsAdmin          bool   `json:"isAdmin" binding:"required"`
}

type CreateOrganization struct {
	OrganizationName string `json:"organizationName" binding:"required"`
}

type CreateOrganizationInvite struct {
	OrganizationId string `json:"organizationId" binding:"required,min=1"`
	UserId         uint32 `json:"userId" binding:"required"`
	IsAdmin        bool   `json:"isAdmin" biding:"required"`
}
