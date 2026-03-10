package repository

import (
	"context"

	"laundryin/internal/repository/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	// CreateOrderWithItems saves the order and order items atomically.
	CreateOrderWithItems(ctx context.Context, order *models.Order, items []models.OrderItem) error

	// FindAllByUserID fetches order history for a customer
	FindAllByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Order, int64, error)

	// FindAllByOutletIDAndOwner fetches orders for a specific outlet, ONLY IF owned by the user
	FindAllByOutletIDAndOwner(ctx context.Context, outletID, userID string, limit, offset int) ([]models.Order, int64, error)

	// FindByIDAndOwner fetches a specific order, verifying via JOIN that the caller owns the outlet
	FindByIDAndOwner(ctx context.Context, orderID, userID string) (*models.Order, error)

	// UpdateOrderStatus updates the order status and simultaneously saves an OrderLog atomically
	UpdateOrderStatus(ctx context.Context, orderID, newStatus string, logEntry *models.OrderLog) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrderWithItems(ctx context.Context, order *models.Order, items []models.OrderItem) error {
	// Execute everything in an ACID transaction
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Create the Order header
		if err := tx.Create(order).Error; err != nil {
			return err // Rollback
		}

		// 2. Insert items sequentially; on any failure, the transaction rolls back all
		for i := range items {
			items[i].OrderID = order.ID // Link item to the freshly created Order (string ID)
			if err := tx.Create(&items[i]).Error; err != nil {
				return err // Rollback ALL previous items and the header
			}
		}

		return nil // Commit
	})
}

func (r *orderRepository) FindAllByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	modelDB := r.db.WithContext(ctx).Model(&models.Order{}).Where("user_id = ?", userID)

	if err := modelDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := modelDB.
		Preload("Items").
		Order("order_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, total, err
}

func (r *orderRepository) FindAllByOutletIDAndOwner(ctx context.Context, outletID, userID string, limit, offset int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	// Anti-IDOR: verify the target outlet belongs to the user requesting the list.
	modelDB := r.db.WithContext(ctx).Model(&models.Order{}).
		Joins("JOIN outlets ON outlets.id = orders.outlet_id").
		Where("orders.outlet_id = ? AND outlets.user_id = ?", outletID, userID)

	if err := modelDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := modelDB.
		Preload("Items").
		Order("orders.order_date DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, total, err
}

func (r *orderRepository) FindByIDAndOwner(ctx context.Context, orderID, userID string) (*models.Order, error) {
	var order models.Order

	// Anti-IDOR: JOIN ensures the order's outlet is owned by the caller
	err := r.db.WithContext(ctx).
		Joins("JOIN outlets ON outlets.id = orders.outlet_id").
		Preload("Items").
		Where("orders.id = ? AND outlets.user_id = ?", orderID, userID).
		First(&order).Error

	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) UpdateOrderStatus(ctx context.Context, orderID, newStatus string, logEntry *models.OrderLog) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Update order status; IDs are now plain strings, no uuid.Parse needed
		if err := tx.Model(&models.Order{}).Where("id = ?", orderID).Update("status", newStatus).Error; err != nil {
			return err
		}

		// Store Who Changed It via OrderLog
		if err := tx.Create(logEntry).Error; err != nil {
			return err
		}

		return nil
	})
}
