package usecase

import (
	"context"
	"errors"

	"laundryin/internal/dto"
	"laundryin/internal/repository"
	"laundryin/internal/repository/models"
	"laundryin/pkg/utils"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ErrServiceNotFound is returned when a service is not found or user lacks access.
var ErrServiceNotFound = errors.New("layanan tidak ditemukan atau akses ditolak")

// ServiceUsecase defines the interface for service business logic.
type ServiceUsecase interface {
	Create(ctx context.Context, userID string, req dto.ServiceRequest) (*dto.ServiceResponse, error)
	GetAllByOutletID(ctx context.Context, outletID, userID string) ([]dto.ServiceResponse, error)
	GetAllByOutletIDPublic(ctx context.Context, outletID string) ([]dto.ServiceResponse, error)
	Update(ctx context.Context, serviceID, userID string, req dto.ServiceRequest) (*dto.ServiceResponse, error)
	Delete(ctx context.Context, serviceID, userID string) error
}

type serviceUsecase struct {
	serviceRepo repository.ServiceRepository
	outletRepo  repository.OutletRepository
}

// NewServiceUsecase creates a new ServiceUsecase instance.
func NewServiceUsecase(serviceRepo repository.ServiceRepository, outletRepo repository.OutletRepository) ServiceUsecase {
	return &serviceUsecase{serviceRepo: serviceRepo, outletRepo: outletRepo}
}

func (u *serviceUsecase) Create(ctx context.Context, userID string, req dto.ServiceRequest) (*dto.ServiceResponse, error) {
	// Sanitize inputs
	req.Name = utils.Sanitize(req.Name)
	req.Unit = utils.Sanitize(req.Unit)

	// Parse and validate price from string to decimal
	price, err := decimal.NewFromString(req.Price)
	if err != nil || price.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("harga harus berupa angka positif yang valid")
	}

	// Step 1: Verify the outlet belongs to the current user (Anti-IDOR Create)
	_, err = u.outletRepo.FindByIDAndUserID(ctx, req.OutletID, userID)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOutletNotFound
		}
		return nil, errors.New("gagal memvalidasi outlet")
	}

	// All IDs are plain strings — no uuid.Parse needed
	service := &models.Service{
		Base: models.Base{
			ID: uuid.New().String(),
		},
		OutletID: req.OutletID,
		Name:     req.Name,
		Price:    price,
		Unit:     req.Unit,
	}

	if err := u.serviceRepo.Create(ctx, service); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal membuat layanan")
	}

	return toServiceResponse(service), nil
}

func (u *serviceUsecase) GetAllByOutletID(ctx context.Context, outletID, userID string) ([]dto.ServiceResponse, error) {
	services, err := u.serviceRepo.FindAllByOutletID(ctx, outletID, userID)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal mengambil daftar layanan")
	}

	// Always return empty array instead of nil
	responses := make([]dto.ServiceResponse, 0, len(services))
	for i := range services {
		responses = append(responses, *toServiceResponse(&services[i]))
	}

	return responses, nil
}

func (u *serviceUsecase) GetAllByOutletIDPublic(ctx context.Context, outletID string) ([]dto.ServiceResponse, error) {
	services, err := u.serviceRepo.FindAllByOutletIDPublic(ctx, outletID)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal mengambil daftar layanan")
	}

	// Always return empty array instead of nil
	responses := make([]dto.ServiceResponse, 0, len(services))
	for i := range services {
		responses = append(responses, *toServiceResponse(&services[i]))
	}

	return responses, nil
}

func (u *serviceUsecase) Update(ctx context.Context, serviceID, userID string, req dto.ServiceRequest) (*dto.ServiceResponse, error) {
	// Sanitize inputs
	req.Name = utils.Sanitize(req.Name)
	req.Unit = utils.Sanitize(req.Unit)

	// Parse and validate price from string to decimal
	price, err := decimal.NewFromString(req.Price)
	if err != nil || price.LessThanOrEqual(decimal.Zero) {
		return nil, errors.New("harga harus berupa angka positif yang valid")
	}

	// Verify the service belongs to an outlet owned by the user
	service, err := u.serviceRepo.FindByIDAndOwner(ctx, serviceID, userID)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrServiceNotFound
		}
		return nil, errors.New("gagal mengambil data layanan")
	}

	// If OutletID changed, verify the new outlet is also owned by the user (Anti-IDOR transfer)
	if req.OutletID != service.OutletID {
		_, err := u.outletRepo.FindByIDAndUserID(ctx, req.OutletID, userID)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrOutletNotFound
			}
			return nil, errors.New("gagal memvalidasi outlet tujuan")
		}
		service.OutletID = req.OutletID
	}

	service.Name = req.Name
	service.Price = price
	service.Unit = req.Unit

	if err := u.serviceRepo.Update(ctx, service); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal mengupdate layanan")
	}

	return toServiceResponse(service), nil
}

func (u *serviceUsecase) Delete(ctx context.Context, serviceID, userID string) error {
	// Deletion itself uses Delete with Anti-IDOR checks inside the repository
	if err := u.serviceRepo.Delete(ctx, serviceID, userID); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrServiceNotFound
		}
		return errors.New("gagal menghapus layanan")
	}
	return nil
}

// toServiceResponse maps a Service model to a ServiceResponse DTO.
// All IDs are plain strings, Price is decimal.Decimal.StringFixed(2).
func toServiceResponse(service *models.Service) *dto.ServiceResponse {
	return &dto.ServiceResponse{
		ID:        service.ID,
		OutletID:  service.OutletID,
		Name:      service.Name,
		Price:     service.Price.StringFixed(2),
		Unit:      service.Unit,
		CreatedAt: service.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: service.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
