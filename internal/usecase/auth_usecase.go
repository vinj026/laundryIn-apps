package usecase

import (
	"errors"

	"laundryin/internal/dto"
	"laundryin/internal/repository"
	"laundryin/internal/repository/models"
	"laundryin/pkg/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuthUsecase defines the interface for authentication business logic.
type AuthUsecase interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
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

func (u *authUsecase) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check for duplicate phone number
	existingUser, _ := u.userRepo.FindByPhone(req.Phone)
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
			ID: uuid.New(),
		},
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     req.Role,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.New("gagal mendaftarkan user")
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  user, // Password excluded via json:"-"
	}, nil
}

func (u *authUsecase) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	// Find user by phone
	user, err := u.userRepo.FindByPhone(req.Phone)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Generic error: don't reveal that the user doesn't exist
			return nil, ErrInvalidCredentials
		}
		return nil, errors.New("terjadi kesalahan internal")
	}

	// Check password
	if !utils.CheckPassword(user.Password, req.Password) {
		// Generic error: don't reveal that the password is wrong
		return nil, ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return nil, errors.New("gagal membuat token")
	}

	return &dto.AuthResponse{
		Token: token,
		User:  user, // Password excluded via json:"-"
	}, nil
}
