package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

	"github.com/LombardiDaniel/gopherbase/common"
	"github.com/LombardiDaniel/gopherbase/controllers"
	"github.com/LombardiDaniel/gopherbase/daemons"
	"github.com/LombardiDaniel/gopherbase/docs"
	"github.com/LombardiDaniel/gopherbase/middlewares"
	"github.com/LombardiDaniel/gopherbase/oauth"
	"github.com/LombardiDaniel/gopherbase/services"
	"github.com/gin-contrib/cors"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	authController         controllers.AuthController
	userController         controllers.UserController
	organizationController controllers.OrganizationController
	billingController      controllers.BillingController

	authMiddleware      middlewares.AuthMiddleware
	telemetryMiddleware middlewares.TelemetryMiddleware

	taskRunner daemons.TaskRunner

	db *sql.DB

	err error
)

func init() {
	common.InitSlogger()
	ctx = context.Background()

	pgConnStr := common.GetEnvVarDefault("POSTGRES_URI", "postgres://user:password@localhost:5432/db?sslmode=disable")
	db, err = sql.Open("postgres", pgConnStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	pgIdleConns, err := strconv.Atoi(common.GetEnvVarDefault("POSTGRES_IDLE_CONNS", "2"))
	if err != nil {
		panic(err)
	}
	pgOpenConns, err := strconv.Atoi(common.GetEnvVarDefault("POSTGRES_OPEN_CONNS", "10"))
	if err != nil {
		panic(err)
	}
	db.SetMaxIdleConns(pgIdleConns)
	db.SetMaxOpenConns(pgOpenConns)
	_, err = db.Exec(fmt.Sprintf("SET TIME ZONE '%s';", common.DefaultTimzone))
	if err != nil {
		panic(err)
	}

	mongoConn := options.Client().ApplyURI(
		common.GetEnvVarDefault("MONGO_URI", "mongodb://localhost:27017"),
	)
	mongoClient, err := mongo.Connect(ctx, mongoConn)
	if err != nil {
		slog.Error(err.Error())
	}
	err = mongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		slog.Error(err.Error())
	}

	tsIdxModel := mongo.IndexModel{
		Keys:    bson.M{"ts": 1},
		Options: options.Index(),
	}

	metricsCol := mongoClient.Database("telemetry").Collection("metrics")
	eventsCol := mongoClient.Database("telemetry").Collection("events")

	_, err = metricsCol.Indexes().CreateOne(ctx, tsIdxModel)
	if err != nil {
		panic(err)
	}
	_, err = eventsCol.Indexes().CreateOne(ctx, tsIdxModel)
	if err != nil {
		panic(err)
	}

	oauthBaseCallback := common.ApiHostUrl + "v1/auth/%s/callback"

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

	s3Host, err := common.ExtractHostFromUrl(common.S3Endpoint)
	if err != nil {
		panic(err)
	}
	s3Secure, err := common.UrlIsSecure(common.S3Endpoint)
	if err != nil {
		panic(err)
	}
	minioClient, err := minio.New(
		s3Host,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				os.Getenv("S3_ACCESS_KEY_ID"),
				os.Getenv("S3_SECRET_ACCESS_KEY"),
				"",
			),
			Region: common.S3Region,
			Secure: s3Secure,
		},
	)
	if err != nil {
		panic(err)
	}

	authService = services.NewAuthServiceJwtImpl(os.Getenv("JWT_SECRET_KEY"), db)
	userService = services.NewUserServicePgImpl(db)
	emailService = services.NewEmailServiceResendImpl(os.Getenv("RESEND_API_KEY"), "./templates")
	organizationService = services.NewOrganizationServicePgImpl(db)
	objectService = services.NewObjectServiceMinioImpl(minioClient)
	billingService = services.NewBillingService(db, os.Getenv("STRIPE_API_KEY"))
	telemetryService = services.NewTelemetryServiceMongoAsyncImpl(mongoClient, metricsCol, eventsCol, 100)

	authMiddleware = middlewares.NewAuthMiddlewareJwt(authService)
	telemetryMiddleware = middlewares.NewTelemetryMiddleware(telemetryService)

	authController = controllers.NewAuthController(authService, userService, emailService, oauthConfigMap)
	userController = controllers.NewUserController(authService, userService, emailService, objectService)
	organizationController = controllers.NewOrganizationController(userService, emailService, organizationService)
	billingController = controllers.NewBillingController(billingService, emailService, userService)

	router = gin.Default()
	router.SetTrustedProxies([]string{"*"})

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{common.ApiHostUrl, common.AppHostUrl}
	corsCfg.AllowCredentials = true
	corsCfg.AddAllowHeaders("Authorization")

	slog.Info(fmt.Sprintf("corsCfg: %+v", corsCfg))

	router.Use(cors.New(corsCfg))
	router.Use(limits.RequestSizeLimiter(common.MaxRequestSize))

	docs.SwaggerInfo.Title = "Generic Forms API"
	docs.SwaggerInfo.Description = "Generic Forms API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Host = strings.Split(common.ApiHostUrl, "://")[1]

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
	taskRunner.RegisterTask(time.Second, func() error {
		return telemetryService.Upload(context.Background())
	}, 1)
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
	authController.RegisterRoutes(basePath, authMiddleware)
	userController.RegisterRoutes(basePath, authMiddleware)
	organizationController.RegisterRoutes(basePath, authMiddleware)
	billingController.RegisterRoutes(basePath, authMiddleware)

	taskRunner.Dispatch()

	slog.Error(router.Run(":8080").Error())
}
