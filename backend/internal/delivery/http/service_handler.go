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

// ServiceHandler handles service HTTP requests.
type ServiceHandler struct {
	serviceUsecase usecase.ServiceUsecase
}

// NewServiceHandler creates a new ServiceHandler instance.
func NewServiceHandler(serviceUsecase usecase.ServiceUsecase) *ServiceHandler {
	return &ServiceHandler{serviceUsecase: serviceUsecase}
}

// CreateService handles creating a new service.
// POST /api/v1/services
func (h *ServiceHandler) CreateService(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)

	var req dto.ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi format data gagal", nil)
		return
	}

	resp, err := h.serviceUsecase.Create(ctx, userID, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", err.Error())
			return
		}
		if errors.Is(err, usecase.ErrOutletNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Outlet tidak ditemukan atau Anda bukan pemiliknya", err.Error())
			return
		}
		if err.Error() == "harga harus berupa angka positif yang valid" || err.Error() == "gagal memvalidasi outlet" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Layanan berhasil dibuat", resp)
}

// GetAllByOutletID handles getting all services for a specific outlet.
// GET /api/v1/outlets/:id/services
func (h *ServiceHandler) GetAllByOutletID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)
	outletID := c.Param("id")

	if _, err := uuid.Parse(outletID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format Outlet ID tidak valid", nil)
		return
	}

	resp, err := h.serviceUsecase.GetAllByOutletID(ctx, outletID, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar layanan berhasil diambil", resp)
}

// UpdateService handles updating a service.
// PUT /api/v1/services/:id
func (h *ServiceHandler) UpdateService(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)
	serviceID := c.Param("id")

	if _, err := uuid.Parse(serviceID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format Service ID tidak valid", nil)
		return
	}

	var req dto.ServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi format data gagal", nil)
		return
	}

	resp, err := h.serviceUsecase.Update(ctx, serviceID, userID, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", err.Error())
			return
		}
		if errors.Is(err, usecase.ErrServiceNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Layanan tidak ditemukan atau akses ditolak", err.Error())
			return
		}
		if errors.Is(err, usecase.ErrOutletNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Outlet tujuan tidak ditemukan atau akses ditolak", err.Error())
			return
		}
		if err.Error() == "harga harus berupa angka positif yang valid" || err.Error() == "gagal memvalidasi outlet" {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Layanan berhasil diupdate", resp)
}

// DeleteService handles soft-deleting a service.
// DELETE /api/v1/services/:id
func (h *ServiceHandler) DeleteService(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)
	serviceID := c.Param("id")

	if _, err := uuid.Parse(serviceID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format Service ID tidak valid", nil)
		return
	}

	err := h.serviceUsecase.Delete(ctx, serviceID, userID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", err.Error())
			return
		}
		if errors.Is(err, usecase.ErrServiceNotFound) {
			utils.ErrorResponse(c, http.StatusNotFound, "Layanan tidak ditemukan atau akses ditolak", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Layanan berhasil dihapus", nil)
}

// GetAllByOutletIDPublic handles getting all services for a specific outlet without auth.
// GET /api/v1/public/outlets/:id/services
func (h *ServiceHandler) GetAllByOutletIDPublic(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	outletID := c.Param("id")

	if _, err := uuid.Parse(outletID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format Outlet ID tidak valid", nil)
		return
	}

	// This specific usecase bypasses the userID ownership check
	resp, err := h.serviceUsecase.GetAllByOutletIDPublic(ctx, outletID)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Daftar layanan public berhasil diambil", resp)
}
