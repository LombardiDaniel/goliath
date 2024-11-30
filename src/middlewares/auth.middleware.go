package middlewares

import (
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	AuthorizeUser() gin.HandlerFunc
	AuthorizeOrganization(needAdmin bool) gin.HandlerFunc
	Reauthorize() gin.HandlerFunc
}
