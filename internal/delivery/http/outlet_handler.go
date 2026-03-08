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

// Create handles creating a new outlet.
// POST /api/v1/outlets
func (h *OutletHandler) Create(c *gin.Context) {
	// 5-second timeout for creation
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID, err := getUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Sesi tidak valid", nil)
		return
	}

	var req dto.OutletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", err.Error())
		return
	}

	resp, err := h.outletUsecase.CreateOutlet(ctx, userID, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, resp)
}

// GetAll handles fetching all outlets belonging to the authenticated owner.
// GET /api/v1/outlets
func (h *OutletHandler) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userID, err := getUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Sesi tidak valid", nil)
		return
	}

	resp, err := h.outletUsecase.GetAllOutlets(ctx, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, resp)
}

// GetByID handles fetching a single outlet by its ID.
// GET /api/v1/outlets/:id
func (h *OutletHandler) GetByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID outlet tidak valid", nil)
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Sesi tidak valid", nil)
		return
	}

	resp, err := h.outletUsecase.GetOutletByID(ctx, id, userID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, resp)
}

// Update handles updating an existing outlet.
// PUT /api/v1/outlets/:id
func (h *OutletHandler) Update(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID outlet tidak valid", nil)
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Sesi tidak valid", nil)
		return
	}

	var req dto.OutletRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", err.Error())
		return
	}

	resp, err := h.outletUsecase.UpdateOutlet(ctx, id, userID, req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, resp)
}

// Delete handles soft deleting an outlet.
// DELETE /api/v1/outlets/:id
func (h *OutletHandler) Delete(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID outlet tidak valid", nil)
		return
	}

	userID, err := getUserIDFromContext(c)
	if err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Sesi tidak valid", nil)
		return
	}

	if err := h.outletUsecase.DeleteOutlet(ctx, id, userID); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Outlet berhasil dihapus"})
}

// getUserIDFromContext extracts and parses user_id from Gin context.
func getUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	idStr, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user_id not found")
	}

	id, err := uuid.Parse(idStr.(string))
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
