package repository

import (
	"context"

	"laundryin/internal/repository/models"

	"gorm.io/gorm"
)

// ServiceRepository defines the interface for service database operations.
type ServiceRepository interface {
	Create(ctx context.Context, service *models.Service) error
	FindAllByOutletID(ctx context.Context, outletID, userID string) ([]models.Service, error)
	FindAllByOutletIDPublic(ctx context.Context, outletID string) ([]models.Service, error)
	FindByIDAndOwner(ctx context.Context, serviceID, userID string) (*models.Service, error)
	FindByIDAndOutletID(ctx context.Context, serviceID, outletID string) (*models.Service, error)
	Update(ctx context.Context, service *models.Service) error
	Delete(ctx context.Context, serviceID, userID string) error
}

type serviceRepository struct {
	db *gorm.DB
}

// NewServiceRepository creates a new ServiceRepository instance.
func NewServiceRepository(db *gorm.DB) ServiceRepository {
	return &serviceRepository{db: db}
}

func (r *serviceRepository) Create(ctx context.Context, service *models.Service) error {
	return r.db.WithContext(ctx).Create(service).Error
}

func (r *serviceRepository) FindAllByOutletID(ctx context.Context, outletID, userID string) ([]models.Service, error) {
	var services []models.Service
	// Anti-IDOR: Ensure the outlet actually belongs to the requesting user
	err := r.db.WithContext(ctx).
		Joins("JOIN outlets ON outlets.id = services.outlet_id").
		Where("services.outlet_id = ? AND outlets.user_id = ?", outletID, userID).
		Find(&services).Error
	return services, err
}

func (r *serviceRepository) FindAllByOutletIDPublic(ctx context.Context, outletID string) ([]models.Service, error) {
	var services []models.Service
	err := r.db.WithContext(ctx).
		Where("outlet_id = ?", outletID).
		Find(&services).Error
	return services, err
}

func (r *serviceRepository) FindByIDAndOwner(ctx context.Context, serviceID, userID string) (*models.Service, error) {
	var service models.Service
	// Anti-IDOR: Strict JOIN to ensure the service's outlet is owned by the user
	err := r.db.WithContext(ctx).
		Joins("JOIN outlets ON outlets.id = services.outlet_id").
		Where("services.id = ? AND outlets.user_id = ?", serviceID, userID).
		First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) FindByIDAndOutletID(ctx context.Context, serviceID, outletID string) (*models.Service, error) {
	var service models.Service
	err := r.db.WithContext(ctx).
		Where("id = ? AND outlet_id = ?", serviceID, outletID).
		First(&service).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceRepository) Update(ctx context.Context, service *models.Service) error {
	// GORM's Save updates all fields, including zero values, for an existing record.
	return r.db.WithContext(ctx).Save(service).Error
}

func (r *serviceRepository) Delete(ctx context.Context, serviceID, userID string) error {
	// First, fetch the service applying the Anti-IDOR check
	service, err := r.FindByIDAndOwner(ctx, serviceID, userID)
	if err != nil {
		return err
	}
	// Proceed to delete (soft delete due to gorm.DeletedAt)
	return r.db.WithContext(ctx).Delete(service).Error
}
