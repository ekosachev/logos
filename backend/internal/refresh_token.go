package internal

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNoActiveToken = errors.New("User does not have an active refresh token")
	ErrTokenNotFound = errors.New("Token was not found")
)

type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Token     string
	Revoked   bool
	ExpiresAt time.Time

	CreatedAt time.Time
}
