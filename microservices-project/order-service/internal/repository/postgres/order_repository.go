package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/example/microservices-project/order-service/internal/domain"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(ctx context.Context, order *domain.Order) error {
	query := `INSERT INTO orders (user_id, amount, status) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, order.UserID, order.Amount, order.Status).Scan(&order.ID)
}

func (r *OrderRepository) GetByID(ctx context.Context, id int64) (*domain.Order, error) {
	order := &domain.Order{}
	query := `SELECT id, user_id, amount, status FROM orders WHERE id = $1`
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&order.ID, &order.UserID, &order.Amount, &order.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return order, nil
}
