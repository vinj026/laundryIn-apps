package http

import (
	"context"
	"errors"
	"net/http"
	"time"

	"laundryin/internal/dto"
	"laundryin/internal/usecase"
	"laundryin/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(orderUsecase usecase.OrderUsecase) *OrderHandler {
	return &OrderHandler{orderUsecase: orderUsecase}
}

// CreateOrder handles placing a new order.
// POST /api/v1/orders
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)

	var req dto.OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi format data gagal", nil)
		return
	}

	resp, err := h.orderUsecase.Create(ctx, userID, req)
	if err != nil {
		// Stop leak: Do not use err.Error() directly unless known
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		// If specific usecase failures vs generic failure
		if err.Error() == "outlet tidak ditemukan" ||
			err.Error() == "akses ditolak: anda bukan pemilik outlet ini" ||
			err.Error() == "satu atau lebih layanan tidak ditemukan atau bukan milik outlet ini" ||
			err.Error() == "qty harus berupa angka positif yang valid" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), nil)
			return
		}

		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memproses pesanan", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Pesanan berhasil dibuat", resp)
}

// GetAllByUserID returns orders for the currently authenticated User
// GET /api/v1/orders
func (h *OrderHandler) GetAllByUserID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)

	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format pagination tidak valid (maksimal limit 100)", nil)
		return
	}

	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}

	resp, err := h.orderUsecase.GetAllByUserID(ctx, userID, pagination.Page, pagination.Limit)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil daftar pesanan", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar pesanan berhasil diambil", resp)
}

// GetAllByOutletID handles getting all orders incoming to a specific outlet.
// GET /api/v1/outlets/:id/orders
func (h *OrderHandler) GetAllByOutletID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)
	outletID := c.Param("id")

	if _, err := uuid.Parse(outletID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format Outlet ID tidak valid", nil)
		return
	}

	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format pagination tidak valid (maksimal limit 100)", nil)
		return
	}

	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}

	resp, err := h.orderUsecase.GetAllByOutletID(ctx, outletID, userID, pagination.Page, pagination.Limit)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data pesanan", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data pesanan outlet berhasil diambil", resp)
}

// UpdateStatus handles advancing the FSM order state.
// PATCH /api/v1/orders/:id/status
func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)
	orderID := c.Param("id")

	if _, err := uuid.Parse(orderID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format Order ID tidak valid", nil)
		return
	}

	var req dto.OrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi format data gagal", nil)
		return
	}

	resp, err := h.orderUsecase.UpdateStatus(ctx, orderID, userID, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		if errors.Is(err, usecase.ErrOrderNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Pesanan tidak ditemukan atau akses ditolak", nil)
			return
		}
		if errors.Is(err, usecase.ErrStateInvalid) {
			utils.ErrorResponse(c, http.StatusBadRequest, "Transisi status pesanan tidak valid", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui status pesanan", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Status pesanan berhasil diperbarui", resp)
}
