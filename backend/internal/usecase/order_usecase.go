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
	orderRepo    repository.OrderRepository
	serviceRepo  repository.ServiceRepository
	outletRepo   repository.OutletRepository
	notifUsecase NotificationUsecase
}

func NewOrderUsecase(or repository.OrderRepository, sr repository.ServiceRepository, ou repository.OutletRepository, nu NotificationUsecase) OrderUsecase {
	return &orderUsecase{orderRepo: or, serviceRepo: sr, outletRepo: ou, notifUsecase: nu}
}

func (u *orderUsecase) Create(ctx context.Context, userID string, req dto.OrderRequest) (*dto.OrderResponse, error) {
	// 1. Verify the target outlet exists.
	_, err := u.outletRepo.FindByID(ctx, req.OutletID)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("outlet tidak ditemukan")
		}
		return nil, errors.New("gagal memvalidasi outlet")
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

	// Fire notification in background
	go u.notifUsecase.NotifyOrderCreated(context.Background(), order)

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

	// Handle items actual qty for 'process' state transitions
	if req.Status == "process" {
		reqItemsMap := make(map[string]decimal.Decimal)
		for _, ri := range req.Items {
			qty, err := decimal.NewFromString(ri.ActualQty)
			if err != nil || qty.LessThanOrEqual(decimal.Zero) {
				return nil, errors.New("berat aktual harus berupa angka positif yang valid")
			}
			reqItemsMap[ri.ID] = qty.Round(2)
		}

		var finalTotalPrice decimal.Decimal

		for i, item := range order.Items {
			if item.Unit == "KG" {
				if actualQty, exists := reqItemsMap[item.ID]; exists {
					order.Items[i].ActualQty = &actualQty
					finalPrice := actualQty.Mul(item.ServicePrice).Round(2)
					order.Items[i].FinalPrice = &finalPrice
					finalTotalPrice = finalTotalPrice.Add(finalPrice)
				} else {
					return nil, errors.New("berat aktual wajib diisi untuk layanan per KG")
				}
			} else {
				// PCS does not need actual_qty, finalPrice = Subtotal
				finalPrice := item.Subtotal
				order.Items[i].FinalPrice = &finalPrice
				finalTotalPrice = finalTotalPrice.Add(finalPrice)
			}
		}

		order.FinalTotalPrice = &finalTotalPrice
	}

	// Create Audit Log
	logEntry := &models.OrderLog{
		ID:        uuid.New().String(),
		OrderID:   orderID,
		UpdatedBy: userID,
		OldStatus: order.Status,
		NewStatus: req.Status,
	}

	order.Status = req.Status

	// Execute Update Transaction (Status change + Final Price update + Log append)
	if err := u.orderRepo.UpdateOrderStatus(ctx, order, logEntry); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal memperbarui status")
	}

	// Fire notification in background
	go u.notifUsecase.NotifyStatusChanged(context.Background(), order, logEntry.OldStatus, req.Status)

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
		item := dto.OrderItemResponse{
			ID:           it.ID,
			ServiceName:  it.ServiceName,
			ServicePrice: it.ServicePrice.StringFixed(2),
			Qty:          it.Qty.StringFixed(2),
			Unit:         it.Unit,
			Subtotal:     it.Subtotal.StringFixed(2),
		}

		if it.ActualQty != nil {
			aq := it.ActualQty.StringFixed(2)
			item.ActualQty = &aq
		}
		if it.FinalPrice != nil {
			fp := it.FinalPrice.StringFixed(2)
			item.FinalPrice = &fp
		}
		
		itemsResp = append(itemsResp, item)
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

	resp := &dto.OrderResponse{
		ID:           o.ID,
		UserID:       o.UserID,
		CustomerName: o.User.Name,
		OutletID:     o.OutletID,
		TotalPrice:   o.TotalPrice.StringFixed(2),
		Status:       o.Status,
		OrderDate:    o.OrderDate.Format("2006-01-02T15:04:05Z07:00"),
		Items:        itemsResp,
		Logs:         logsResp,
	}

	if o.FinalTotalPrice != nil {
		ftp := o.FinalTotalPrice.StringFixed(2)
		resp.FinalTotalPrice = &ftp
	}

	return resp
}
