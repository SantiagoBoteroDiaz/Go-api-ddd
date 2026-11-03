package http

import (
	"github.com/gin-gonic/gin"

	"api-ddd/internal/delivery/http/handlers"
)

func RegisterRoutes(
	router *gin.Engine,
	orderHandler *handlers.OrderHandler,
) {
	api := router.Group("/api/v1")
	{
		api.POST("/orders", orderHandler.CreateOrder)
		api.GET("/orders/:id", orderHandler.GetOrderByID)
		api.POST("/orders/:id/confirm", orderHandler.ConfirmOrder)
	}
}
