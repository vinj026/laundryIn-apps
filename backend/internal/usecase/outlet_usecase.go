package usecase

import (
	"context"
	"errors"
	"math"

	"laundryin/internal/dto"
	"laundryin/internal/repository"
	"laundryin/internal/repository/models"
	"laundryin/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ErrOutletNotFound is returned when the outlet is not found or not owned by the user.
var ErrOutletNotFound = errors.New("outlet tidak ditemukan")

// OutletUsecase defines the interface for outlet business logic.
type OutletUsecase interface {
	Create(ctx context.Context, userID string, req dto.OutletRequest) (*dto.OutletResponse, error)
	GetAll(ctx context.Context, userID string, page, limit int) (*dto.PaginatedResponse, error)
	GetAllPublic(ctx context.Context, page, limit int) (*dto.PaginatedResponse, error)
	GetByID(ctx context.Context, outletID, userID string) (*dto.OutletResponse, error)
	GetByIDPublic(ctx context.Context, outletID string) (*dto.OutletResponse, error)
	Update(ctx context.Context, outletID, userID string, req dto.OutletRequest) (*dto.OutletResponse, error)
	Delete(ctx context.Context, outletID, userID string) error
}

type outletUsecase struct {
	outletRepo repository.OutletRepository
}

// NewOutletUsecase creates a new OutletUsecase instance.
func NewOutletUsecase(outletRepo repository.OutletRepository) OutletUsecase {
	return &outletUsecase{outletRepo: outletRepo}
}

func (u *outletUsecase) Create(ctx context.Context, userID string, req dto.OutletRequest) (*dto.OutletResponse, error) {
	// Sanitize inputs
	req.Name = utils.Sanitize(req.Name)
	req.Address = utils.Sanitize(req.Address)
	req.Phone = utils.Sanitize(req.Phone)

	outlet := &models.Outlet{
		Base: models.Base{
			ID: uuid.New().String(),
		},
		UserID:  userID,
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
	}

	if err := u.outletRepo.Create(ctx, outlet); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal membuat outlet")
	}

	return toOutletResponse(outlet), nil
}

func (u *outletUsecase) GetAllPublic(ctx context.Context, page, limit int) (*dto.PaginatedResponse, error) {
	offset := (page - 1) * limit

	outlets, total, err := u.outletRepo.FindAll(ctx, limit, offset)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		// Return original error so it can be logged by the handler
		return nil, err
	}

	// Always return empty array, never nil
	responses := make([]dto.OutletResponse, 0, len(outlets))
	for i := range outlets {
		responses = append(responses, *toOutletResponse(&outlets[i]))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.PaginatedResponse{
		Data:       responses,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (u *outletUsecase) GetAll(ctx context.Context, userID string, page, limit int) (*dto.PaginatedResponse, error) {
	offset := (page - 1) * limit

	outlets, total, err := u.outletRepo.FindAllByUserID(ctx, userID, limit, offset)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal mengambil data outlet")
	}

	// Always return empty array, never nil
	responses := make([]dto.OutletResponse, 0, len(outlets))
	for i := range outlets {
		responses = append(responses, *toOutletResponse(&outlets[i]))
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &dto.PaginatedResponse{
		Data:       responses,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (u *outletUsecase) GetByID(ctx context.Context, outletID, userID string) (*dto.OutletResponse, error) {
	outlet, err := u.outletRepo.FindByIDAndUserID(ctx, outletID, userID)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOutletNotFound
		}
		return nil, errors.New("gagal mengambil data outlet")
	}

	return toOutletResponse(outlet), nil
}

func (u *outletUsecase) GetByIDPublic(ctx context.Context, outletID string) (*dto.OutletResponse, error) {
	outlet, err := u.outletRepo.FindByID(ctx, outletID)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOutletNotFound
		}
		return nil, errors.New("gagal mengambil data outlet")
	}

	return toOutletResponse(outlet), nil
}

func (u *outletUsecase) Update(ctx context.Context, outletID, userID string, req dto.OutletRequest) (*dto.OutletResponse, error) {
	// Verify ownership first
	outlet, err := u.outletRepo.FindByIDAndUserID(ctx, outletID, userID)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOutletNotFound
		}
		return nil, errors.New("gagal mengambil data outlet")
	}

	// Sanitize and update fields
	outlet.Name = utils.Sanitize(req.Name)
	outlet.Address = utils.Sanitize(req.Address)
	outlet.Phone = utils.Sanitize(req.Phone)

	if err := u.outletRepo.Update(ctx, outlet); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal mengupdate outlet")
	}

	return toOutletResponse(outlet), nil
}

func (u *outletUsecase) Delete(ctx context.Context, outletID, userID string) error {
	// Verify ownership first
	_, err := u.outletRepo.FindByIDAndUserID(ctx, outletID, userID)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrOutletNotFound
		}
		return errors.New("gagal mengambil data outlet")
	}

	if err := u.outletRepo.Delete(ctx, outletID, userID); err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return errors.New("gagal menghapus outlet")
	}

	return nil
}

// toOutletResponse converts a model to a response DTO.
func toOutletResponse(o *models.Outlet) *dto.OutletResponse {
	return &dto.OutletResponse{
		ID:        o.ID,
		Name:      o.Name,
		Address:   o.Address,
		Phone:     o.Phone,
		CreatedAt: o.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: o.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
