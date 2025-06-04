package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/LombardiDaniel/goliath/src/internal/di"
	"github.com/LombardiDaniel/goliath/src/pkg/constants"
	"github.com/LombardiDaniel/goliath/src/pkg/it"
	"github.com/LombardiDaniel/goliath/src/pkg/logger"
	_ "github.com/lib/pq"

	"github.com/gin-contrib/cors"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"

	"github.com/LombardiDaniel/goliath/src/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	router *gin.Engine

	container *di.Container

	port *string
	env  *string
)

func init() {
	logger.InitSlogger()

	port = flag.String("p", "8080", "Port to run the server on")
	env = flag.String("e", "dev", "Sets the env as dev/prod")
	flag.Parse()

	allowedEnvs := map[di.Env]bool{
		di.DevEnv:  true,
		di.ProdEnv: true,
	}

	envVal := di.Env(*env)
	if !allowedEnvs[envVal] {
		panic(fmt.Sprintf("invalid environment: %s", envVal))
	}

	container = it.Must(di.NewContainer(envVal))

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

	container.Setup(router)
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
	container.DispatchAsync()

	go func() {
		if err := router.Run(":" + *port); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()
	slog.Info("Server running on port " + *port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")
}
