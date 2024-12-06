package models

import (
	"time"

	"github.com/golang-jwt/jwt"
)

// type OauthUser struct {
// 	Email         string
// 	UserId        uint32
// 	OauthProvider string
// }

type JwtClaims struct {
	UserId         uint32  `json:"userId" binding:"required"`
	Email          string  `json:"email" binding:"required"`
	OrganizationId *string `json:"organizationId" binding:"required"`
	IsAdmin        *bool   `json:"isAdmin" binding:"required"`

	jwt.StandardClaims
}

// only here because swaggo cant expand the above example (but same thing, KEEP IN SYNC!!)
type JwtClaimsOutput struct {
	UserId         uint32  `json:"userId" binding:"required"`
	Email          string  `json:"email" binding:"required"`
	OrganizationId *string `json:"organizationId" binding:"required"`

	Audience  string `json:"aud"`
	ExpiresAt int64  `json:"exp"`
	Id        string `json:"jti"`
	IssuedAt  int64  `json:"iat"`
	Issuer    string `json:"iss"`
	NotBefore int64  `json:"nbf"`
	Subject   string `json:"sub"`
}

type PasswordReset struct {
	UserId uint32
	Otp    string
	Exp    time.Time
}

type JwtPasswordResetClaims struct {
	UserId  uint32 `json:"userId" binding:"required"`
	Allowed bool   `json:"allowrd" binding:"required"`

	jwt.StandardClaims
}
