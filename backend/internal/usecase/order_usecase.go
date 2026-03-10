package usecase

import (
	"context"
	"errors"

	"laundryin/internal/dto"
	"laundryin/internal/repository"
	"laundryin/internal/repository/models"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// ErrOrderNotFound is returned when an order is not found or user lacks access.
var ErrOrderNotFound = errors.New("pesanan tidak ditemukan atau akses ditolak")
var ErrStateInvalid = errors.New("transisi status pesanan tidak valid")

type OrderUsecase interface {
	Create(ctx context.Context, userID string, req dto.OrderRequest) (*dto.OrderResponse, error)
	GetAllByUserID(ctx context.Context, userID string, page, limit int) (*dto.PaginatedResponse, error)
	GetAllByOutletID(ctx context.Context, outletID, userID string, page, limit int) (*dto.PaginatedResponse, error)
	UpdateStatus(ctx context.Context, orderID, userID string, req dto.OrderStatusRequest) (*dto.OrderResponse, error)
}

type orderUsecase struct {
	orderRepo   repository.OrderRepository
	serviceRepo repository.ServiceRepository
	outletRepo  repository.OutletRepository
}

func NewOrderUsecase(or repository.OrderRepository, sr repository.ServiceRepository, ou repository.OutletRepository) OrderUsecase {
	return &orderUsecase{orderRepo: or, serviceRepo: sr, outletRepo: ou}
}

func (u *orderUsecase) Create(ctx context.Context, userID string, req dto.OrderRequest) (*dto.OrderResponse, error) {
	// 1. Anti-IDOR: Verify the target outlet exists and is owned by the current user.
	outlet, err := u.outletRepo.FindByID(ctx, req.OutletID)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("outlet tidak ditemukan")
		}
		return nil, errors.New("gagal memvalidasi outlet")
	}

	// Ownership check
	if outlet.UserID != userID {
		return nil, errors.New("akses ditolak: anda bukan pemilik outlet ini")
	}

	// 2. Deep Anti-IDOR & Zero-Trust Pricing.
	var items []models.OrderItem
	var grandTotal decimal.Decimal

	for _, itemReq := range req.Items {
		s, err := u.serviceRepo.FindByIDAndOutletID(ctx, itemReq.ServiceID, req.OutletID)
		if err != nil {
			if ctx.Err() != nil {
				return nil, ctx.Err()
			}
			return nil, errors.New("satu atau lebih layanan tidak ditemukan atau bukan milik outlet ini")
		}

		// Parse qty from string to decimal (DTO fields are now strings)
		qty, err := decimal.NewFromString(itemReq.Qty)
		if err != nil || qty.LessThanOrEqual(decimal.Zero) {
			return nil, errors.New("qty harus berupa angka positif yang valid")
		}
		qty = qty.Round(2)

		// Service.Price is already decimal.Decimal from the model
		servicePrice := s.Price.Round(2)
		subtotal := qty.Mul(servicePrice).Round(2)

		grandTotal = grandTotal.Add(subtotal)

		items = append(items, models.OrderItem{
			ID:           uuid.New().String(),
			ServiceName:  s.Name,
			ServicePrice: servicePrice,
			Qty:          qty,
			Unit:         s.Unit,
			Subtotal:     subtotal,
		})
	}

	// 3. Build the Order header
	order := &models.Order{
		Base: models.Base{
			ID: uuid.New().String(),
		},
		UserID:     userID,
		OutletID:   req.OutletID,
		Status:     "pending",
		TotalPrice: grandTotal,
		Logs: []models.OrderLog{
			{
				ID:        uuid.New().String(),
				UpdatedBy: userID,
				OldStatus: "",
				NewStatus: "pending",
			},
		},
	}

	// 4. ACID Transaction Insert
	if err := u.orderRepo.CreateOrderWithItems(ctx, order, items); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal memproses transaksi pesanan")
	}

	order.Items = items
	return toOrderResponse(order), nil
}

func (u *orderUsecase) GetAllByUserID(ctx context.Context, userID string, page, limit int) (*dto.PaginatedResponse, error) {
	offset := (page - 1) * limit
	orders, total, err := u.orderRepo.FindAllByUserID(ctx, userID, limit, offset)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal mengambil riwayat pesanan")
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.PaginatedResponse{
		Data:       toOrderResponseList(orders),
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (u *orderUsecase) GetAllByOutletID(ctx context.Context, outletID, userID string, page, limit int) (*dto.PaginatedResponse, error) {
	offset := (page - 1) * limit
	orders, total, err := u.orderRepo.FindAllByOutletIDAndOwner(ctx, outletID, userID, limit, offset)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal mengambil data pesanan outlet")
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.PaginatedResponse{
		Data:       toOrderResponseList(orders),
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (u *orderUsecase) UpdateStatus(ctx context.Context, orderID, userID string, req dto.OrderStatusRequest) (*dto.OrderResponse, error) {
	// Get Current State + Ownership (Anti-IDOR)
	order, err := u.orderRepo.FindByIDAndOwner(ctx, orderID, userID)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, errors.New("gagal memverifikasi pesanan")
	}

	// FSM Rule Check
	if !isValidTransition(order.Status, req.Status) {
		return nil, ErrStateInvalid
	}

	// Create Audit Log
	logEntry := &models.OrderLog{
		ID:        uuid.New().String(),
		OrderID:   orderID,
		UpdatedBy: userID,
		OldStatus: order.Status,
		NewStatus: req.Status,
	}

	// Execute Update Transaction (Status change + Log append)
	if err := u.orderRepo.UpdateOrderStatus(ctx, orderID, req.Status, logEntry); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal memperbarui status")
	}

	order.Status = req.Status
	return toOrderResponse(order), nil
}

// -- Helpers --

func isValidTransition(current, next string) bool {
	switch current {
	case "pending":
		return next == "process" || next == "cancelled"
	case "process":
		return next == "completed"
	case "completed":
		return next == "picked_up"
	default:
		return false
	}
}

func toOrderResponseList(orders []models.Order) []dto.OrderResponse {
	resp := make([]dto.OrderResponse, 0, len(orders))
	for i := range orders {
		resp = append(resp, *toOrderResponse(&orders[i]))
	}
	return resp
}

// toOrderResponse — all financial fields serialized as decimal strings.
func toOrderResponse(o *models.Order) *dto.OrderResponse {
	itemsResp := make([]dto.OrderItemResponse, 0, len(o.Items))
	for _, it := range o.Items {
		itemsResp = append(itemsResp, dto.OrderItemResponse{
			ID:           it.ID,
			ServiceName:  it.ServiceName,
			ServicePrice: it.ServicePrice.StringFixed(2),
			Qty:          it.Qty.StringFixed(2),
			Unit:         it.Unit,
			Subtotal:     it.Subtotal.StringFixed(2),
		})
	}

	logsResp := make([]dto.OrderLogResponse, 0, len(o.Logs))
	for _, lg := range o.Logs {
		logsResp = append(logsResp, dto.OrderLogResponse{
			OldStatus: lg.OldStatus,
			NewStatus: lg.NewStatus,
			UpdatedBy: lg.UpdatedBy,
			CreatedAt: lg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &dto.OrderResponse{
		ID:         o.ID,
		UserID:     o.UserID,
		OutletID:   o.OutletID,
		TotalPrice: o.TotalPrice.StringFixed(2),
		Status:     o.Status,
		OrderDate:  o.OrderDate.Format("2006-01-02T15:04:05Z07:00"),
		Items:      itemsResp,
		Logs:       logsResp,
	}
}
