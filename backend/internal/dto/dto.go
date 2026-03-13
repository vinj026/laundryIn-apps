package dto

// RegisterRequest is the request body for user registration.
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Phone    string `json:"phone" binding:"required,e164_strict"`
	Email    string `json:"email" binding:"omitempty,email"`
	Password string `json:"password" binding:"required,min=8,max=64"`
	Role     string `json:"role" binding:"required,oneof=owner customer"`
}

// LoginRequest is the request body for user login.
type LoginRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserResponse is a safe, filtered representation of a user.
// Password is NEVER included — this is the only user shape exposed via API.
type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email,omitempty"`
	Role  string `json:"role"`
}

// AuthResponse is the response body after successful auth.
type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// === Outlet DTOs ===

// OutletRequest is the request body for creating/updating an outlet.
// UserID is NEVER accepted from frontend — injected from JWT context.
type OutletRequest struct {
	Name    string `json:"name" binding:"required,min=3,max=100"`
	Address string `json:"address" binding:"required,min=10,max=500"`
	Phone   string `json:"phone" binding:"required,e164_strict"`
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

// === Service DTOs ===

// ServiceRequest — Price is string to preserve decimal precision.
// Backend will parse to shopspring/decimal for validation and calculation.
type ServiceRequest struct {
	OutletID string `json:"outlet_id" binding:"required,uuid"`
	Name     string `json:"name" binding:"required,min=3,max=100"`
	Price    string `json:"price" binding:"required"`
	Unit     string `json:"unit" binding:"required,oneof=KG PCS METER"`
}

// ServiceResponse — Price is string (decimal representation) to avoid floating-point issues.
type ServiceResponse struct {
	ID        string `json:"id"`
	OutletID  string `json:"outlet_id"`
	Name      string `json:"name"`
	Price     string `json:"price"`
	Unit      string `json:"unit"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// === Order DTOs ===

// OrderItemRequest — Qty is string for decimal precision.
type OrderItemRequest struct {
	ServiceID string `json:"service_id" binding:"required,uuid"`
	Qty       string `json:"qty" binding:"required"`
}

// OrderRequest is the request body for creating an order.
type OrderRequest struct {
	OutletID string             `json:"outlet_id" binding:"required,uuid"`
	Items    []OrderItemRequest `json:"items" binding:"required,min=1,dive"`
}

// OrderStatusRequest for advancing the FSM.
type OrderStatusRequest struct {
	Status string          `json:"status" binding:"required,oneof=pending process completed picked_up cancelled"`
	Items  []ActualQtyItem `json:"items,omitempty"`
}

type ActualQtyItem struct {
	ID        string `json:"id" binding:"required,uuid"`
	ActualQty string `json:"actual_qty" binding:"required"`
}

// OrderItemResponse — all financial fields are strings (decimal precision).
type OrderItemResponse struct {
	ID           string  `json:"id"`
	ServiceName  string  `json:"service_name"`
	ServicePrice string  `json:"service_price"`
	Qty          string  `json:"qty"`
	ActualQty    *string `json:"actual_qty,omitempty"`
	Unit         string  `json:"unit"`
	Subtotal     string  `json:"subtotal"`
	FinalPrice   *string `json:"final_price,omitempty"`
}

type OrderLogResponse struct {
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
	UpdatedBy string `json:"updated_by"`
	CreatedAt string `json:"created_at"`
}

// OrderResponse — TotalPrice is string (decimal precision).
type OrderResponse struct {
	ID              string              `json:"id"`
	UserID          string              `json:"user_id"`
	CustomerName    string              `json:"customer_name,omitempty"`
	OutletID        string              `json:"outlet_id"`
	TotalPrice      string              `json:"total_price"`
	FinalTotalPrice *string             `json:"final_total_price,omitempty"`
	Status          string              `json:"status"`
	OrderDate       string              `json:"order_date"`
	Items           []OrderItemResponse `json:"items,omitempty"`
	Logs            []OrderLogResponse  `json:"logs,omitempty"`
}

// === Notification DTOs ===

type NotificationResponse struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Title     string      `json:"title"`
	Body      string      `json:"body"`
	Data      interface{} `json:"data"`
	IsRead    bool        `json:"is_read"`
	CreatedAt string      `json:"created_at"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}
