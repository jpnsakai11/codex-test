package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/example/microservices-project/user-service/internal/domain"
)

var ErrInvalidInput = errors.New("invalid user input")

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id int64) (*domain.User, error)
}

type UserUsecase struct {
	repo UserRepository
}

func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{repo: repo}
}

func (u *UserUsecase) Create(ctx context.Context, user *domain.User) error {
	if strings.TrimSpace(user.Name) == "" || !strings.Contains(user.Email, "@") {
		return ErrInvalidInput
	}
	return u.repo.Create(ctx, user)
}

func (u *UserUsecase) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}
	return u.repo.GetByID(ctx, id)
}
