package http

import (
	"context"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"laundryin/pkg/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// AuthMiddleware validates JWT token from the Authorization header.
// On success, it sets "user_id" and "role" in the Gin context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// Fallback to query param for WebSockets
			authHeader = c.Query("token")
			if authHeader != "" {
				authHeader = "Bearer " + authHeader
			}
		}

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
		c.Set("user_id", claims.UserID)
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

// visitor struct to track last seen time for proper memory cleanup.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter implements IP-based rate limiting using a token bucket algorithm.
// Each IP address gets its own limiter: 1 request/second with burst up to 60.
// This prevents brute-force attacks and API abuse while allowing normal usage.
// Includes a background cleaner to prevent memory leaks (OOM).
func RateLimiter() gin.HandlerFunc {
	var (
		mu       sync.Mutex
		visitors = make(map[string]*visitor)
	)

	// Fix BUG-018: Add context support for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	// Usually this would be tied to server shutdown, 
	// for simplicity in middleware we'll just ensure it has a ticker
	
	// Background worker to clean up inactive IPs
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		
		for {
			select {
			case <-ticker.C:
				mu.Lock()
				for ip, v := range visitors {
					if time.Since(v.lastSeen) > 3*time.Minute {
						delete(visitors, ip)
					}
				}
				mu.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	// NOTE: In a real Gin app, you'd call cancel() when the server stops.
	// For this middleware instance, it will run for the process life.
	_ = cancel 

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			// 1.66 token per second, burst of 100 = ~100 requests/minute max
			limiter := rate.NewLimiter(1.66, 100)
			visitors[ip] = &visitor{limiter: limiter, lastSeen: time.Now()}
			v = visitors[ip]
		}
		v.lastSeen = time.Now()
		mu.Unlock()

		if !v.limiter.Allow() {
			fmt.Printf("⚠️  RATE LIMIT EXCEEDED for IP: %s\n", ip)
			utils.ErrorResponse(c, http.StatusTooManyRequests, "Terlalu banyak request, silakan coba lagi nanti", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecurityHeaders adds security-related HTTP headers to every response.
// Prevents XSS, clickjacking, MIME sniffing, and cache leaking.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Cache-Control", "no-store")
		c.Header("Content-Security-Policy", "default-src 'none'")
		c.Header("Referrer-Policy", "no-referrer")
		c.Next()
	}
}
// CORSMiddleware handles Cross-Origin Resource Sharing.
// Dynamically allowed origins from ALLOWED_ORIGINS environment variable.
func CORSMiddleware() gin.HandlerFunc {
	// Pre-parse allowed origins for performance
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	allowedOrigins := make(map[string]bool)

	if allowedOriginsStr != "" {
		origins := strings.Split(allowedOriginsStr, ",")
		for _, o := range origins {
			allowedOrigins[strings.TrimSpace(o)] = true
		}
	} else {
		// Sensible defaults for development
		allowedOrigins["http://localhost:3000"] = true
		allowedOrigins["http://localhost:3001"] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Check if origin is allowed by ENV list or is a Vercel/Railway preview
		isVercelSubdomain := strings.HasSuffix(origin, ".vercel.app")
		isRailwaySubdomain := strings.HasSuffix(origin, ".up.railway.app")

		isAllowed := allowedOrigins[origin] || isVercelSubdomain || isRailwaySubdomain

		// For production mode, be strict
		if origin != "" && !isAllowed && os.Getenv("GIN_MODE") == "release" {
			fmt.Printf("❌ CORS REJECTED: Origin '%s' is not in ALLOWED_ORIGINS\n", origin)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		// Set CORS headers
		effectiveOrigin := origin
		if effectiveOrigin == "" {
			effectiveOrigin = "*"
		} else if !isAllowed {
			effectiveOrigin = "*" // Fallback for dev
		}

		c.Header("Access-Control-Allow-Origin", effectiveOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
