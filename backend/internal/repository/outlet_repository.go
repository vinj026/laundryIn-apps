package repository

import (
	"context"

	"laundryin/internal/repository/models"

	"gorm.io/gorm"
)

// OutletRepository defines the interface for outlet database operations.
type OutletRepository interface {
	Create(ctx context.Context, outlet *models.Outlet) error
	FindAll(ctx context.Context, limit, offset int) ([]models.Outlet, int64, error)
	FindAllByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Outlet, int64, error)
	FindByID(ctx context.Context, outletID string) (*models.Outlet, error)
	FindByIDAndUserID(ctx context.Context, outletID, userID string) (*models.Outlet, error)
	Update(ctx context.Context, outlet *models.Outlet) error
	Delete(ctx context.Context, outletID, userID string) error
}

type outletRepository struct {
	db *gorm.DB
}

// NewOutletRepository creates a new OutletRepository instance.
func NewOutletRepository(db *gorm.DB) OutletRepository {
	return &outletRepository{db: db}
}

func (r *outletRepository) Create(ctx context.Context, outlet *models.Outlet) error {
	return r.db.WithContext(ctx).Create(outlet).Error
}

func (r *outletRepository) FindAll(ctx context.Context, limit, offset int) ([]models.Outlet, int64, error) {
	var outlets []models.Outlet
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Outlet{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Offset(offset).Find(&outlets).Error; err != nil {
		return nil, 0, err
	}

	return outlets, total, nil
}

func (r *outletRepository) FindAllByUserID(ctx context.Context, userID string, limit, offset int) ([]models.Outlet, int64, error) {
	var outlets []models.Outlet
	var total int64

	// Count total records for pagination metadata
	if err := r.db.WithContext(ctx).Model(&models.Outlet{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated records
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&outlets).Error; err != nil {
		return nil, 0, err
	}

	return outlets, total, nil
}

func (r *outletRepository) FindByID(ctx context.Context, outletID string) (*models.Outlet, error) {
	var outlet models.Outlet
	err := r.db.WithContext(ctx).Where("id = ?", outletID).First(&outlet).Error
	if err != nil {
		return nil, err
	}
	return &outlet, nil
}

func (r *outletRepository) FindByIDAndUserID(ctx context.Context, outletID, userID string) (*models.Outlet, error) {
	var outlet models.Outlet
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", outletID, userID).First(&outlet).Error
	if err != nil {
		return nil, err
	}
	return &outlet, nil
}

func (r *outletRepository) Update(ctx context.Context, outlet *models.Outlet) error {
	return r.db.WithContext(ctx).Save(outlet).Error
}

func (r *outletRepository) Delete(ctx context.Context, outletID, userID string) error {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", outletID, userID).Delete(&models.Outlet{}).Error
}
