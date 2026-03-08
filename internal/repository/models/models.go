package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
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

type Outlet struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Address   string         `gorm:"type:text;not null" json:"address"`
	Phone     string         `gorm:"type:varchar(20);not null" json:"phone"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Services  []Service      `gorm:"foreignKey:OutletID" json:"services,omitempty"`
}

type Service struct {
	Base
	OutletID uuid.UUID `gorm:"type:uuid;not null" json:"outlet_id"`
	Services string    `gorm:"type:text;not null" json:"services"`
	Price    float64   `gorm:"type:numeric(12,2);not null" json:"price"`
	Units    string    `gorm:"type:text;not null" json:"units"`
	Outlet   Outlet    `gorm:"foreignKey:OutletID" json:"outlet,omitempty"`
}

type Order struct {
	Base
	CustomerID uuid.UUID `gorm:"type:uuid;not null" json:"customer_id"`
	OutletID   uuid.UUID `gorm:"type:uuid;not null" json:"outlet_id"`
	ServicesID uuid.UUID `gorm:"type:uuid;not null" json:"services_id"`
	Quantity   float64   `gorm:"type:numeric(12,2);not null;default:1" json:"quantity"`
	Status     string    `gorm:"type:text;not null;default:'pending'" json:"status"`
	TotalPrice float64   `gorm:"type:numeric(12,2);not null" json:"total_price"`
	Customer   User      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Outlet     Outlet    `gorm:"foreignKey:OutletID" json:"outlet,omitempty"`
	Service    Service   `gorm:"foreignKey:ServicesID" json:"service,omitempty"`
}
