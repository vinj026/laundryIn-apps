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
		// Fix BUG-006: Whitelist origins for WS
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true // Allow mobile/client directly
		}

		allowedOrigins := map[string]bool{
			"https://laundryin.vercel.app":     true,
			"https://www.laundryin.vercel.app": true,
			"http://localhost:3000":            true,
			"http://localhost:3001":            true,
		}

		return allowedOrigins[origin] || strings.HasSuffix(origin, ".vercel.app")
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
