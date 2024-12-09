package models

import "time"

type User struct {
	UserId       uint32
	Email        string
	PasswordHash string
	FirstName    string
	LastName     string
	AvatarUrl    *string
	DateOfBirth  *time.Time
	// LastLogin    *time.Time
	CreatedAt *time.Time
	UpdatedAt *time.Time
	IsActive  bool
}

type UnconfirmedUser struct {
	Email        string
	Otp          string
	PasswordHash string
	FirstName    string
	LastName     string
	DateOfBirth  *time.Time
}

// type User struct {
// 	UserId       int        `json:"userId" binding:"required"`
// 	Email        string     `json:"email" binding:"required,email,max=100"`
// 	PasswordHash string     `json:"passwordHash" binding:"required,max=255"`
// 	FirstName    string     `json:"firstName" binding:"required,min=1,max=50"`
// 	LastName     string     `json:"lastName" binding:"required,min=1,max=50"`
// 	DateOfBirth  *time.Time `json:"dateOfBirth" binding:"required"`
// 	LastLogin    *time.Time `json:"lastLogin" binding:"required"`
// 	CreatedAt    *time.Time `json:"createdAt" binding:"required"`
// 	IsActive     bool       `json:"isActive" binding:"required"`
// }

// type Invite struct {
// 	// Id					primitive.ObjectID			`json:"_id" bson:"_id" binding:"required"`
// 	OrganizationDetails OrganizationDetails `json:"organizationDetails" bson:"organizationDetails" binding:"required"`
// 	Email               string              `json:"email" bson:"email" binding:"required"`
// 	Declined            bool                `json:"declined" bson:"declined" binding:"required"`
// 	Accepted            bool                `json:"accepted" bson:"accepted" binding:"required"`
// 	Ts                  time.Time           `json:"ts" bson:"ts" binding:"required"`
// }

// type OrganizationDetails struct {
// 	OrganizationID string `json:"organizationID" bson:"organizationID" binding:"required,min=1"`
// 	Role           string `json:"role" bson:"role" binding:"required"`
// 	Enabled        bool   `json:"enabled" bson:"enabled" binding:"required"`
// 	RegistrationNo string `json:"registrationNo" bson:"registrationNo" binding:"required"`
// }

// type UnconfirmedUser struct {
// 	Email    string `json:"email" bson:"email" binding:"required,email"`
// 	Name     string `json:"name" bson:"name" binding:"required"`
// 	Password string `json:"password" bson:"password" binding:"required"`
// 	Otp      string `json:"otp" bson:"otp" binding:"required,min=1,max=256"`
// }
