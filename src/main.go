package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/LombardiDaniel/go-gin-template/common"
	"github.com/LombardiDaniel/go-gin-template/controllers"
	"github.com/LombardiDaniel/go-gin-template/docs"
	"github.com/LombardiDaniel/go-gin-template/middlewares"
	"github.com/LombardiDaniel/go-gin-template/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	router *gin.Engine

	// Services
	authService         services.AuthService
	userService         services.UserService
	emailService        services.EmailService
	organizationService services.OrganizationService

	// Controllers
	authController         controllers.AuthController
	userController         controllers.UserController
	organizationController controllers.OrganizationController

	// Middlewares
	authMiddleware middlewares.AuthMiddleware

	db *sql.DB

	err error
)

func init() {
	common.InitSlogger()

	pgConnStr := common.GetEnvVarDefault("POSTGRES_URI", "postgres://user:password@localhost:5432/db?sslmode=disable")

	db, err = sql.Open("postgres", pgConnStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Services
	authService = services.NewAuthServiceJwtImpl(os.Getenv("JWT_SECRET_KEY"))
	userService = services.NewUserServicePgImpl(db)
	emailService = services.NewEmailServiceResendImpl(os.Getenv("RESEND_API_KEY"), "./templates")
	organizationService = services.NewOrganizationServicePgImpl(db)

	// Middleware
	authMiddleware = middlewares.NewAuthMiddlewareJwt(authService)

	// Controllers
	authController = controllers.NewAuthController(authService, userService)
	userController = controllers.NewUserController(userService, emailService)
	organizationController = controllers.NewOrganizationController(userService, emailService, organizationService)

	router = gin.Default()
	router.SetTrustedProxies([]string{"*"})

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowOrigins = []string{common.API_HOST_URL, common.APP_HOST_URL}
	corsCfg.AllowCredentials = true
	corsCfg.AddAllowHeaders("Authorization")

	fmt.Printf("len(corsCfg.AllowOrigins): %v\n", len(corsCfg.AllowOrigins))

	slog.Info(fmt.Sprintf("corsCfg: %+v", corsCfg))

	router.Use(cors.New(corsCfg))

	docs.SwaggerInfo.Title = "Generic Forms API"
	docs.SwaggerInfo.Description = "Generic Forms API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Host = strings.Split(common.API_HOST_URL, "://")[1]

	if os.Getenv("GIN_MODE") == "release" {
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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

	basePath := router.Group("/v1")
	authController.RegisterRoutes(basePath, authMiddleware)
	userController.RegisterRoutes(basePath, authMiddleware)
	organizationController.RegisterRoutes(basePath, authMiddleware)

	slog.Error(router.Run(":8080").Error())
}
