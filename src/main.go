package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/LombardiDaniel/go-gin-template/controllers"
	"github.com/LombardiDaniel/go-gin-template/docs"
	"github.com/LombardiDaniel/go-gin-template/middlewares"
	"github.com/LombardiDaniel/go-gin-template/services"
	"github.com/LombardiDaniel/go-gin-template/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	router *gin.Engine

	usersCol    *mongo.Collection
	sessionsCol *mongo.Collection

	// Services
	authService services.AuthService

	// Controllers
	userController controllers.UserController
	// authController controllers.AuthController

	// Middlewares
	authMiddleware middlewares.AuthMiddleware

	mongoClient *mongo.Client
	ctx         context.Context
)

func init() {
	ctx = context.TODO()

	utils.InitSlogger()

	pgConnStr := utils.GetEnvVarDefault("POSTGRES_URI", "postgres://user:password@localhost:5432/db?sslmode=disable")

	db, err := sql.Open("postgres", pgConnStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Services
	authService = services.NewAuthServiceImpl(db)

	// Middleware
	authMiddleware = middlewares.NewAuthMiddleware(db)

	// Controllers
	// formsController = controllers.NewFormsController(formsService)

	router = gin.Default()
	router.SetTrustedProxies([]string{"*"})

	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AddAllowHeaders("Authorization")

	slog.Info(fmt.Sprintf("corsCfg: %+v\n", corsCfg))

	router.Use(cors.New(corsCfg))

	docs.SwaggerInfo.Title = "Generic Forms API"
	docs.SwaggerInfo.Description = "Generic Forms API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = ""

	if os.Getenv("GIN_MODE") == "release" {
		docs.SwaggerInfo.Host = os.Getenv("SWAGGER_HOST")
		docs.SwaggerInfo.Schemes = []string{"https"}
	} else {
		docs.SwaggerInfo.Host = "localhost:8080"
		docs.SwaggerInfo.Schemes = []string{"http"}
	}

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
// @description "Type 'Bearer $TOKEN' to correctly set the API Key"
func main() {
	defer mongoClient.Disconnect(ctx)

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	basePath := router.Group("/v1")
	formsController.RegisterRoutes(basePath, authMiddleware)

	slog.Error(router.Run(":8080").Error())
}
