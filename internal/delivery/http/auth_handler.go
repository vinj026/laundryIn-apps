package http

import (
	"errors"
	"net/http"

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
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", err.Error())
		return
	}

	resp, err := h.authUsecase.Register(req)
	if err != nil {
		if errors.Is(err, usecase.ErrDuplicatePhone) {
			utils.ErrorResponse(c, http.StatusConflict, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, resp)
}

// Login handles user login.
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Validasi gagal", err.Error())
		return
	}

	resp, err := h.authUsecase.Login(req)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			utils.ErrorResponse(c, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Terjadi kesalahan internal", nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, resp)
}
