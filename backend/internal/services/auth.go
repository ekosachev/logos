package services

import (
	"context"
	"errors"
	"time"

	"github.com/ekosachev/logos/internal"
	"github.com/ekosachev/logos/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RefreshTokenStorer interface {
	StoreRefreshToken(ctx context.Context, token *internal.RefreshToken) error
}

type AuthService struct {
	refreshTokenRepository RefreshTokenStorer
	userRepository         UserStorer
}

func NewAuthService(
	refreshTokenRepository RefreshTokenStorer,
	userRepository UserStorer,
) *AuthService {
	return &AuthService{
		refreshTokenRepository: refreshTokenRepository,
		userRepository:         userRepository,
	}
}

func (s *AuthService) Login(
	ctx context.Context,
	username string,
	password string,
) (accessToken, refreshToken string, err error) {
	user, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, internal.ErrUserNotFound) {
			return "", "", internal.ErrAuthFailed
		}

		return "", "", err
	}

	if err = validatePassword(user.PasswordHash, password); err != nil {
		return "", "", internal.ErrAuthFailed
	}

	accessToken, err = generateAccessToken(user.ID)
	if err != nil {
		return "", "", internal.ErrAuthFailed
	}
	exp, refreshToken, err := generateRefreshToken(user.ID)
	if err != nil {
		return "", "", internal.ErrAuthFailed
	}

	token := internal.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: exp,
	}
	if err = s.refreshTokenRepository.StoreRefreshToken(ctx, &token); err != nil {
		return "", "", internal.ErrAuthFailed
	}

	return accessToken, refreshToken, nil
}

func validatePassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func generateAccessToken(userID uuid.UUID) (string, error) {
	cfg := config.GetConfig()
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Duration(cfg.JWTAccessExpiration) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTAccessSecret))
}

func generateRefreshToken(userID uuid.UUID) (time.Time, string, error) {
	cfg := config.GetConfig()
	exp := time.Now().Add(time.Duration(cfg.JWTRefreshExpiration) * time.Second)

	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWTRefreshSecret))
	return exp, tokenString, err
}
