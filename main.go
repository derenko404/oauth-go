package main

import (
	"oauth-go/internal/app"
	"oauth-go/internal/controllers"
	"oauth-go/internal/middleware"
	"oauth-go/internal/types"
	"oauth-go/pkg/configurator"
	"oauth-go/pkg/logger"
	"os"

	docs "oauth-go/docs"

	files "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

// @title           Swagger oauth-go API
// @version         1.0
// @description     This is a sample oauth-go server.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	var config types.AppConfig
	err := configurator.Load(&config, &configurator.Options{})

	logger := logger.NewLogger(config.AppLogLevel)

	if err != nil {
		logger.Error("error loading app configuration", "error", err)
		os.Exit(1)
	}

	app, err := app.New(&config, logger)

	if err != nil {
		logger.Error("error creating app instance", "error", err)
		os.Exit(1)
	}

	authController := controllers.NewAuthController(app)
	healthController := controllers.NewHelathController(app)

	api := app.Router.Group("/api/v1")

	api.GET("/health", healthController.HealthCheck)

	api.GET("/auth/sign-in/:provider", authController.SignIn)
	api.GET("/auth/callback/:provider", authController.HandleCallback)
	api.GET("/auth/me", middleware.AuthMiddleware(app.Store, app.Services, app.Logger), authController.GetMe)
	api.POST("/auth/refresh", authController.RefreshToken)
	api.GET("/auth/sign-out", middleware.AuthMiddleware(app.Store, app.Services, app.Logger), authController.SignOut)

	docs.SwaggerInfo.BasePath = "/api/v1"
	app.Router.GET("/swagger/*any", swagger.WrapHandler(files.Handler))

	if err := app.Start(); err != nil {
		logger.Error("cannot start app", "error", err)
		os.Exit(1)
	}
}
