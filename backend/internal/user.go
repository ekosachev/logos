package internal

import (
	"errors"

	"github.com/google/uuid"
)

var ErrUsernameExists = errors.New("User with this username already exists")

type User struct {
	ID           uuid.UUID
	Username     string
	PasswordHash string
}
