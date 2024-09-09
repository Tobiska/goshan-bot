package user

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/sqlite3"

	"goshan-bot/internal/clients/sqlite"
	"goshan-bot/internal/models"
)

const (
	notificationsTableName      = "notifications"
	notificationBuildsTableName = "notifications_build"
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

func (r *Repository) CreateBuild(ctx context.Context, build models.NotificationBuild) error {
	query := driver.Insert(notificationBuildsTableName).Rows(goqu.Record{
		"chat_id": build.ChatID,
		"user_id": build.UserID,
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

func (r *Repository) UpdateBuild(ctx context.Context, chatID, userID int64, build models.NotificationBuild) error {
	updateRec := goqu.Record{}

	if build.Tag != nil {
		updateRec["tag"] = build.Tag
	}

	if build.Description != nil {
		updateRec["description"] = build.Description
	}

	if build.RemindIn != nil {
		updateRec["notify_at"] = build.RemindIn
	}

	if build.EventAt != nil {
		updateRec["event_at"] = build.EventAt
	}

	query := driver.Update(notificationBuildsTableName).Where(goqu.Ex{"chat_id": chatID, "user_id": userID}).
		Set(updateRec)

	sqlQuery, args, err := query.ToSQL()
	if err != nil {
		return fmt.Errorf("build to sql error: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, sqlQuery, args...); err != nil {
		return fmt.Errorf("error while execution query: %w", err)
	}

	return nil
}

func (r *Repository) CreateNotification(ctx context.Context, notification *models.Notification) error {
	query := driver.Insert(notificationsTableName).Rows(goqu.Record{
		"chat_id":     notification.ChatID,
		"user_id":     notification.UserID,
		"tag":         notification.Tag,
		"description": notification.Description,
		"event_at":    notification.EventAt,
		"notify_at":   notification.NotifyAt,
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

func (r *Repository) FindBuildByUserID(ctx context.Context, userID int64) (*models.NotificationBuild, error) {
	query := driver.Select(
		goqu.C("user_id"),
		goqu.C("chat_id"),
		goqu.C("tag"),
		goqu.C("description"),
		goqu.C("notify_at"),
		goqu.C("event_at"),
		goqu.C("created_at"),
		goqu.C("updated_at"),
	).From(notificationBuildsTableName).Where(goqu.Ex{"user_id": userID}).Prepared(true)

	sqlQuery, args, err := query.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("sql build error: %w", err)
	}

	b := &models.NotificationBuild{}

	if err := r.db.QueryRowContext(ctx, sqlQuery, args...).Scan(&b.UserID, &b.ChatID,
		&b.Tag, &b.Description,
		&b.RemindIn, &b.EventAt,
		&b.CreatedAt, &b.UpdatedAt); err != nil {
		return nil, fmt.Errorf("execution error: %w", err)
	}

	return b, nil
}

func (r *Repository) DeleteBuild(ctx context.Context, chatID, userID int64) error {
	query := driver.Delete(notificationBuildsTableName).Where(goqu.Ex{"user_id": userID, "chat_id": chatID}).Prepared(true)

	sqlQuery, args, err := query.ToSQL()
	if err != nil {
		return fmt.Errorf("sql build error: %w", err)
	}

	if _, err := r.db.ExecContext(ctx, sqlQuery, args...); err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	return nil
}
