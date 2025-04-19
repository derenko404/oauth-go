package logger

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
)

type loggingTracer struct {
	logger *slog.Logger
}

func NewTracer(logLevel string) *loggingTracer {
	logger := getLogger(logLevel, false)

	return &loggingTracer{
		logger: logger,
	}
}

type contextKey string

const queryDataKey contextKey = "queryData"

type queryData struct {
	SQL       string    `json:"sql"`
	Args      []any     `json:"args"`
	StartTime time.Time `json:"startTime"`
}

func (lt *loggingTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	queryData := queryData{
		SQL:       data.SQL,
		Args:      data.Args,
		StartTime: time.Now(),
	}

	return context.WithValue(ctx, queryDataKey, queryData)
}

func (lt *loggingTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	queryData, ok := ctx.Value(queryDataKey).(queryData)
	if !ok {
		lt.logger.Error("cannot get query data")
		return
	}

	duration := time.Since(queryData.StartTime).Milliseconds()

	if data.Err != nil {
		lt.logger.Info(
			"sql query error",
			"error", data.Err,
			"sql", queryData.SQL,
			"args", queryData.Args,
		)
		return
	}

	lt.logger.Info(
		"sql query completed",
		"sql", queryData.SQL,
		"args", queryData.Args,
		"duration", fmt.Sprintf("%dms", duration),
	)
}
