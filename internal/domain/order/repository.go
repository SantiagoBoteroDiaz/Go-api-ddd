package order

import "context"

type OrderRepository interface {
	Save(ctx context.Context, order *Order) error

	FindByID(ctx context.Context, id int64) (*Order, error)

	Update(ctx context.Context, order *Order) error
}
