package common

import (
	"github.com/LombardiDaniel/go-gin-template/models"
	"github.com/gin-gonic/gin"
)

func GetClaimsFromGinCtx(ctx *gin.Context) (models.JwtClaims, error) {
	claims, ok := ctx.Get(GIN_CTX_JWT_CLAIM_KEY_NAME)
	if !ok {
		return models.JwtClaims{}, ErrAuth
	}

	parsedClaims, ok := claims.(models.JwtClaims)
	if !ok {
		return models.JwtClaims{}, ErrAuth
	}

	return parsedClaims, nil
}
