package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"

	"goshan-bot/internal/clients/sqlite"
	"goshan-bot/internal/models"
)

const (
	userTableName = "users"
)

var (
	driver = goqu.Dialect("sqlite3")
)

type Repository struct {
	db sqlite.SQLDatabase
}

func New(db sqlite.SQLDatabase) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Save(ctx context.Context, u models.User) error {
	query := driver.Insert(userTableName).Rows(goqu.Record{
		"username": u.Username,
		"chat_id":  u.ChatID,
		"user_id":  u.ID,
	}).Prepared(true)

	sqlQuery, args, err := query.ToSQL()
	if err != nil {
		return fmt.Errorf("sql build error: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, sqlQuery, args...); err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	return nil
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*models.User, error) {
	query := driver.Select(
		goqu.C("user_id"),
		goqu.C("chat_id"),
		goqu.C("username"),
		goqu.C("created_at"),
		goqu.C("updated_at"),
	).From(userTableName).Where(goqu.Ex{"user_id": id}).Prepared(true)

	sqlQuery, args, err := query.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("sql build error: %w", err)
	}

	u := &models.User{}

	if err := r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&u.ID, &u.ChatID, &u.Username, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("execution error: %w", err)
	}

	return u, nil
}
