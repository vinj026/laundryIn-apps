package repository

import (
	"context"
	"laundryin/internal/repository/models"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(ctx context.Context, notif *models.Notification) error
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Notification, error)
	GetUnreadCount(ctx context.Context, userID string) (int64, error)
	GetTotalCount(ctx context.Context, userID string) (int64, error)
	MarkAsRead(ctx context.Context, notifID string, userID string) error
	MarkAllAsRead(ctx context.Context, userID string) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notif *models.Notification) error {
	return r.db.WithContext(ctx).Create(notif).Error
}

func (r *notificationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Notification, error) {
	var notifs []models.Notification
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifs).Error
	return notifs, err
}

func (r *notificationRepository) GetUnreadCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) GetTotalCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, notifID string, userID string) error {
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notifID, userID).
		Update("is_read", true).Error
}

func (r *notificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).
		Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}
