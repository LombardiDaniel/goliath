package middlewares

import (
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
	// the user is authorized to access organization-specific resources. The
	// needAdmin parameter determines if administrative privileges are required.
	AuthorizeOrganization(needAdmin bool) gin.HandlerFunc

	// Reauthorize returns a middleware handler function that handles
	// reauthorization logic, such as refreshing tokens or revalidating
	// user sessions.
	Reauthorize() gin.HandlerFunc
}
