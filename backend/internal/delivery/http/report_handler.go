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
)

type ReportHandler struct {
	reportUsecase usecase.ReportUsecase
}

func NewReportHandler(reportUsecase usecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{reportUsecase: reportUsecase}
}

// GetOmzet handles GET /api/v1/reports/omzet
func (h *ReportHandler) GetOmzet(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)

	var query dto.ReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format query parameter tidak valid", nil)
		return
	}

	resp, err := h.reportUsecase.GetOmzet(ctx, userID, query)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memproses data omzet", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Berhasil mengambil data omzet", resp)
}

// GetOrderSummary handles GET /api/v1/reports/orders/summary
func (h *ReportHandler) GetOrderSummary(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)

	var query dto.ReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format query parameter tidak valid", nil)
		return
	}

	resp, err := h.reportUsecase.GetOrderStatusSummary(ctx, userID, query)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memproses ringkasan order", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Berhasil mengambil ringkasan status pesanan", resp)
}

// GetTopServices handles GET /api/v1/reports/services/top
func (h *ReportHandler) GetTopServices(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID := c.MustGet("user_id").(string)

	var query dto.ReportQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format query parameter tidak valid", nil)
		return
	}

	resp, err := h.reportUsecase.GetTopServices(ctx, userID, query)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memproses data layanan terlaris", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Berhasil mengambil data layanan terlaris", resp)
}
