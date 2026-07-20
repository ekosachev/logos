package services

import (
	"context"

	"github.com/ekosachev/logos/internal"
	"golang.org/x/crypto/bcrypt"
)

type UserStorer interface {
	StoreUser(ctx context.Context, user *internal.User) error
}

type UserService struct {
	repository UserStorer
}

func NewUserService(repository UserStorer) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) CreateUser(ctx context.Context, user *internal.User) error {
	hashedPassword, err := hashPassword(user.PasswordHash)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.repository.StoreUser(ctx, user)
}

func hashPassword(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(passwordBytes), nil
}
