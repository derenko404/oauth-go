package app

import (
	"fmt"
	"go-auth/internal/services"
	"go-auth/internal/store"
	"go-auth/internal/types"
	"go-auth/pkg/database"
	"log/slog"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type App struct {
	Config   *types.AppConfig
	DB       *pgxpool.Pool
	Logger   *slog.Logger
	Router   *gin.Engine
	Store    *store.Store
	Services *services.Services
	rdb      *redis.Client
}

func New(config *types.AppConfig, logger *slog.Logger) (*App, error) {
	app := &App{
		Config: config,
		Router: gin.Default(),
		Logger: logger,
	}

	db, err := database.Connect(
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
		config.AppLogLevel,
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	app.DB = db

	app.rdb = redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(app.Config.RedisHost, app.Config.RedisPort),
		Password: "",
		DB:       0,
	})

	app.Store = &store.Store{
		User:    store.NewUserStore(app.DB),
		Session: store.NewSessionStore(app.DB),
	}

	app.Services = services.New(app.Config)

	return app, nil
}

func (app *App) Start() error {
	return app.Router.Run(net.JoinHostPort(app.Config.AppHost, app.Config.AppPort))
}

func (app *App) HealthCheck(c *gin.Context) {
	err := app.DB.Ping(c.Request.Context())
	if err != nil {
		app.Logger.Error("db ping error", "error", err)

		c.JSON(500, gin.H{
			"status": "db error",
		})
		return
	}

	_, err = app.rdb.Ping(c.Request.Context()).Result()

	if err != nil {
		app.Logger.Error("redis ping error", "error", err)

		c.JSON(500, gin.H{
			"status": "redis error",
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "ok",
	})
}
