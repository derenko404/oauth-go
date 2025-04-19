package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type loggingTracer struct{}

func (lt *loggingTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	slog.Debug("sql query start", "sql", data.SQL)
	return ctx
}

func (lt *loggingTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		slog.Debug("sql query error", "error", data.Err)
	} else {
		slog.Debug("sql query end", "affected", data.CommandTag.RowsAffected())
	}
}

func Connect(user, password, host, port, name string) (*pgxpool.Pool, error) {
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

	config.ConnConfig.Tracer = &loggingTracer{}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)

	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
