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
	outletRepo := repository.NewOutletRepository(db)

	// Usecase layer
	authUsecase := usecase.NewAuthUsecase(userRepo)
	outletUsecase := usecase.NewOutletUsecase(outletRepo)

	// Handler layer
	authHandler := handler.NewAuthHandler(authUsecase)
	outletHandler := handler.NewOutletHandler(outletUsecase)

	// === Router Setup ===
	r := gin.Default()

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "LaundryIn API is running! 🚀"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	v1.Use(handler.PayloadLimit(1024 * 1024)) // Limit payload to 1MB
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes — Owner only
		outlets := v1.Group("/outlets")
		outlets.Use(handler.AuthMiddleware(), handler.RoleMiddleware("owner"))
		{
			outlets.POST("", outletHandler.CreateOutlet)
			outlets.GET("", outletHandler.GetAllOutlets)
			outlets.GET("/:id", outletHandler.GetOutletByID)
			outlets.PUT("/:id", outletHandler.UpdateOutlet)
			outlets.DELETE("/:id", outletHandler.DeleteOutlet)
		}
	}

	r.Run() // defaults to :8080
}
