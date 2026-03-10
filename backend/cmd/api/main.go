package main

import (
	"log"

	"laundryin/internal/database"
	handler "laundryin/internal/delivery/http"
	"laundryin/internal/repository"
	"laundryin/internal/repository/models"
	"laundryin/internal/usecase"
	"laundryin/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Register custom validators (must be before any request handling)
	utils.RegisterCustomValidators()

	// Connect to database
	db := database.ConnectDB()

	// Auto-migrate models
	db.AutoMigrate(&models.User{}, &models.Outlet{}, &models.Service{}, &models.Order{}, &models.OrderItem{}, &models.OrderLog{})

	// === Dependency Injection ===
	// Repository layer
	userRepo := repository.NewUserRepository(db)
	outletRepo := repository.NewOutletRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	reportRepo := repository.NewReportRepository(db)

	// Usecase layer
	authUsecase := usecase.NewAuthUsecase(userRepo)
	outletUsecase := usecase.NewOutletUsecase(outletRepo)
	serviceUsecase := usecase.NewServiceUsecase(serviceRepo, outletRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, serviceRepo, outletRepo)
	reportUsecase := usecase.NewReportUsecase(reportRepo)

	// Handler layer
	authHandler := handler.NewAuthHandler(authUsecase)
	outletHandler := handler.NewOutletHandler(outletUsecase)
	serviceHandler := handler.NewServiceHandler(serviceUsecase)
	orderHandler := handler.NewOrderHandler(orderUsecase)
	reportHandler := handler.NewReportHandler(reportUsecase)

	// === Router Setup ===
	r := gin.Default()

	// Global security middleware
	r.Use(handler.SecurityHeaders())

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "LaundryIn API is running! 🚀"})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	v1.Use(handler.PayloadLimit(1024 * 1024)) // Limit payload to 1MB
	v1.Use(handler.RateLimiter())             // 60 req/min per IP
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes — Owner only
		protected := v1.Group("")
		protected.Use(handler.AuthMiddleware(), handler.RoleMiddleware("owner"))
		{
			// Outlets
			protected.POST("/outlets", outletHandler.CreateOutlet)
			protected.GET("/outlets", outletHandler.GetAllOutlets)
			protected.GET("/outlets/:id", outletHandler.GetOutletByID)
			protected.PUT("/outlets/:id", outletHandler.UpdateOutlet)
			protected.DELETE("/outlets/:id", outletHandler.DeleteOutlet)

			// Services
			protected.POST("/services", serviceHandler.CreateService)
			protected.GET("/outlets/:id/services", serviceHandler.GetAllByOutletID)
			protected.PUT("/services/:id", serviceHandler.UpdateService)
			protected.DELETE("/services/:id", serviceHandler.DeleteService)

			// Orders (Owner specific)
			protected.POST("/orders", orderHandler.CreateOrder)
			protected.GET("/orders", orderHandler.GetAllByUserID)
			protected.GET("/outlets/:id/orders", orderHandler.GetAllByOutletID)
			protected.PATCH("/orders/:id/status", orderHandler.UpdateStatus)

			// Reports & Analytics
			protected.GET("/reports/omzet", reportHandler.GetOmzet)
			protected.GET("/reports/orders/summary", reportHandler.GetOrderSummary)
			protected.GET("/reports/services/top", reportHandler.GetTopServices)
		}
	}

	r.Run() // defaults to :8080
}
