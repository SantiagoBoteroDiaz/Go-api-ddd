
package order

import (
	"errors"
	"time"
)


// ORDER STATUS


type OrderStatus string

const (
	StatusPending   OrderStatus = "PENDING"
	StatusConfirmed OrderStatus = "CONFIRMED"
)


// ORDER ITEM


type OrderItem struct {
	productID int64
	quantity  int
}

// Constructor
func NewOrderItem(productID int64, quantity int) (OrderItem, error) {
	if productID <= 0 {
		return OrderItem{}, ErrInvalidProductID
	}

	if quantity <= 0 {
		return OrderItem{}, ErrInvalidQuantity
	}

	return OrderItem{
		productID: productID,
		quantity:  quantity,
	}, nil
}

// Getters
func (i OrderItem) ProductID() int64 {
	return i.productID
}

func (i OrderItem) Quantity() int {
	return i.quantity
}


// ORDER


type Order struct {
	id        int64
	items     []OrderItem
	status    OrderStatus
	createdAt time.Time
}

// Constructores
func NewOrder(items []OrderItem) (*Order, error) {
	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}

	copiedItems := append([]OrderItem(nil), items...)

	return &Order{
		items:     copiedItems,
		status:    StatusPending,
		createdAt: time.Now(),
	}, nil
}

func RehydrateOrder(
	id int64,
	items []OrderItem,
	status OrderStatus,
	createdAt time.Time,
) (*Order, error) {
	if id <= 0 {
		return nil, errors.New("ID de pedido inválido")
	}

	if len(items) == 0 {
		return nil, ErrEmptyOrder
	}

	if status != StatusPending && status != StatusConfirmed {
		return nil, errors.New("estado de pedido inválido")
	}

	return &Order{
		id:        id,
		items:     append([]OrderItem(nil), items...),
		status:    status,
		createdAt: createdAt,
	}, nil
}

// Operaciones de negocio
func (o *Order) Confirm() error {
	if len(o.items) == 0 {
		return ErrEmptyOrder
	}

	if o.status == StatusConfirmed {
		return ErrOrderAlreadyConfirmed
	}

	o.status = StatusConfirmed

	return nil
}

// Operaciones de persistencia
func (o *Order) AssignID(id int64) error {
	if id <= 0 {
		return errors.New("ID de pedido inválido")
	}

	if o.id != 0 {
		return errors.New("el pedido ya tiene un ID")
	}

	o.id = id
	return nil
}

// Getters
func (o *Order) ID() int64 {
	return o.id
}

func (o *Order) Items() []OrderItem {
	return append([]OrderItem(nil), o.items...)
}

func (o *Order) Status() OrderStatus {
	return o.status
}

func (o *Order) CreatedAt() time.Time {
	return o.createdAt
}
