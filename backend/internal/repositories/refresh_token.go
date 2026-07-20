package repositories

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID  uuid.UUID
	Token   string
	Revoked bool

	CreatedAt time.Time
}
