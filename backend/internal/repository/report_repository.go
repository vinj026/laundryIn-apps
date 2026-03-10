package repository

import (
	"context"

	"laundryin/internal/dto"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetTotalOmzet(ctx context.Context, userID string, req dto.ReportQuery) (decimal.Decimal, error)
	GetOrderStatusSummary(ctx context.Context, userID string, req dto.ReportQuery) (map[string]int64, error)
	GetTopServices(ctx context.Context, userID string, req dto.ReportQuery) ([]TopServiceRow, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

// Internal Structs for GORM Scan
type TopServiceRow struct {
	ServiceName  string
	OutletName   string
	TotalQty     decimal.Decimal
	TotalRevenue decimal.Decimal
}

type OrderStatusRow struct {
	Status string
	Count  int64
}

// Build base query with dynamic filters (StartDate, EndDate, OutletID)
// Table prefix is applied for disambiguation when JOINs are involved.
func (r *reportRepository) buildBaseQuery(ctx context.Context, query *gorm.DB, userID string, tablePrefix string, req dto.ReportQuery) *gorm.DB {
	prefix := ""
	if tablePrefix != "" {
		prefix = tablePrefix + "."
	}

	q := query.Where(prefix+"user_id = ?", userID)

	if req.OutletID != "" {
		q = q.Where(prefix+"outlet_id = ?", req.OutletID)
	}

	if req.StartDate != "" {
		q = q.Where("DATE("+prefix+"order_date) >= ?", req.StartDate)
	}

	if req.EndDate != "" {
		q = q.Where("DATE("+prefix+"order_date) <= ?", req.EndDate)
	}

	return q
}

func (r *reportRepository) GetTotalOmzet(ctx context.Context, userID string, req dto.ReportQuery) (decimal.Decimal, error) {
	var total decimal.Decimal

	db := r.db.WithContext(ctx).Table("orders")
	q := r.buildBaseQuery(ctx, db, userID, "", req)

	err := q.Where("status != ?", "cancelled").
		Select("COALESCE(SUM(total_price), 0)").
		Scan(&total).Error

	return total, err
}

func (r *reportRepository) GetOrderStatusSummary(ctx context.Context, userID string, req dto.ReportQuery) (map[string]int64, error) {
	var rows []OrderStatusRow

	db := r.db.WithContext(ctx).Table("orders")
	q := r.buildBaseQuery(ctx, db, userID, "", req)

	err := q.Select("status, COUNT(id) as count").
		Group("status").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	summary := make(map[string]int64)
	for _, row := range rows {
		summary[row.Status] = row.Count
	}

	return summary, nil
}

func (r *reportRepository) GetTopServices(ctx context.Context, userID string, req dto.ReportQuery) ([]TopServiceRow, error) {
	var rows []TopServiceRow

	db := r.db.WithContext(ctx).Table("order_items oi")
	query := db.Joins("JOIN orders ord ON ord.id = oi.order_id").
		Joins("JOIN outlets o ON o.id = ord.outlet_id")

	// Apply filtering to "ord"
	q := r.buildBaseQuery(ctx, query, userID, "ord", req)

	err := q.Where("ord.status != ?", "cancelled").
		Select("oi.service_name, o.name as outlet_name, SUM(oi.qty) as total_qty, SUM(oi.subtotal) as total_revenue").
		Group("oi.service_name, o.name").
		Order("total_revenue DESC").
		Limit(5).
		Scan(&rows).Error

	return rows, err
}
