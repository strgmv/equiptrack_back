package models

import (
	"github.com/google/uuid"
)

type Session struct {
	SessionID    int       `db:"id"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
}
