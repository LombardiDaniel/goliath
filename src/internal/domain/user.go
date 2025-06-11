package domain

import (
	"time"

	"github.com/LombardiDaniel/goliath/src/pkg/common"
	"github.com/LombardiDaniel/goliath/src/pkg/constants"
	"github.com/LombardiDaniel/goliath/src/pkg/token"
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

func NewUnconfirmedUser(email string, password string, firstName string, lastName string, dateOfBirth *time.Time) (*UnconfirmedUser, error) {
	hash, err := token.HashPassword(password)
	if err != nil {
		return nil, err
	}

	otp, err := common.GenerateRandomString(constants.OtpLen)
	if err != nil {
		return nil, err
	}

	return &UnconfirmedUser{
		Email:        email,
		Otp:          otp,
		PasswordHash: hash,
		FirstName:    firstName,
		LastName:     lastName,
		DateOfBirth:  dateOfBirth,
	}, nil
}
