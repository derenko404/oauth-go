package main

import (
	"go-auth/internal/app"
	"go-auth/internal/controllers"
	"go-auth/internal/middleware"
	"go-auth/internal/types"
	"go-auth/pkg/configurator"
	"go-auth/pkg/logger"
	"os"
)

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

	api := app.Router.Group("/api/v1")

	api.GET("/health", app.HealthCheck)

	api.GET("/sign-in/:provider", authController.SignIn)
	api.GET("/auth/callback/:provider", authController.HandleCallback)
	api.POST("/auth/refresh", authController.RefreshToken)
	api.GET("/auth/me", middleware.AuthMiddleware(app.Store, app.Services), authController.GetMe)
	api.GET("/auth/sign-out", middleware.AuthMiddleware(app.Store, app.Services), authController.SignOut)

	if err := app.Start(); err != nil {
		logger.Error("cannot start app", "error", err)
		os.Exit(1)
	}
}
