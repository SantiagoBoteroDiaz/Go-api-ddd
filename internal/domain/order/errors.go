package order

import "errors"

var (
	ErrEmptyOrder = errors.New("el pedido debe contener al menos un producto")

	ErrInvalidProductID = errors.New("el ID del producto no es válido")

	ErrInvalidQuantity = errors.New("la cantidad debe ser mayor a cero")

	ErrOrderAlreadyConfirmed = errors.New("el pedido ya fue confirmado")

	ErrOrderNotFound = errors.New("pedido no encontrado")
)
