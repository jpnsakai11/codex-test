package usecase

import (
	"context"
	"errors"

	"github.com/example/microservices-project/order-service/internal/domain"
)

var ErrInvalidInput = errors.New("invalid order input")
var ErrUserNotFound = errors.New("user not found")

type OrderRepository interface {
	Create(ctx context.Context, order *domain.Order) error
	GetByID(ctx context.Context, id int64) (*domain.Order, error)
}

type UserClient interface {
	UserExists(ctx context.Context, userID int64) (bool, error)
}

type OrderUsecase struct {
	repo       OrderRepository
	userClient UserClient
}

func NewOrderUsecase(repo OrderRepository, userClient UserClient) *OrderUsecase {
	return &OrderUsecase{repo: repo, userClient: userClient}
}

func (u *OrderUsecase) Create(ctx context.Context, order *domain.Order) error {
	if order.UserID <= 0 || order.Amount <= 0 {
		return ErrInvalidInput
	}
	exists, err := u.userClient.UserExists(ctx, order.UserID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrUserNotFound
	}
	if order.Status == "" {
		order.Status = "created"
	}
	return u.repo.Create(ctx, order)
}

func (u *OrderUsecase) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	if id <= 0 {
		return nil, ErrInvalidInput
	}
	return u.repo.GetByID(ctx, id)
}
