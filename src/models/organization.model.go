package models

import "time"

type Organization struct {
	OrganizationID   string         `json:"organizationID" bson:"organizationID" binding:"required,min=1"`
	OrganizationName string         `json:"organizationName" bson:"organizationName" binding:"required,min=1"`
	BillingPlanId    string         `json:"billingPlan" bson:"billingPlan" binding:"required"`
	FrontendConfig   FrontendConfig `json:"frontendConfig" bson:"frontendConfig" binding:"required"`
	CreatedTs        time.Time      `json:"createdTs" bson:"createdTs" binding:"required"`
	Deleted          Deleted        `json:"deleted" bson:"deleted" binding:"required"`
	Owner            string         `json:"owner" bson:"owner" binding:"required,email"`
}

type Deleted struct {
	Deleted bool      `json:"createdTs" bson:"createdTs" binding:"required"`
	Ts      time.Time `json:"ts" bson:"ts" binding:"required"`
}

type FrontendConfig struct {
	PrimaryColor   string
	SecondaryColor string
}
