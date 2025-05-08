package middlewares

import (
	"github.com/LombardiDaniel/gopherbase/src/internal/domain"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware defines an interface for authentication and authorization
// middleware. It includes methods for user authorization, organization
// authorization, and reauthorization.
type AuthMiddleware interface {
	// AuthorizeUser returns a middleware handler function that ensures
	// the user is authorized to access the requested resource.
	AuthorizeUser() gin.HandlerFunc

	// AuthorizeOrganization returns a middleware handler function that ensures
	// the user is authorized to access organization-specific resources.
	AuthorizeOrganization(need map[string]domain.Permission) gin.HandlerFunc

	// Reauthorize returns a middleware handler function that handles
	// reauthorization logic, such as refreshing tokens or revalidating
	// user sessions.
	Reauthorize() gin.HandlerFunc
}
