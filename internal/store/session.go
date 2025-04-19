package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionStore interface {
	CreateSession(ctx context.Context, dto *UserSessionDto) (*UserSession, error)
	GetSessionBy(ctx context.Context, filters map[string]any) (*UserSession, error)
	DeleteSessionBy(ctx context.Context, filters map[string]any) error
}

type SessionStoreImpl struct {
	db *pgxpool.Pool
}

type UserSessionDto struct {
	UserID    int
	IPAddress string
	UserAgent string
	Location  string
	DeviceID  string
}

type UserSession struct {
	ID     int   `db:"id" json:"id"`
	UserID int64 `db:"user_id" json:"user_id"`

	IPAddress    string    `db:"ip_address" json:"ip_address"`
	UserAgent    string    `db:"user_agent" json:"user_agent"`
	Location     *string   `db:"location" json:"location,omitempty"`
	DeviceID     string    `db:"device_id" json:"device_id"`
	LastActiveAt time.Time `db:"last_active_at" json:"last_active_at"`

	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

func NewSessionStore(db *pgxpool.Pool) *SessionStoreImpl {
	return &SessionStoreImpl{
		db: db,
	}
}

func (store *SessionStoreImpl) CreateSession(ctx context.Context, dto *UserSessionDto) (*UserSession, error) {
	sql, _, _ := goqu.Insert("user_sessions").
		Rows(goqu.Record{
			"user_id":    dto.UserID,
			"ip_address": dto.IPAddress,
			"location":   dto.Location,
			"user_agent": dto.UserAgent,
			"device_id":  dto.DeviceID,
		}).Returning("*").ToSQL()

	rows, err := store.db.Query(ctx, sql)

	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	session, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[UserSession])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}

		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	return session, nil
}

func (repo *SessionStoreImpl) GetSessionBy(ctx context.Context, filters map[string]any) (*UserSession, error) {
	query := goqu.From("user_sessions")

	for key, value := range filters {
		query = query.Where(goqu.I(key).Eq(value))
	}

	query = query.Where(goqu.I("deleted_at").Is(nil))

	sql, _, _ := query.ToSQL()

	rows, err := repo.db.Query(ctx, sql)

	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	session, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[UserSession])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("session not found: %s", filters)
		}

		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	return session, nil
}

func (repo *SessionStoreImpl) DeleteSessionBy(ctx context.Context, filters map[string]any) error {
	query := goqu.Update("user_sessions").Set(goqu.Record{"deleted_at": time.Now()})

	for key, value := range filters {
		query = query.Where(goqu.I(key).Eq(value))
	}

	query = query.Where(goqu.I("deleted_at").Is(nil))

	sql, _, _ := query.ToSQL()

	_, err := repo.db.Exec(ctx, sql)

	if err != nil {
		return fmt.Errorf("query execution failed: %w", err)
	}

	return nil
}
