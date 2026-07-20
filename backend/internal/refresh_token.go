package internal

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Token   string
	Revoked bool

	CreatedAt time.Time
}
