package http

import (
	"log"
	"net/http"
	"strings"
	"time"

	"laundryin/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// AuthMiddleware validates JWT token from the Authorization header.
// On success, it sets "user_id" and "role" in the Gin context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token tidak ditemukan", nil)
			c.Abort()
			return
		}

		// Expect format: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Format token tidak valid", nil)
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Token tidak valid atau sudah kadaluarsa", nil)
			c.Abort()
			return
		}

		// Set user info in context for downstream handlers
		c.Set("user_id", claims.UserID.String())
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RoleMiddleware checks if the authenticated user has one of the allowed roles.
// Must be used AFTER AuthMiddleware.
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "Akses ditolak", nil)
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			utils.ErrorResponse(c, http.StatusForbidden, "Akses ditolak", nil)
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusForbidden, "Anda tidak memiliki izin untuk mengakses resource ini", nil)
		c.Abort()
	}
}

// PayloadLimit limits the size of the request body to prevent Denial of Service attacks.
func PayloadLimit(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}

// RateLimiter creates a middleware to limit requests.
// Configured here for 5 requests per 5 minutes for brute-force protection.
func RateLimiter() gin.HandlerFunc {
	// 5 requests per 5 minutes
	rate := limiter.Rate{
		Period: 5 * time.Minute,
		Limit:  5,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return mgin.NewMiddleware(instance, mgin.WithLimitReachedHandler(func(c *gin.Context) {
		utils.ErrorResponse(c, http.StatusTooManyRequests, "Terlalu banyak percobaan login, silakan coba lagi dalam 5 menit", nil)
	}))
}

func init() {
	log.Println("Middleware initialized")
}
