package dto

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

// === Outlet DTOs ===

// OutletRequest is the request body for creating/updating an outlet.
// UserID is NEVER accepted from frontend — injected from JWT context.
type OutletRequest struct {
	Name    string `json:"name" binding:"required,min=3,max=100"`
	Address string `json:"address" binding:"required,min=10"`
	Phone   string `json:"phone" binding:"required,min=10,max=15,numeric"`
}

// OutletResponse is the filtered response for outlet data.
type OutletResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// PaginationQuery is used to parse page and limit from URL query params.
type PaginationQuery struct {
	Page  int `form:"page,default=1" binding:"min=1"`
	Limit int `form:"limit,default=10" binding:"min=1,max=100"`
}

// PaginatedResponse wraps a list with pagination metadata.
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}
