package dto

// RegisterRequest is the request body for user registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Phone    string `json:"phone" binding:"required,e164,min=10,max=15"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=6"`
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
