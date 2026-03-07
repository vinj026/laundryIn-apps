package http

import (
	"net/http"
	"strings"

	"laundryin/pkg/utils"

	"github.com/gin-gonic/gin"
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
