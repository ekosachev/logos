package repositories

import (
	"context"
	"errors"

	"github.com/ekosachev/logos/internal"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string    `gorm:"uniqueIndex"`
	PasswordHash string

	RefreshTokens []RefreshToken
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) StoreUser(ctx context.Context, user *internal.User) error {
	model := User{Username: user.Username, PasswordHash: user.PasswordHash}
	err := gorm.G[User](r.db).Create(ctx, &model)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return internal.ErrUsernameExists
		}

		return err
	}

	user.ID = model.ID
	return nil
}
