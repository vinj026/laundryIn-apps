package usecase

import (
	"context"
	"errors"
	"regexp"

	"laundryin/internal/dto"
	"laundryin/internal/repository"
	"laundryin/internal/repository/models"
	"laundryin/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthUsecase defines the interface for authentication business logic.
type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
}

// NewAuthUsecase creates a new AuthUsecase instance.
func NewAuthUsecase(userRepo repository.UserRepository) AuthUsecase {
	return &authUsecase{userRepo: userRepo}
}

// ErrDuplicatePhone is returned when the phone number is already registered.
var ErrDuplicatePhone = errors.New("nomor HP sudah terdaftar")

// ErrInvalidCredentials is a generic error for failed login (anti user enumeration).
var ErrInvalidCredentials = errors.New("nomor HP atau password salah")

// ErrWeakPassword is returned when the password doesn't meet complexity requirements.
var ErrWeakPassword = errors.New("password harus mengandung setidaknya satu huruf besar, satu huruf kecil, dan satu angka")

func (u *authUsecase) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Sanitize inputs — using shared utils.Sanitize
	req.Name = utils.Sanitize(req.Name)
	req.Phone = utils.Sanitize(req.Phone)
	req.Email = utils.Sanitize(req.Email)

	// Password complexity check
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(req.Password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(req.Password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(req.Password)

	if !hasUpper || !hasLower || !hasDigit {
		return nil, ErrWeakPassword
	}

	// Check for duplicate phone number
	existingUser, _ := u.userRepo.FindByPhone(ctx, req.Phone)
	if existingUser != nil {
		return nil, ErrDuplicatePhone
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	// Create user
	user := &models.User{
		Base: models.Base{
			ID: uuid.New().String(),
		},
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("gagal mendaftarkan user")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

func (u *authUsecase) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Sanitize phone
	req.Phone = utils.Sanitize(req.Phone)

	// Find user by phone
	user, err := u.userRepo.FindByPhone(ctx, req.Phone)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, errors.New("terjadi kesalahan internal")
	}

	// Check password
	if !utils.CheckPassword(user.Password, req.Password) {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  toUserResponse(user),
	}, nil
}

// toUserResponse maps a User model to a safe UserResponse DTO.
// Password is NEVER exposed.
func toUserResponse(u *models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Phone: u.Phone,
		Email: u.Email,
		Role:  u.Role,
	}
}
