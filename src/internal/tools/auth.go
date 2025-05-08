package tools

import (
	"github.com/LombardiDaniel/goliath/src/pkg/constants"
	"github.com/gin-gonic/gin"
)

func GetClaimsFromGinCtx[T any](ctx *gin.Context) (T, error) {
	claims, ok := ctx.Get(constants.GinCtxJwtClaimKeyName)
	var zero T
	if !ok {
		return zero, constants.ErrAuth
	}

	parsedClaims, ok := claims.(T)
	if !ok {
		return zero, constants.ErrAuth
	}

	return parsedClaims, nil
}
