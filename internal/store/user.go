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

type UserStore interface {
	CreateUser(ctx context.Context, dto *UserDto) (*User, error)
	GetUserBy(ctx context.Context, filters map[string]any) (*User, error)
}

type userStore struct {
	db *pgxpool.Pool
}

type User struct {
	ID              int        `db:"id" json:"id"`
	Name            *string    `db:"name" json:"name,omitempty"`
	Email           string     `db:"email" json:"email"`
	AvatarURL       *string    `db:"avatar_url" json:"avatar_url,omitempty"`
	IsEmailVerified bool       `db:"is_email_verified" json:"is_email_verified"`
	Provider        string     `db:"provider" json:"provider"`
	ProviderUserID  string     `db:"provider_user_id" json:"-"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"-"`
	DeletedAt       *time.Time `db:"deleted_at" json:"-"`
}

type UserDto struct {
	Name           string `db:"name" json:"name,omitempty"`
	Email          string `db:"email" json:"email"`
	AvatarURL      string `db:"avatar_url" json:"avatar_url,omitempty"`
	Provider       string
	ProviderUserID string
}

func NewUserStore(db *pgxpool.Pool) *userStore {
	return &userStore{
		db: db,
	}
}

func (store *userStore) CreateUser(ctx context.Context, dto *UserDto) (*User, error) {
	sql, _, _ := goqu.Insert("users").
		Rows(goqu.Record{
			"name":              dto.Name,
			"email":             dto.Email,
			"avatar_url":        dto.AvatarURL,
			"is_email_verified": true,
			"provider":          dto.Provider,
			"provider_user_id":  dto.ProviderUserID,
		}).Returning("*").ToSQL()

	rows, err := store.db.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[User])
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	return user, nil
}

func (store *userStore) GetUserBy(ctx context.Context, filters map[string]any) (*User, error) {
	query := goqu.From("users")

	for key, value := range filters {
		query = query.Where(goqu.I(key).Eq(value))
	}

	query = query.Where(goqu.I("deleted_at").Is(nil))

	sql, _, _ := query.ToSQL()

	rows, err := store.db.Query(ctx, sql)

	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByPos[User])

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("user not found: %s", filters)
		}

		return nil, fmt.Errorf("query execution failed: %w", err)
	}

	return user, nil
}
