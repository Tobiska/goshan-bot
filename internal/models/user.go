package models

import "time"

type User struct {
	ID        int64
	Username  string
	ChatID    int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
