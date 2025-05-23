package app

import (
	"fmt"
	"log/slog"
	"net"
	"oauth-go/internal/services"
	"oauth-go/internal/store"
	"oauth-go/internal/types"
	"oauth-go/pkg/database"

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
	RDB      *redis.Client
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

	app.RDB = redis.NewClient(&redis.Options{
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
