package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/ekosachev/logos/internal"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID
	Token     string
	Revoked   bool
	ExpiresAt time.Time

	CreatedAt time.Time
}

type RefreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) StoreRefreshToken(ctx context.Context, token *internal.RefreshToken) error {
	model := RefreshToken{
		UserID:  token.UserID,
		Token:   token.Token,
		Revoked: false,
	}

	err := gorm.G[RefreshToken](r.db).Create(ctx, &model)
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return internal.ErrUserNotFound
		}

		return err
	}

	token.ID = model.ID
	token.CreatedAt = model.CreatedAt
	return nil
}

func (r *RefreshTokenRepository) GetActiveRefreshToken(ctx context.Context, userID uuid.UUID) (*internal.RefreshToken, error) {
	token, err := gorm.G[RefreshToken](r.db).
		Where("user_id = ? AND revoked = false AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, internal.ErrNoActiveToken
		}

		return nil, err
	}

	return &internal.RefreshToken{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		Revoked:   token.Revoked,
		ExpiresAt: token.ExpiresAt,
		CreatedAt: token.CreatedAt,
	}, nil
}

func (r *RefreshTokenRepository) RevokeToken(ctx context.Context, tokenID uuid.UUID) error {
	rowsUpdated, err := gorm.G[RefreshToken](r.db).Where("id = ?", tokenID).Update(ctx, "revoked", false)
	if rowsUpdated == 0 {
		return internal.ErrTokenNotFound
	}

	return err
}
