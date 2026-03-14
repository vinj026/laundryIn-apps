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

// AuthHandler handles authentication HTTP requests.
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{authUsecase: authUsecase}
}

// Register handles user registration.
// POST /api/v1/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	// Set 5-second timeout for registration
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Stop leak: hide gin struct field name details
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi format data gagal", nil)
		return
	}

	resp, err := h.authUsecase.Register(ctx, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", err.Error())
			return
		}
		if errors.Is(err, usecase.ErrWeakPassword) {
			utils.ErrorResponse(c, http.StatusBadRequest, err.Error(), err.Error())
			return
		}
		if errors.Is(err, usecase.ErrDuplicatePhone) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), err.Error())
			return
		}
		// Log external error for cloud debugging (Railway)
		fmt.Printf("🔴 REGISTRATION ERROR: %v\n", err)
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Registrasi berhasil", resp)
}

// Login handles user login.
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	// Set 5-second timeout for login
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi format data gagal", nil)
		return
	}

	resp, err := h.authUsecase.Login(ctx, req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			utils.ErrorResponse(c, http.StatusRequestTimeout, "Proses terlalu lama, silakan coba lagi", err.Error())
			return
		}
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Login berhasil", resp)
}
