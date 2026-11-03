package order

import (
	"context"

	domain "api-ddd/internal/domain/order"
)

type OrderService struct {
	repository domain.OrderRepository
}

func NewOrderService(repository domain.OrderRepository) *OrderService {
	return &OrderService{
		repository: repository,
	}
}

// CreateOrder crea y guarda un pedido.
func (s *OrderService) CreateOrder(
	ctx context.Context,
	dto CreateOrderDTO,
) (*domain.Order, error) {

	var items []domain.OrderItem

	for _, itemDTO := range dto.Items {
		item, err := domain.NewOrderItem(
			itemDTO.ProductID,
			itemDTO.Quantity,
		)

		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	newOrder, err := domain.NewOrder(items)
	if err != nil {
		return nil, err
	}

	err = s.repository.Save(ctx, newOrder)
	if err != nil {
		return nil, err
	}

	return newOrder, nil
}

// GetOrderByID consulta un pedido.
func (s *OrderService) GetOrderByID(
	ctx context.Context,
	id int64,
) (*domain.Order, error) {

	return s.repository.FindByID(ctx, id)
}

// ConfirmOrder confirma un pedido.
func (s *OrderService) ConfirmOrder(
	ctx context.Context,
	id int64,
) (*domain.Order, error) {

	existingOrder, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := existingOrder.Confirm(); err != nil {
		return nil, err
	}

	if err := s.repository.Update(ctx, existingOrder); err != nil {
		return nil, err
	}

	return existingOrder, nil
}
