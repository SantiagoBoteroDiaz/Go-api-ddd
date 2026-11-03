package order

import (
	"errors"
	"testing"
)

func TestNewOrder(t *testing.T) {
	item, err := NewOrderItem(1, 2)
	if err != nil {
		t.Fatalf("no se esperaba un error: %v", err)
	}

	tests := []struct {
		name    string
		items   []OrderItem
		wantErr error
	}{
		{
			name:  "crear pedido correctamente",
			items: []OrderItem{item},
		},
		{
			name:    "crear pedido vacío",
			items:   []OrderItem{},
			wantErr: ErrEmptyOrder,
		},
		{
			name:    "crear pedido con nil",
			items:   nil,
			wantErr: ErrEmptyOrder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewOrder(tt.items)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("se esperaba %v, se obtuvo %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("no se esperaba un error: %v", err)
			}

			if result.Status() != StatusPending {
				t.Errorf("se esperaba PENDING, se obtuvo %s", result.Status())
			}

			if len(result.Items()) == 0 {
				t.Error("el pedido debería contener productos")
			}
		})
	}
}

func TestNewOrderItem(t *testing.T) {
	tests := []struct {
		name      string
		productID int64
		quantity  int
		wantErr   error
	}{
		{
			name:      "producto válido",
			productID: 1,
			quantity:  2,
		},
		{
			name:      "ID inválido",
			productID: -1,
			quantity:  2,
			wantErr:   ErrInvalidProductID,
		},
		{
			name:      "cantidad inválida",
			productID: 1,
			quantity:  0,
			wantErr:   ErrInvalidQuantity,
		},
		{
			name:      "cantidad negativa",
			productID: 1,
			quantity:  -2,
			wantErr:   ErrInvalidQuantity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewOrderItem(tt.productID, tt.quantity)

			if !errors.Is(err, tt.wantErr) {
				if tt.wantErr == nil && err == nil {
					return
				}

				t.Fatalf("error esperado: %v, obtenido: %v", tt.wantErr, err)
			}
		})
	}
}

func TestConfirmOrder(t *testing.T) {
	item, _ := NewOrderItem(1, 2)

	order, err := NewOrder([]OrderItem{item})
	if err != nil {
		t.Fatalf("error creando pedido: %v", err)
	}

	err = order.Confirm()
	if err != nil {
		t.Fatalf("no se esperaba un error al confirmar: %v", err)
	}

	if order.Status() != StatusConfirmed {
		t.Errorf("se esperaba CONFIRMED, se obtuvo %s", order.Status())
	}
}

func TestConfirmOrderTwice(t *testing.T) {
	item, _ := NewOrderItem(1, 2)

	order, _ := NewOrder([]OrderItem{item})

	if err := order.Confirm(); err != nil {
		t.Fatalf("error en primera confirmación: %v", err)
	}

	err := order.Confirm()

	if !errors.Is(err, ErrOrderAlreadyConfirmed) {
		t.Fatalf("se esperaba ErrOrderAlreadyConfirmed, se obtuvo: %v", err)
	}
}
