package usecase

import (
	"context"
	"errors"

	"laundryin/internal/dto"
	"laundryin/internal/repository"
)

type ReportUsecase interface {
	GetOmzet(ctx context.Context, userID string, req dto.ReportQuery) (*dto.OmzetResponse, error)
	GetOrderStatusSummary(ctx context.Context, userID string, req dto.ReportQuery) (*dto.OrderStatusSummaryResponse, error)
	GetTopServices(ctx context.Context, userID string, req dto.ReportQuery) ([]dto.TopServiceResponse, error)
}

type reportUsecase struct {
	reportRepo repository.ReportRepository
}

func NewReportUsecase(reportRepo repository.ReportRepository) ReportUsecase {
	return &reportUsecase{reportRepo: reportRepo}
}

func (u *reportUsecase) GetOmzet(ctx context.Context, userID string, req dto.ReportQuery) (*dto.OmzetResponse, error) {
	totalDec, err := u.reportRepo.GetTotalOmzet(ctx, userID, req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal menghitung total omzet")
	}

	return &dto.OmzetResponse{
		// Force exact 2 decimal precision for stability
		TotalOmzet: totalDec.StringFixed(2),
	}, nil
}

func (u *reportUsecase) GetOrderStatusSummary(ctx context.Context, userID string, req dto.ReportQuery) (*dto.OrderStatusSummaryResponse, error) {
	summaryMap, err := u.reportRepo.GetOrderStatusSummary(ctx, userID, req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal mengambil ringkasan status pesanan")
	}

	return &dto.OrderStatusSummaryResponse{
		Pending:   summaryMap["pending"],
		Process:   summaryMap["process"],
		Completed: summaryMap["completed"],
		PickedUp:  summaryMap["picked_up"],
		Cancelled: summaryMap["cancelled"],
	}, nil
}

func (u *reportUsecase) GetTopServices(ctx context.Context, userID string, req dto.ReportQuery) ([]dto.TopServiceResponse, error) {
	rows, err := u.reportRepo.GetTopServices(ctx, userID, req)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, errors.New("gagal mengambil data layanan terlaris")
	}

	responses := make([]dto.TopServiceResponse, 0, len(rows))
	for _, r := range rows {
		responses = append(responses, dto.TopServiceResponse{
			ServiceName:  r.ServiceName,
			OutletName:   r.OutletName,
			TotalQty:     r.TotalQty.StringFixed(2),
			TotalRevenue: r.TotalRevenue.StringFixed(2),
		})
	}

	return responses, nil
}
