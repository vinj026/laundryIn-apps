package dto

// ReportQuery holds the optional query parameters for all report endpoints.
type ReportQuery struct {
	StartDate string `form:"start_date"` // YYYY-MM-DD
	EndDate   string `form:"end_date"`   // YYYY-MM-DD
	OutletID  string `form:"outlet_id" binding:"omitempty,uuid"`
}

// OmzetResponse is the response format for total revenue.
type OmzetResponse struct {
	TotalOmzet string `json:"total_omzet"`
}

// OrderStatusSummaryResponse represents the count of orders for each status.
type OrderStatusSummaryResponse struct {
	Pending   int64 `json:"pending"`
	Process   int64 `json:"process"`
	Completed int64 `json:"completed"`
	PickedUp  int64 `json:"picked_up"`
	Cancelled int64 `json:"cancelled"`
}

// TopServiceResponse represents one record in the top 5 services list.
type TopServiceResponse struct {
	ServiceName  string `json:"service_name"`
	OutletName   string `json:"outlet_name"`
	TotalQty     string `json:"total_qty"`
	TotalRevenue string `json:"total_revenue"`
}
