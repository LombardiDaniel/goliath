package di

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/LombardiDaniel/goliath/src/internal/handlers"
	"github.com/LombardiDaniel/goliath/src/internal/middlewares"
	"github.com/LombardiDaniel/goliath/src/internal/services"
	"github.com/LombardiDaniel/goliath/src/pkg/daemons"

	"github.com/gin-gonic/gin"
)

type Env string

const (
	DevEnv  Env = "dev"
	ProdEnv Env = "prod"
)

type Container struct {
	Ctx context.Context

	AuthService         services.AuthService
	UserService         services.UserService
	EmailService        services.EmailService
	OrganizationService services.OrganizationService
	ObjectService       services.ObjectService
	BillingService      services.BillingService
	TelemetryService    services.TelemetryService

	AuthHandler         handlers.AuthHandler
	UserHandler         handlers.UserHandler
	OrganizationHandler handlers.OrganizationHandler
	BillingHandler      handlers.BillingHandler

	AuthMiddleware      middlewares.AuthMiddleware
	TelemetryMiddleware middlewares.TelemetryMiddleware

	TaskRunner daemons.TaskRunner
}

func NewContainer(env Env) (*Container, error) {
	c := newClients()
	authService := services.NewAuthServiceJwtImpl(os.Getenv("JWT_SECRET_KEY"), c.PostgressConn)
	userService := services.NewUserServicePgImpl(c.PostgressConn)
	var emailService services.EmailService
	if os.Getenv("RESEND_API_KEY") == "mock" {
		emailService = &services.EmailServiceMock{}
	} else {
		emailService = services.NewEmailServiceResendImpl(c.ResendClient)
	}
	organizationService := services.NewOrganizationServicePgImpl(c.PostgressConn)
	objectService := services.NewObjectServiceMinioImpl(c.MinioClient)
	billingService := services.NewBillingService(c.PostgressConn, os.Getenv("STRIPE_API_KEY"))
	telemetryService := services.NewTelemetryServiceMongoAsyncImpl(c.MongoClient, c.MetricsMongoCol, c.EventsMongoCol, 100)

	authMiddleware := middlewares.NewAuthMiddlewareJwt(authService)
	telemetryMiddleware := middlewares.NewTelemetryMiddleware(telemetryService)

	authHandler := handlers.NewAuthHandler(authService, userService, emailService, c.OauthConfigMap)
	userHandler := handlers.NewUserHandler(authService, userService, emailService, objectService)
	organizationHandler := handlers.NewOrganizationHandler(userService, emailService, organizationService)
	billingHandler := handlers.NewBillingHandler(billingService, emailService, userService)

	taskRunner := daemons.TaskRunner{}

	return &Container{
		Ctx:                 context.Background(),
		AuthService:         authService,
		UserService:         userService,
		EmailService:        emailService,
		OrganizationService: organizationService,
		ObjectService:       objectService,
		BillingService:      billingService,
		TelemetryService:    telemetryService,

		AuthHandler:         authHandler,
		UserHandler:         userHandler,
		OrganizationHandler: organizationHandler,
		BillingHandler:      billingHandler,

		AuthMiddleware:      authMiddleware,
		TelemetryMiddleware: telemetryMiddleware,

		TaskRunner: taskRunner,
	}, nil
}

func (c *Container) Setup(router *gin.Engine) {
	c.TaskRunner.RegisterTask(24*time.Hour, c.UserService.DeleteExpiredPwResets, 1)
	c.TaskRunner.RegisterTask(24*time.Hour, c.OrganizationService.DeleteExpiredOrgInvites, 1)
	c.TaskRunner.RegisterTask(time.Second, c.TelemetryService.Upload, 1)

	router.Use(c.TelemetryMiddleware.CollectApiCalls())

	router.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	basePath := router.Group("/v1")
	c.AuthHandler.RegisterRoutes(basePath, c.AuthMiddleware)
	c.UserHandler.RegisterRoutes(basePath, c.AuthMiddleware)
	c.OrganizationHandler.RegisterRoutes(basePath, c.AuthMiddleware)
	c.BillingHandler.RegisterRoutes(basePath, c.AuthMiddleware)
}

func (c *Container) DispatchAsync() {
	c.TaskRunner.Dispatch()
}
