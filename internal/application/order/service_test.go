package order

import (
	"context"
	"testing"

	domain "api-ddd/internal/domain/order"
)

type FakeOrderRepository struct {
	orders map[int64]*domain.Order
	nextID int64
}

func NewFakeOrderRepository() *FakeOrderRepository {
	return &FakeOrderRepository{
		orders: make(map[int64]*domain.Order),
		nextID: 1,
	}
}

func (f *FakeOrderRepository) Save(
	ctx context.Context,
	o *domain.Order,
) error {
	id := f.nextID
	f.nextID++

	if err := o.AssignID(id); err != nil {
		return err
	}

	f.orders[id] = o
	return nil
}

func (f *FakeOrderRepository) FindByID(
	ctx context.Context,
	id int64,
) (*domain.Order, error) {
	o, exists := f.orders[id]

	if !exists {
		return nil, domain.ErrOrderNotFound
	}

	return o, nil
}

func (f *FakeOrderRepository) Update(
	ctx context.Context,
	o *domain.Order,
) error {
	if _, exists := f.orders[o.ID()]; !exists {
		return domain.ErrOrderNotFound
	}

	f.orders[o.ID()] = o
	return nil
}

func TestCreateOrder(t *testing.T) {
	repository := NewFakeOrderRepository()
	service := NewOrderService(repository)

	dto := CreateOrderDTO{
		Items: []CreateOrderItemDTO{
			{
				ProductID: 1,
				Quantity:  2,
			},
		},
	}

	result, err := service.CreateOrder(context.Background(), dto)

	if err != nil {
		t.Fatalf("error creando pedido: %v", err)
	}

	if result.ID() == 0 {
		t.Error("el pedido debería tener un ID")
	}

	if result.Status() != domain.StatusPending {
		t.Errorf("estado esperado PENDING, obtenido %s", result.Status())
	}
}

func TestCreateEmptyOrder(t *testing.T) {
	repository := NewFakeOrderRepository()
	service := NewOrderService(repository)

	_, err := service.CreateOrder(
		context.Background(),
		CreateOrderDTO{},
	)

	if err != domain.ErrEmptyOrder {
		t.Errorf("se esperaba ErrEmptyOrder, obtenido %v", err)
	}
}

func TestConfirmOrder(t *testing.T) {
	repository := NewFakeOrderRepository()
	service := NewOrderService(repository)

	dto := CreateOrderDTO{
		Items: []CreateOrderItemDTO{
			{ProductID: 1, Quantity: 2},
		},
	}

	created, err := service.CreateOrder(context.Background(), dto)
	if err != nil {
		t.Fatal(err)
	}

	confirmed, err := service.ConfirmOrder(
		context.Background(),
		created.ID(),
	)

	if err != nil {
		t.Fatalf("error confirmando pedido: %v", err)
	}

	if confirmed.Status() != domain.StatusConfirmed {
		t.Errorf("se esperaba CONFIRMED, obtenido %s", confirmed.Status())
	}
}
