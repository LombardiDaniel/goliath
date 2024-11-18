package middlewares

import (
	"github.com/gin-gonic/gin"
)

type AuthMiddleware interface {
	Authorize() gin.HandlerFunc
	// AuthorizeOrganiaztion(needAdmin bool) gin.HandlerFunc
	// AuthorizeNoOrganization() gin.HandlerFunc
}
