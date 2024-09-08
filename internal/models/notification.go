package models

import "time"

type Notification struct {
	ID          int64
	UserID      int64
	ChatID      int64
	Tag         string
	Description string
	NotifyAt    time.Time
	EventAt     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func CreateFromBuild(build *NotificationBuild) *Notification {
	if build == nil {
		return nil
	}

	return &Notification{
		UserID:   build.UserID,
		ChatID:   build.ChatID,
		Tag:      *build.Tag,
		EventAt:  *build.EventAt,
		NotifyAt: build.EventAt.Add(-1 * *build.RemindIn),
	}
}

type NotificationBuild struct {
	UserID      int64
	ChatID      int64
	Tag         *string
	Description *string
	RemindIn    *time.Duration
	EventAt     *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
