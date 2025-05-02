package models

import (
	"time"
)

// User represents a user in the system.
type User struct {
	UserId       uint32
	Perms        map[string]Permission
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

// UnconfirmedUser represents a user who has not yet confirmed their email address.
type UnconfirmedUser struct {
	Email        string
	Otp          string
	PasswordHash string
	FirstName    string
	LastName     string
	DateOfBirth  *time.Time
}
