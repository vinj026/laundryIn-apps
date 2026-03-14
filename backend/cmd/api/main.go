package main

import (
	"fmt"
	"log"
	"os"

	"laundryin/internal/database"
	handler "laundryin/internal/delivery/http"
	"laundryin/internal/repository"
	"laundryin/internal/repository/models"
	"laundryin/internal/usecase"
	"laundryin/internal/websocket"
	"laundryin/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables (optional for production environments like Railway)
	_ = godotenv.Load()

	// Register custom validators (must be before any request handling)
	utils.RegisterCustomValidators()

	// Connect to database
	db := database.ConnectDB()

	// Auto-migrate models
	err := db.AutoMigrate(&models.User{}, &models.Outlet{}, &models.Service{}, &models.Order{}, &models.OrderItem{}, &models.OrderLog{}, &models.Notification{})
	if err != nil {
		log.Fatalf("❌ CRITICAL: Database migration failed: %v", err)
	}
	fmt.Println("🚀 Database Migration Successful!")

	// Initialize WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// === Dependency Injection ===
	// Repository layer
	userRepo := repository.NewUserRepository(db)
	outletRepo := repository.NewOutletRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	reportRepo := repository.NewReportRepository(db)
	notifRepo := repository.NewNotificationRepository(db)

	// Usecase layer
	authUsecase := usecase.NewAuthUsecase(userRepo)
	outletUsecase := usecase.NewOutletUsecase(outletRepo)
	serviceUsecase := usecase.NewServiceUsecase(serviceRepo, outletRepo)
	notifUsecase := usecase.NewNotificationUsecase(notifRepo, userRepo, outletRepo, hub)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, serviceRepo, outletRepo, notifUsecase)
	reportUsecase := usecase.NewReportUsecase(reportRepo)

	// Handler layer
	authHandler := handler.NewAuthHandler(authUsecase)
	outletHandler := handler.NewOutletHandler(outletUsecase)
	serviceHandler := handler.NewServiceHandler(serviceUsecase)
	orderHandler := handler.NewOrderHandler(orderUsecase)
	reportHandler := handler.NewReportHandler(reportUsecase)
	notifHandler := handler.NewNotificationHandler(notifUsecase, hub)

	// === Router Setup ===
	// Production hardening: Respect GIN_MODE=release and set TrustedProxies
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Railway/Cloudflare/Vercel proxies are usually fine to trust for header mapping in this context
	_ = r.SetTrustedProxies(nil)

	// Global security middleware
	r.Use(handler.CORSMiddleware())
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

		public := v1.Group("/public")
		{
			public.GET("/outlets", outletHandler.GetAllOutletsPublic)
			public.GET("/outlets/:id", outletHandler.GetOutletByIDPublic)
			public.GET("/outlets/:id/services", serviceHandler.GetAllByOutletIDPublic)
		}

		// Protected routes — Customer
		customer := v1.Group("")
		customer.Use(handler.AuthMiddleware(), handler.RoleMiddleware("customer"))
		{
			// Orders (Customer)
			customer.POST("/orders", orderHandler.CreateOrder)
			customer.GET("/orders", orderHandler.GetAllByUserID)
		}

		// Protected routes — Owner
		owner := v1.Group("")
		owner.Use(handler.AuthMiddleware(), handler.RoleMiddleware("owner"))
		{
			// Outlets
			owner.POST("/outlets", outletHandler.CreateOutlet)
			owner.GET("/outlets", outletHandler.GetAllOutlets)
			owner.GET("/outlets/:id", outletHandler.GetOutletByID)
			owner.PUT("/outlets/:id", outletHandler.UpdateOutlet)
			owner.DELETE("/outlets/:id", outletHandler.DeleteOutlet)

			// Services
			owner.POST("/services", serviceHandler.CreateService)
			owner.GET("/outlets/:id/services", serviceHandler.GetAllByOutletID)
			owner.PUT("/services/:id", serviceHandler.UpdateService)
			owner.DELETE("/services/:id", serviceHandler.DeleteService)

			// Orders (Owner)
			owner.GET("/outlets/:id/orders", orderHandler.GetAllByOutletID)
			owner.PATCH("/orders/:id/status", orderHandler.UpdateStatus)

			// Reports & Analytics
			owner.GET("/reports/omzet", reportHandler.GetOmzet)
			owner.GET("/reports/orders/summary", reportHandler.GetOrderSummary)
			owner.GET("/reports/services/top", reportHandler.GetTopServices)
		}

		// Notifications (Both roles)
		authorized := v1.Group("/notifications")
		authorized.Use(handler.AuthMiddleware())
		{
			authorized.GET("", notifHandler.GetNotifications)
			authorized.GET("/unread-count", notifHandler.GetUnreadCount)
			authorized.PATCH("/:id/read", notifHandler.MarkAsRead)
			authorized.PATCH("/read-all", notifHandler.MarkAllAsRead)
		}

		// WebSocket
		v1.GET("/ws/connect", handler.AuthMiddleware(), notifHandler.Connect)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
