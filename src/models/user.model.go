package models

import "time"

type User struct {
	Email                string                `json:"email" bson:"email" binding:"required,email"`
	Name                 string                `json:"name" bson:"name" binding:"required,min=1,max=256"`
	Password             string                `json:"password" bson:"password" binding:"required"`
	OrganizationsDetails []OrganizationDetails `json:"organizationsDetails" bson:"organizationsDetails" binding:"required"`
	LastLogin            *time.Time            `json:"lastLogin" bson:"lastLogin" binding:"required"`
}

type Invite struct {
	// Id					primitive.ObjectID			`json:"_id" bson:"_id" binding:"required"`
	OrganizationDetails OrganizationDetails `json:"organizationDetails" bson:"organizationDetails" binding:"required"`
	Email               string              `json:"email" bson:"email" binding:"required"`
	Declined            bool                `json:"declined" bson:"declined" binding:"required"`
	Accepted            bool                `json:"accepted" bson:"accepted" binding:"required"`
	Ts                  time.Time           `json:"ts" bson:"ts" binding:"required"`
}

type OrganizationDetails struct {
	OrganizationID string `json:"organizationID" bson:"organizationID" binding:"required,min=1"`
	Role           string `json:"role" bson:"role" binding:"required"`
	Enabled        bool   `json:"enabled" bson:"enabled" binding:"required"`
	RegistrationNo string `json:"registrationNo" bson:"registrationNo" binding:"required"`
}

type UnconfirmedUser struct {
	Email    string `json:"email" bson:"email" binding:"required,email"`
	Name     string `json:"name" bson:"name" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
	Otp      string `json:"otp" bson:"otp" binding:"required,min=1,max=256"`
}
