package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/LombardiDaniel/goliath/src/internal/handlers"
	"github.com/LombardiDaniel/goliath/src/internal/middlewares"
	"github.com/LombardiDaniel/goliath/src/internal/services"
	"github.com/LombardiDaniel/goliath/src/pkg/common"
	"github.com/LombardiDaniel/goliath/src/pkg/constants"
	"github.com/LombardiDaniel/goliath/src/pkg/daemons"
	"github.com/LombardiDaniel/goliath/src/pkg/it"
	"github.com/LombardiDaniel/goliath/src/pkg/logger"
	"github.com/LombardiDaniel/goliath/src/pkg/oauth"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"github.com/gin-contrib/cors"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/basic/docs"
)

var (
	ctx    context.Context
	router *gin.Engine

	authService         services.AuthService
	userService         services.UserService
	emailService        services.EmailService
	organizationService services.OrganizationService
	objectService       services.ObjectService
	billingService      services.BillingService
	telemetryService    services.TelemetryService

	authHandler         handlers.AuthHandler
	userHandler         handlers.UserHandler
	organizationHandler handlers.OrganizationHandler
	billingHandler      handlers.BillingHandler

	authMiddleware      middlewares.AuthMiddleware
	telemetryMiddleware middlewares.TelemetryMiddleware

	taskRunner daemons.TaskRunner
)

func init() {
	logger.InitSlogger()
	ctx = context.Background()

	pgConnStr := common.GetEnvVarDefault("POSTGRES_URI", "postgres://user:password@localhost:5432/db?sslmode=disable")
	db := it.Must(sql.Open("postgres", pgConnStr))
	if err := db.Ping(); err != nil {
		panic(errors.Join(err, errors.New("could not ping pgsql")))
	}

	pgIdleConns := it.Must(strconv.Atoi(common.GetEnvVarDefault("POSTGRES_IDLE_CONNS", "2")))
	pgOpenConns := it.Must(strconv.Atoi(common.GetEnvVarDefault("POSTGRES_OPEN_CONNS", "10")))
	db.SetMaxIdleConns(pgIdleConns)
	db.SetMaxOpenConns(pgOpenConns)
	it.Must(db.Exec(fmt.Sprintf("SET TIME ZONE '%s';", constants.DefaultTimzone)))

	mongoConn := options.Client().ApplyURI(
		common.GetEnvVarDefault("MONGO_URI", "mongodb://localhost:27017"),
	)
	mongoClient := it.Must(mongo.Connect(ctx, mongoConn))
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		panic(errors.Join(err, errors.New("could not ping mongodb")))
	}

	tsIdxModel := mongo.IndexModel{
		Keys:    bson.M{"ts": 1},
		Options: options.Index(),
	}
	metricsCol := mongoClient.Database("telemetry").Collection("metrics")
	eventsCol := mongoClient.Database("telemetry").Collection("events")
	it.Must(metricsCol.Indexes().CreateOne(ctx, tsIdxModel))
	it.Must(eventsCol.Indexes().CreateOne(ctx, tsIdxModel))

	oauthBaseCallback := constants.ApiHostUrl + "v1/auth/%s/callback"

	oauthConfigMap := make(map[string]oauth.Provider)
	oauthConfigMap[oauth.GOOGLE_PROVIDER] = oauth.NewGoogleProvider(&oauth2.Config{
		RedirectURL:  fmt.Sprintf(oauthBaseCallback, oauth.GOOGLE_PROVIDER),
		ClientID:     os.Getenv("OAUTH_GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_GOOGLE_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	})
	oauthConfigMap[oauth.GITHUB_PROVIDER] = oauth.NewGithubProvider(&oauth2.Config{
		RedirectURL:  fmt.Sprintf(oauthBaseCallback, oauth.GITHUB_PROVIDER),
		ClientID:     os.Getenv("OAUTH_GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
		Scopes: []string{
			"read:user",
		},
		Endpoint: github.Endpoint,
	})

	s3Host := it.Must(common.ExtractHostFromUrl(constants.S3Endpoint))
	s3Secure := it.Must(common.UrlIsSecure(constants.S3Endpoint))
	minioClient := it.Must(minio.New(
		s3Host,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				os.Getenv("S3_ACCESS_KEY_ID"),
				os.Getenv("S3_SECRET_ACCESS_KEY"),
				"",
			),
			Region: constants.S3Region,
			Secure: s3Secure,
		},
	))

	authService = services.NewAuthServiceJwtImpl(os.Getenv("JWT_SECRET_KEY"), db)
	userService = services.NewUserServicePgImpl(db)
	if os.Getenv("RESEND_API_KEY") == "mock" {
		emailService = &services.EmailServiceMock{}
	} else {
		emailService = services.NewEmailServiceResendImpl(os.Getenv("RESEND_API_KEY"), "internal/templates")
	}
	organizationService = services.NewOrganizationServicePgImpl(db)
	objectService = services.NewObjectServiceMinioImpl(minioClient)
	billingService = services.NewBillingService(db, os.Getenv("STRIPE_API_KEY"))
	telemetryService = services.NewTelemetryServiceMongoAsyncImpl(mongoClient, metricsCol, eventsCol, 100)

	authMiddleware = middlewares.NewAuthMiddlewareJwt(authService)
	telemetryMiddleware = middlewares.NewTelemetryMiddleware(telemetryService)

	authHandler = handlers.NewAuthHandler(authService, userService, emailService, oauthConfigMap)
	userHandler = handlers.NewUserHandler(authService, userService, emailService, objectService)
	organizationHandler = handlers.NewOrganizationHandler(userService, emailService, organizationService)
	billingHandler = handlers.NewBillingHandler(billingService, emailService, userService)

	router = gin.Default()
	router.SetTrustedProxies([]string{"*"})

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{constants.ApiHostUrl, constants.AppHostUrl}
	corsCfg.AllowCredentials = true
	corsCfg.AddAllowHeaders("Authorization")
	corsCfg.MaxAge = 24 * time.Hour

	slog.Info(fmt.Sprintf("corsCfg: %+v", corsCfg))

	router.Use(cors.New(corsCfg))
	router.Use(limits.RequestSizeLimiter(constants.MaxRequestSize))

	docs.SwaggerInfo.Title = "Goliath"
	docs.SwaggerInfo.Description = "Goliath"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Host = strings.Split(constants.ApiHostUrl, "://")[1]

	if os.Getenv("GIN_MODE") == "release" {
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs", func(ctx *gin.Context) {
		ctx.Header("location", "/docs/index.html")
		ctx.String(http.StatusMovedPermanently, "MovedPermanently")
	})

	// Daemons
	taskRunner.RegisterTask(24*time.Hour, userService.DeleteExpiredPwResets, 1)
	taskRunner.RegisterTask(24*time.Hour, organizationService.DeleteExpiredOrgInvites, 1)
	taskRunner.RegisterTask(time.Second, telemetryService.Upload, 1)
}

// @securityDefinitions.apiKey JWT
// @in cookie
// @name Authorization
// @description JWT
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer $TOKEN' to correctly set the API Key"
func main() {

	// LB healthcheck
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	router.Use(telemetryMiddleware.CollectApiCalls())

	basePath := router.Group("/v1")
	authHandler.RegisterRoutes(basePath, authMiddleware)
	userHandler.RegisterRoutes(basePath, authMiddleware)
	organizationHandler.RegisterRoutes(basePath, authMiddleware)
	billingHandler.RegisterRoutes(basePath, authMiddleware)

	taskRunner.Dispatch()

	slog.Error(router.Run(":8080").Error())
}
