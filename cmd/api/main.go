package main

import (
	"log"

	"laundryin/internal/database"
	handler "laundryin/internal/delivery/http"
	"laundryin/internal/repository"
	"laundryin/internal/repository/models"
	"laundryin/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to database
	db := database.ConnectDB()

	// Auto-migrate models
	db.AutoMigrate(&models.User{}, &models.Outlet{}, &models.Service{}, &models.Order{})

	// === Dependency Injection ===
	// Repository layer
	userRepo := repository.NewUserRepository(db)

	// Usecase layer
	authUsecase := usecase.NewAuthUsecase(userRepo)

	// Handler layer
	authHandler := handler.NewAuthHandler(authUsecase)

	// === Router Setup ===
	r := gin.Default()

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "LaundryIn API is running! 🚀"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	v1.Use(handler.PayloadLimit(1024 * 1024)) // Limit paylod to 1MB
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes (Phase 3+) — ready for future use
		// protected := v1.Group("/")
		// protected.Use(handler.AuthMiddleware())
		// {
		//     // Outlet, Service, Order routes will go here
		// }
	}

	r.Run() // defaults to :8080
}
