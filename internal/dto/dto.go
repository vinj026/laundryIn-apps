package dto

import "time"

// RegisterRequest is the request body for user registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Phone    string `json:"phone" binding:"required,e164,min=10,max=15"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
	Role     string `json:"role" binding:"required,oneof=owner customer"`
}

// LoginRequest is the request body for user login.
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse is the response body after successful auth.
type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

// OutletRequest is the request body for creating or updating an outlet.
type OutletRequest struct {
	Name    string `json:"name" binding:"required,min=3,max=100"`
	Address string `json:"address" binding:"required,min=10"`
	Phone   string `json:"phone" binding:"required,min=10,max=15,numeric"`
}

// OutletResponse is the response body for outlet operations.
type OutletResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
