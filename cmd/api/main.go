package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"api-ddd/internal/application/order"
	httpDelivery "api-ddd/internal/delivery/http"
	"api-ddd/internal/delivery/http/handlers"
	"api-ddd/internal/infrastructure/persistence"
)

func main() {
	// Cargar variables del archivo .env
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró .env; se usarán las variables del entorno")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL no está configurada")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	db, err := persistence.NewPostgresPool(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Error conectando a PostgreSQL: %v", err)
	}
	defer db.Close()

	orderRepository := persistence.NewPostgresOrderRepository(db)
	orderService := order.NewOrderService(orderRepository)
	orderHandler := handlers.NewOrderHandler(orderService)

	router := gin.Default()

	httpDelivery.RegisterRoutes(router, orderHandler)

	log.Printf("Servidor ejecutándose en http://localhost:%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}