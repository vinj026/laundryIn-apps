package main

import (
	"log"
	"os"

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
	authHandler := handler.NewAuthHandler(authUsecase)

	// Phase 3: Outlets
	outletRepo := repository.NewOutletRepository(db)
	outletUsecase := usecase.NewOutletUsecase(outletRepo)
	outletHandler := handler.NewOutletHandler(outletUsecase)

	// === Router Setup ===
	// Set Gin mode
	if mode := os.Getenv("GIN_MODE"); mode != "" {
		gin.SetMode(mode)
	}

	r := gin.Default()

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "LaundryIn API is running! 🚀"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	v1.Use(handler.PayloadLimit(1024 * 1024)) // Limit paylod to 1MB
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			// Brute force protection for login & registration
			auth.Use(handler.RateLimiter())

			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Outlet routes (Phase 3)
		outlets := v1.Group("/outlets")
		outlets.Use(handler.AuthMiddleware(), handler.RoleMiddleware("owner"))
		{
			outlets.POST("", outletHandler.Create)
			outlets.GET("", outletHandler.GetAll)
			outlets.GET("/:id", outletHandler.GetByID)
			outlets.PUT("/:id", outletHandler.Update)
			outlets.DELETE("/:id", outletHandler.Delete)
		}
	}

	r.Run() // defaults to :8080
}
