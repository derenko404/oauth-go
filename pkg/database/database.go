package database

import (
	"context"
	"fmt"
	"oauth-go/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(user, password, host, port, name string, logLevel string) (*pgxpool.Pool, error) {
	connection := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user, password, host, port, name,
	)

	config, err := pgxpool.ParseConfig(connection)

	if err != nil {
		return nil, err
	}

	config.MaxConns = 10
	config.MinConns = 2

	tracer := logger.NewTracer(logLevel)
	config.ConnConfig.Tracer = tracer

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
