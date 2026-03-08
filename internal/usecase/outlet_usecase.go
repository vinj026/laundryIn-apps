package usecase

import (
	"context"
	"errors"

	"laundryin/internal/dto"
	"laundryin/internal/repository"
	"laundryin/internal/repository/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// OutletUsecase defines the interface for outlet business logic.
type OutletUsecase interface {
	CreateOutlet(ctx context.Context, userID uuid.UUID, req dto.OutletRequest) (*dto.OutletResponse, error)
	GetAllOutlets(ctx context.Context, userID uuid.UUID) ([]dto.OutletResponse, error)
	GetOutletByID(ctx context.Context, id, userID uuid.UUID) (*dto.OutletResponse, error)
	UpdateOutlet(ctx context.Context, id, userID uuid.UUID, req dto.OutletRequest) (*dto.OutletResponse, error)
	DeleteOutlet(ctx context.Context, id, userID uuid.UUID) error
}

type outletUsecase struct {
	outletRepo repository.OutletRepository
}

// NewOutletUsecase creates a new OutletUsecase instance.
func NewOutletUsecase(outletRepo repository.OutletRepository) OutletUsecase {
	return &outletUsecase{outletRepo: outletRepo}
}

func (u *outletUsecase) CreateOutlet(ctx context.Context, userID uuid.UUID, req dto.OutletRequest) (*dto.OutletResponse, error) {
	outlet := &models.Outlet{
		UserID:  userID,
		Name:    req.Name,
		Address: req.Address,
		Phone:   req.Phone,
	}

	if err := u.outletRepo.Create(ctx, outlet); err != nil {
		return nil, errors.New("gagal membuat outlet")
	}

	return toOutletResponse(outlet), nil
}

func (u *outletUsecase) GetAllOutlets(ctx context.Context, userID uuid.UUID) ([]dto.OutletResponse, error) {
	outlets, err := u.outletRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("gagal mengambil data outlet")
	}

	responses := make([]dto.OutletResponse, 0)
	for _, o := range outlets {
		responses = append(responses, *toOutletResponse(&o))
	}

	return responses, nil
}

func (u *outletUsecase) GetOutletByID(ctx context.Context, id, userID uuid.UUID) (*dto.OutletResponse, error) {
	outlet, err := u.outletRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("outlet tidak ditemukan atau bukan milik Anda")
		}
		return nil, errors.New("terjadi kesalahan internal")
	}

	return toOutletResponse(outlet), nil
}

func (u *outletUsecase) UpdateOutlet(ctx context.Context, id, userID uuid.UUID, req dto.OutletRequest) (*dto.OutletResponse, error) {
	// 1. Verify ownership
	outlet, err := u.outletRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("outlet tidak ditemukan atau bukan milik Anda")
		}
		return nil, errors.New("terjadi kesalahan internal")
	}

	// 2. Update fields
	outlet.Name = req.Name
	outlet.Address = req.Address
	outlet.Phone = req.Phone

	if err := u.outletRepo.Update(ctx, outlet); err != nil {
		return nil, errors.New("gagal memperbarui outlet")
	}

	return toOutletResponse(outlet), nil
}

func (u *outletUsecase) DeleteOutlet(ctx context.Context, id, userID uuid.UUID) error {
	// 1. Verify ownership
	outlet, err := u.outletRepo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("outlet tidak ditemukan atau bukan milik Anda")
		}
		return errors.New("terjadi kesalahan internal")
	}

	// 2. Soft delete
	if err := u.outletRepo.Delete(ctx, outlet); err != nil {
		return err
	}

	return nil
}

func toOutletResponse(o *models.Outlet) *dto.OutletResponse {
	return &dto.OutletResponse{
		ID:        o.ID.String(),
		Name:      o.Name,
		Address:   o.Address,
		Phone:     o.Phone,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
	}
}
