package models

import (
	"math"
	"time"

	"github.com/golang-jwt/jwt"
)

type Permissions int32

const (
	NonePermissions      Permissions = 0
	ReadPermissions      Permissions = 1 << 0
	WritePermissions     Permissions = 1 << 1
	ReadWritePermissions Permissions = ReadPermissions | WritePermissions
	AllPermissions       Permissions = math.MaxInt32
)

type OrganizationPermission struct {
	OrganizationId string      `json:"organizationId" binding:"required,min=1"`
	UserId         uint32      `json:"userId" binding:"required"`
	ActionName     string      `json:"actionName" binding:"required"`
	Perms          Permissions `json:"perms" binding:"required"`
}

// JwtClaims represents the claims in a JWT token.
type JwtClaims struct {
	UserId         uint32                 `json:"userId" binding:"required"`
	Email          string                 `json:"email" binding:"required"`
	OrganizationId *string                `json:"organizationId" binding:"required"`
	IsAdmin        *bool                  `json:"isAdmin" binding:"required"`
	ActionPerms    map[string]Permissions `json:"actionPerms" binding:"required"`

	jwt.StandardClaims
}

// only here because swaggo cant expand the above example (but same thing, KEEP IN SYNC!!)
type JwtClaimsOutput struct {
	UserId         uint32                 `json:"userId" binding:"required"`
	Email          string                 `json:"email" binding:"required"`
	OrganizationId *string                `json:"organizationId" binding:"required"`
	IsAdmin        *bool                  `json:"isAdmin" binding:"required"`
	ActionPerms    map[string]Permissions `json:"actionPerms" binding:"required"`

	Audience  string `json:"aud"`
	ExpiresAt int64  `json:"exp"`
	Id        string `json:"jti"`
	IssuedAt  int64  `json:"iat"`
	Issuer    string `json:"iss"`
	NotBefore int64  `json:"nbf"`
	Subject   string `json:"sub"`
}

// PasswordReset represents a password reset struct.
type PasswordReset struct {
	UserId uint32
	Otp    string
	Exp    time.Time
}

// JwtPasswordResetClaims represents the claims in a JWT token for password reset.
type JwtPasswordResetClaims struct {
	UserId  uint32 `json:"userId" binding:"required"`
	Allowed bool   `json:"allowrd" binding:"required"`

	jwt.StandardClaims
}
