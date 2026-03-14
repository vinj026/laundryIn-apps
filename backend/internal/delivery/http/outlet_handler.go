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

// OutletHandler handles outlet HTTP requests.
type OutletHandler struct {
	outletUsecase usecase.OutletUsecase
}

// NewOutletHandler creates a new OutletHandler instance.
func NewOutletHandler(outletUsecase usecase.OutletUsecase) *OutletHandler {
	return &OutletHandler{outletUsecase: outletUsecase}
}

// CreateOutlet handles creating a new outlet.
// POST /api/v1/outlets
func (h *OutletHandler) CreateOutlet(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)

	var req dto.OutletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi format data gagal", nil)
		return
	}

	resp, err := h.outletUsecase.Create(ctx, userID, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Outlet berhasil dibuat", resp)
}

// GetAllOutlets handles getting all outlets for the authenticated owner.
// GET /api/v1/outlets?page=1&limit=10
func (h *OutletHandler) GetAllOutlets(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)

	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format pagination tidak valid (maksimal limit 100)", nil)
		return
	}

	// Ensure defaults if missing
	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}

	resp, err := h.outletUsecase.GetAll(ctx, userID, pagination.Page, pagination.Limit)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar outlet berhasil diambil", resp)
}

// GetOutletByID handles getting a specific outlet by ID.
// GET /api/v1/outlets/:id
func (h *OutletHandler) GetOutletByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)
	outletID := c.Param("id")

	// Validate UUID format
	if _, err := uuid.Parse(outletID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid", nil)
		return
	}

	resp, err := h.outletUsecase.GetByID(ctx, outletID, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		if errors.Is(err, usecase.ErrOutletNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data outlet berhasil diambil", resp)
}

// UpdateOutlet handles updating an outlet.
// PUT /api/v1/outlets/:id
func (h *OutletHandler) UpdateOutlet(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)
	outletID := c.Param("id")

	// Validate UUID format
	if _, err := uuid.Parse(outletID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid", nil)
		return
	}

	var req dto.OutletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi format data gagal", nil)
		return
	}

	resp, err := h.outletUsecase.Update(ctx, outletID, userID, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		if errors.Is(err, usecase.ErrOutletNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Outlet berhasil diupdate", resp)
}

// DeleteOutlet handles soft-deleting an outlet.
// DELETE /api/v1/outlets/:id
func (h *OutletHandler) DeleteOutlet(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)
	outletID := c.Param("id")

	// Validate UUID format
	if _, err := uuid.Parse(outletID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid", nil)
		return
	}

	err := h.outletUsecase.Delete(ctx, outletID, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		if errors.Is(err, usecase.ErrOutletNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Outlet berhasil dihapus", nil)
}

// GetAllOutletsPublic handles getting all outlets without auth.
// GET /api/v1/public/outlets?page=1&limit=10
func (h *OutletHandler) GetAllOutletsPublic(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format pagination tidak valid (maksimal limit 100)", nil)
		return
	}

	// Ensure defaults if missing
	if pagination.Page == 0 {
		pagination.Page = 1
	}
	if pagination.Limit == 0 {
		pagination.Limit = 10
	}

	resp, err := h.outletUsecase.GetAllPublic(ctx, pagination.Page, pagination.Limit)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar outlet public berhasil diambil", resp)
}

// GetOutletByIDPublic handles getting a specific outlet by ID without auth.
// GET /api/v1/public/outlets/:id
func (h *OutletHandler) GetOutletByIDPublic(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	outletID := c.Param("id")

	// Validate UUID format
	if _, err := uuid.Parse(outletID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid", nil)
		return
	}

	resp, err := h.outletUsecase.GetByIDPublic(ctx, outletID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		if errors.Is(err, usecase.ErrOutletNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Data outlet public berhasil diambil", resp)
}
