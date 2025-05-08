package middlewares

import "github.com/gin-gonic/gin"

// TelemetryMiddleware defines an interface for telemetry-related middleware.
// It includes methods for collecting and processing telemetry data.
type TelemetryMiddleware interface {
	// CollectApiCall returns a middleware handler function that collects
	// telemetry data for API calls, such as request metrics or logging.
	CollectApiCalls() gin.HandlerFunc
}
