package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	appOrder "api-ddd/internal/application/order"
	domain "api-ddd/internal/domain/order"
)

type OrderHandler struct {
	service *appOrder.OrderService
}

func NewOrderHandler(service *appOrder.OrderService) *OrderHandler {
	return &OrderHandler{
		service: service,
	}
}

type CreateOrderRequest struct {
	Items []CreateOrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

type CreateOrderItemRequest struct {
	ProductID int64 `json:"product_id" binding:"required,gt=0"`
	Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

type OrderResponse struct {
	ID        int64               `json:"id"`
	Status    domain.OrderStatus  `json:"status"`
	Items     []OrderItemResponse `json:"items"`
	CreatedAt time.Time           `json:"created_at"`
}

type OrderItemResponse struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

func toOrderResponse(o *domain.Order) OrderResponse {
	items := o.Items()
	responseItems := make([]OrderItemResponse, 0, len(items))

	for _, item := range items {
		responseItems = append(responseItems, OrderItemResponse{
			ProductID: item.ProductID(),
			Quantity:  item.Quantity(),
		})
	}

	return OrderResponse{
		ID:        o.ID(),
		Status:    o.Status(),
		Items:     responseItems,
		CreatedAt: o.CreatedAt(),
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var request CreateOrderRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"error": gin.H{
				"code":    "INVALID_REQUEST",
				"message": "El formato de entrada no es válido",
			},
		})
		return
	}

	dto := appOrder.CreateOrderDTO{}

	for _, item := range request.Items {
		dto.Items = append(dto.Items, appOrder.CreateOrderItemDTO{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	createdOrder, err := h.service.CreateOrder(c.Request.Context(), dto)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(201, toOrderResponse(createdOrder))
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "El ID del pedido no es válido",
			},
		})
		return
	}

	result, err := h.service.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(200, toOrderResponse(result))
}

func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{
			"error": gin.H{
				"code":    "INVALID_ID",
				"message": "El ID del pedido no es válido",
			},
		})
		return
	}

	result, err := h.service.ConfirmOrder(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(200, toOrderResponse(result))
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrOrderNotFound):
		c.JSON(404, gin.H{
			"error": gin.H{
				"code":    "ORDER_NOT_FOUND",
				"message": "El pedido no existe",
			},
		})

	case errors.Is(err, domain.ErrEmptyOrder),
		errors.Is(err, domain.ErrInvalidProductID),
		errors.Is(err, domain.ErrInvalidQuantity):
		c.JSON(400, gin.H{
			"error": gin.H{
				"code":    "INVALID_ORDER",
				"message": err.Error(),
			},
		})

	case errors.Is(err, domain.ErrOrderAlreadyConfirmed):
		c.JSON(409, gin.H{
			"error": gin.H{
				"code":    "ORDER_ALREADY_CONFIRMED",
				"message": "El pedido ya fue confirmado",
			},
		})

	default:
		c.JSON(500, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Ocurrió un error interno",
			},
		})
	}
}

func (h *OrderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.service.GetAllOrders(c.Request.Context())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error" : "No se pudieron consultar los pedidos", 
		})
		return
	}

	response := make([]gin.H, 0, len(orders))

	for _, o := range orders { 
		items := make([]gin.H, 0)

		for _, item := range o.Items() {
            items = append(items, gin.H{
                "product_id": item.ProductID(),
                "quantity":   item.Quantity(),
            })
        }

		response = append(response, gin.H{
			"id":         o.ID(),
            "status":     o.Status(),
            "created_at": o.CreatedAt(),
            "items":      items,
		})
	}
	
	c.JSON(http.StatusOK, response) 
}