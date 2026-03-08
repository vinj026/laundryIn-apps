package repository

import (
	"context"
	"laundryin/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OutletRepository defines the interface for outlet database operations.
type OutletRepository interface {
	Create(ctx context.Context, outlet *models.Outlet) error
	FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.Outlet, error)
	FindByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (*models.Outlet, error)
	Update(ctx context.Context, outlet *models.Outlet) error
	Delete(ctx context.Context, outlet *models.Outlet) error
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

func (r *outletRepository) FindAllByUserID(ctx context.Context, userID uuid.UUID) ([]models.Outlet, error) {
	var outlets []models.Outlet
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&outlets).Error
	if err != nil {
		return nil, err
	}
	// Return empty array instead of nil
	if outlets == nil {
		outlets = []models.Outlet{}
	}
	return outlets, nil
}

func (r *outletRepository) FindByIDAndUserID(ctx context.Context, id, userID uuid.UUID) (*models.Outlet, error) {
	var outlet models.Outlet
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&outlet).Error
	if err != nil {
		return nil, err
	}
	return &outlet, nil
}

func (r *outletRepository) Update(ctx context.Context, outlet *models.Outlet) error {
	return r.db.WithContext(ctx).Save(outlet).Error
}

func (r *outletRepository) Delete(ctx context.Context, outlet *models.Outlet) error {
	return r.db.WithContext(ctx).Delete(outlet).Error
}
