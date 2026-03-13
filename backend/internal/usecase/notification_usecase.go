package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"laundryin/internal/dto"
	"laundryin/internal/repository"
	"laundryin/internal/repository/models"
	"laundryin/internal/websocket"
	"strings"
	"time"

	"github.com/google/uuid"
)

type NotificationUsecase interface {
	NotifyOrderCreated(ctx context.Context, order *models.Order) error
	NotifyStatusChanged(ctx context.Context, order *models.Order, oldStatus, newStatus string) error

	GetNotifications(ctx context.Context, userID string, page, limit int) ([]dto.NotificationResponse, int64, error)
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	MarkAsRead(ctx context.Context, notifID string, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}

type notificationUsecase struct {
	repo       repository.NotificationRepository
	userRepo   repository.UserRepository
	outletRepo repository.OutletRepository
	hub        *websocket.Hub
}

func NewNotificationUsecase(repo repository.NotificationRepository, userRepo repository.UserRepository, outletRepo repository.OutletRepository, hub *websocket.Hub) NotificationUsecase {
	return &notificationUsecase{
		repo:       repo,
		userRepo:   userRepo,
		outletRepo: outletRepo,
		hub:        hub,
	}
}

func (u *notificationUsecase) NotifyOrderCreated(ctx context.Context, order *models.Order) error {
	// Notif to owner
	outlet, _ := u.outletRepo.FindByID(ctx, order.OutletID)
	if outlet == nil {
		return nil
	}

	customer, _ := u.userRepo.FindByID(ctx, order.UserID)
	customerName := "Customer"
	if customer != nil {
		customerName = customer.Name
	}

	// List services for body
	var services []string
	for _, item := range order.Items {
		services = append(services, item.ServiceName)
	}
	servicesStr := "laundry"
	if len(services) > 0 {
		servicesStr = strings.Join(services, ", ")
	}

	title := "Pesanan Baru Masuk"
	body := fmt.Sprintf("%s memesan %s di %s", customerName, servicesStr, outlet.Name)
	
	notifData, _ := json.Marshal(map[string]interface{}{
		"order_id":      order.ID,
		"outlet_id":     order.OutletID,
		"customer_name": customerName,
		"total_price":   order.TotalPrice.StringFixed(2),
	})

	notif := &models.Notification{
		ID:     uuid.New().String(),
		UserID: outlet.UserID,
		Type:   "new_order",
		Title:  title,
		Body:   body,
		Data:   string(notifData),
	}

	if err := u.repo.Create(ctx, notif); err != nil {
		fmt.Printf("Failed to create notification: %v\n", err)
		return err
	}

	// WS
	u.hub.SendToUser(outlet.UserID, websocket.Message{
		ID:        notif.ID,
		Type:      "new_order",
		Title:     title,
		Body:      body,
		Data:      map[string]interface{}{"order_id": order.ID, "outlet_id": order.OutletID, "customer_name": customerName, "total_price": order.TotalPrice.StringFixed(2)},
		Timestamp: time.Now(),
	})

	return nil
}

func (u *notificationUsecase) NotifyStatusChanged(ctx context.Context, order *models.Order, oldStatus, newStatus string) error {
	outlet, _ := u.outletRepo.FindByID(ctx, order.OutletID)
	outletName := "Outlet"
	if outlet != nil {
		outletName = outlet.Name
	}

	orderIDShort := order.ID
	if len(orderIDShort) > 8 {
		orderIDShort = orderIDShort[:8]
	}

	var title, body, nType string
	var data map[string]interface{}

	switch newStatus {
	case "process":
		nType = "order_status"
		title = "Pesananmu Sedang Diproses"
		body = fmt.Sprintf("Outlet %s mulai memproses pesanan #%s", outletName, orderIDShort)
		data = map[string]interface{}{"order_id": order.ID, "new_status": newStatus}
		
		if order.FinalTotalPrice != nil {
			// Also price updated notif
			u.notifyPriceUpdated(ctx, order)
		}
	case "completed":
		nType = "order_status"
		title = "Pesananmu Siap Diambil! 🎉"
		body = fmt.Sprintf("Cucian kamu sudah selesai, silakan ambil di %s", outletName)
		data = map[string]interface{}{"order_id": order.ID, "new_status": newStatus}
	case "cancelled":
		nType = "order_cancelled"
		title = "Pesananmu Dibatalkan"
		body = fmt.Sprintf("Pesanan #%s di %s telah dibatalkan", orderIDShort, outletName)
		data = map[string]interface{}{"order_id": order.ID}
	default:
		return nil
	}

	notifData, _ := json.Marshal(data)

	notif := &models.Notification{
		ID:     uuid.New().String(),
		UserID: order.UserID,
		Type:   nType,
		Title:  title,
		Body:   body,
		Data:   string(notifData),
	}

	if err := u.repo.Create(ctx, notif); err != nil {
		fmt.Printf("Failed to create notification: %v\n", err)
		return err
	}

	u.hub.SendToUser(order.UserID, websocket.Message{
		ID:        notif.ID,
		Type:      nType,
		Title:     title,
		Body:      body,
		Data:      data,
		Timestamp: time.Now(),
	})

	return nil
}

func (u *notificationUsecase) notifyPriceUpdated(ctx context.Context, order *models.Order) {
	orderIDShort := order.ID
	if len(orderIDShort) > 8 {
		orderIDShort = orderIDShort[:8]
	}

	title := "Harga Final Pesananmu Sudah Diketahui"
	body := fmt.Sprintf("Total pembayaran pesanan #%s adalah Rp %s", orderIDShort, order.FinalTotalPrice.StringFixed(2))
	
	data := map[string]interface{}{
		"order_id":        order.ID,
		"estimated_price": order.TotalPrice.StringFixed(2),
		"final_price":     order.FinalTotalPrice.StringFixed(2),
	}
	notifData, _ := json.Marshal(data)

	notif := &models.Notification{
		ID:     uuid.New().String(),
		UserID: order.UserID,
		Type:   "price_updated",
		Title:  title,
		Body:   body,
		Data:   string(notifData),
	}

	if err := u.repo.Create(ctx, notif); err != nil {
		fmt.Printf("Failed to create notification: %v\n", err)
		return
	}

	u.hub.SendToUser(order.UserID, websocket.Message{
		ID:        notif.ID,
		Type:      "price_updated",
		Title:     title,
		Body:      body,
		Data:      data,
		Timestamp: time.Now(),
	})
}

func (u *notificationUsecase) GetNotifications(ctx context.Context, userID string, page, limit int) ([]dto.NotificationResponse, int64, error) {
	offset := (page - 1) * limit
	notifs, err := u.repo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, _ := u.repo.GetTotalCount(ctx, userID) 

	res := make([]dto.NotificationResponse, len(notifs))
	for i, n := range notifs {
		var data interface{}
		json.Unmarshal([]byte(n.Data), &data)
		res[i] = dto.NotificationResponse{
			ID:        n.ID,
			Type:      n.Type,
			Title:     n.Title,
			Body:      n.Body,
			Data:      data,
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt.Format(time.RFC3339),
		}
	}

	return res, total, nil
}

func (u *notificationUsecase) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	return u.repo.GetUnreadCount(ctx, userID)
}

func (u *notificationUsecase) MarkAsRead(ctx context.Context, notifID string, userID string) error {
	return u.repo.MarkAsRead(ctx, notifID, userID)
}

func (u *notificationUsecase) MarkAllAsRead(ctx context.Context, userID string) error {
	return u.repo.MarkAllAsRead(ctx, userID)
}
