package auth

import "time"

type AuthToken struct {
	ID        string
	UserID    string
	Token     string
	ExpiredAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
