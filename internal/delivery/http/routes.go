package http

import (
    "github.com/gin-gonic/gin"

    "api-ddd/internal/delivery/http/handlers"
    "api-ddd/internal/delivery/http/middleware"
)

func RegisterRoutes(
    router *gin.Engine,
    orderHandler *handlers.OrderHandler,
) {
    authHandler := handlers.NewAuthHandler()

    api := router.Group("/api/v1")
    {
        // Rutas públicas
        api.POST("/auth/token", authHandler.GenerateToken)

        // Rutas protegidas
        protected := api.Group("")
        protected.Use(middleware.JWTAuth())

        protected.POST("/orders", orderHandler.CreateOrder)
        protected.GET("/orders/:id", orderHandler.GetOrderByID)
        protected.GET("/orders", orderHandler.GetAllOrders)
        protected.POST("/orders/:id/confirm", orderHandler.ConfirmOrder)
    }
}