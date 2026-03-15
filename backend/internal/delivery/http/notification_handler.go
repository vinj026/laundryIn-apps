package http

import (
	"laundryin/internal/usecase"
	"laundryin/internal/websocket"
	"net/http"
	"strconv"
	"strings"

	"laundryin/pkg/utils"

	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

type NotificationHandler struct {
	notifUsecase usecase.NotificationUsecase
	hub          *websocket.Hub
}

func NewNotificationHandler(notifUsecase usecase.NotificationUsecase, hub *websocket.Hub) *NotificationHandler {
	return &NotificationHandler{
		notifUsecase: notifUsecase,
		hub:          hub,
	}
}

var upgrader = gorilla.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true // Allow mobile/client directly
		}

		// Use dynamic validation similar to CORSMiddleware
		allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
		isVercelSubdomain := strings.HasSuffix(origin, ".vercel.app")
		isRailwaySubdomain := strings.HasSuffix(origin, ".up.railway.app")

		if isVercelSubdomain || isRailwaySubdomain {
			return true
		}

		if allowedOriginsStr != "" {
			origins := strings.Split(allowedOriginsStr, ",")
			for _, o := range origins {
				if strings.TrimSpace(o) == origin {
					return true
				}
			}
		}

		// Fallback for localhost in non-release mode
		if os.Getenv("GIN_MODE") != "release" {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
		}

		fmt.Printf("❌ WS ORIGIN REJECTED: Origin '%s' is not allowed\n", origin)
		return false
	},
}

func (h *NotificationHandler) Connect(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &websocket.Client{
		UserID: userID.(string),
		Role:   role.(string),
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.hub.Register() <- client

	// Start pumping messages
	go client.WritePump()
	go client.ReadPump(h.hub)
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	notifs, total, err := h.notifUsecase.GetNotifications(c.Request.Context(), userID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil notifikasi"})
		return
	}

	unread, _ := h.notifUsecase.GetUnreadCount(c.Request.Context(), userID)

	totalPages := 0
	if limit > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Notifikasi berhasil diambil",
		"data": gin.H{
			"data":        notifs,
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
		"unread_count": unread,
	})
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	count, err := h.notifUsecase.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil jumlah notifikasi"})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Jumlah unread count berhasil diambil", gin.H{"count": count})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	notifID := c.Param("id")

	err := h.notifUsecase.MarkAsRead(c.Request.Context(), notifID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal update notifikasi"})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notifikasi ditandai sebagai sudah baca", nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	err := h.notifUsecase.MarkAllAsRead(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal update notifikasi"})
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Semua notifikasi ditandai sebagai sudah baca", nil)
}
