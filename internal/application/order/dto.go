package order

type CreateOrderItemDTO struct {
	ProductID int64
	Quantity  int
}

type CreateOrderDTO struct {
	Items []CreateOrderItemDTO
}
