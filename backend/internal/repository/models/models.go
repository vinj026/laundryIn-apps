package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Base provides a consistent string-based UUID primary key for all models.
// All IDs are plain Go strings — no uuid.UUID anywhere.
type Base struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

type User struct {
	Base
	Name     string   `gorm:"type:text;not null" json:"name"`
	Phone    string   `gorm:"type:text;uniqueIndex;not null" json:"phone"`
	Email    string   `gorm:"type:text" json:"email,omitempty"`
	Password string   `gorm:"type:text;not null" json:"-"`
	Role     string   `gorm:"type:text;not null" json:"role"`
	Outlets  []Outlet `gorm:"foreignKey:UserID" json:"outlets,omitempty"`
}

// Outlet now embeds Base for consistency with all other models.
type Outlet struct {
	Base
	UserID   string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Name     string    `gorm:"type:varchar(100);not null" json:"name"`
	Address  string    `gorm:"type:text;not null" json:"address"`
	Phone    string    `gorm:"type:varchar(20);not null" json:"phone"`
	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Services []Service `gorm:"foreignKey:OutletID;constraint:OnDelete:CASCADE" json:"services,omitempty"`
}

// Service — Price is now decimal.Decimal (not float64) for financial precision.
type Service struct {
	Base
	OutletID string          `gorm:"type:uuid;not null;index" json:"outlet_id"`
	Name     string          `gorm:"type:varchar(100);not null" json:"name"`
	Price    decimal.Decimal `gorm:"type:numeric(10,2);not null" json:"price"`
	Unit     string          `gorm:"type:varchar(20);not null" json:"unit"`
	Outlet   Outlet          `gorm:"foreignKey:OutletID;constraint:OnDelete:CASCADE" json:"outlet,omitempty"`
}

type Order struct {
	Base
	UserID     string          `gorm:"type:uuid;not null;index" json:"user_id"`
	OutletID   string          `gorm:"type:uuid;not null;index" json:"outlet_id"`
	TotalPrice decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"total_price"`
	Status     string          `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	OrderDate  time.Time       `gorm:"autoCreateTime" json:"order_date"`
	Items      []OrderItem     `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"items"`
	Logs       []OrderLog      `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"logs,omitempty"`
	User       User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Outlet     Outlet          `gorm:"foreignKey:OutletID" json:"outlet,omitempty"`
}

type OrderItem struct {
	ID           string          `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderID      string          `gorm:"type:uuid;not null;index" json:"order_id"`
	ServiceName  string          `gorm:"type:varchar(100);not null" json:"service_name"`
	ServicePrice decimal.Decimal `gorm:"type:numeric(10,2);not null" json:"service_price"`
	Qty          decimal.Decimal `gorm:"type:numeric(6,2);not null" json:"qty"`
	Unit         string          `gorm:"type:varchar(20);not null" json:"unit"`
	Subtotal     decimal.Decimal `gorm:"type:numeric(12,2);not null" json:"subtotal"`
}

type OrderLog struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	OrderID   string    `gorm:"type:uuid;not null;index" json:"order_id"`
	UpdatedBy string    `gorm:"type:uuid;not null" json:"updated_by"`
	OldStatus string    `gorm:"type:varchar(20)" json:"old_status"`
	NewStatus string    `gorm:"type:varchar(20);not null" json:"new_status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	User      User      `gorm:"foreignKey:UpdatedBy" json:"user,omitempty"`
}
