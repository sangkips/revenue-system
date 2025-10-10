package revenue

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sangkips/revenue-system/internal/domain/revenue/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}


func (s *Service) CreateRevenue(ctx context.Context, req CreateRevenueRequest) (models.Revenue, error) {
	if req.TaxpayerID == "" || req.CountyID == 0 || req.Amount <= 0 || req.RevenueType == "" || req.TransactionDate.IsZero() {
		return models.Revenue{}, errors.New("taxpayer_id, county_id, amount, revenue_type, and transaction_date are required")
	}

	if req.RevenueType != "tax" && req.RevenueType != "fee" && req.RevenueType != "fine" {
		return models.Revenue{}, errors.New("revenue_type must be 'tax', 'fee', or 'fine'")
	}

	taxpayerID , err := uuid.Parse(req.TaxpayerID)
	if err != nil {
		return models.Revenue{}, err
	}

	params := models.InsertRevenueParams{
		TaxpayerID:      taxpayerID,
		CountyID:        req.CountyID,
		Amount:          fmt.Sprintf("%.2f", req.Amount),
		RevenueType:     req.RevenueType,
		TransactionDate: req.TransactionDate, // To research this date furher
		Description:     sql.NullString{String: req.Description, Valid: req.Description != ""},
	}

	return s.repo.CreateRevenue(ctx, params)
}

func (s *Service) GetRevenue(ctx context.Context, id string) (models.Revenue, error) {
	return s.repo.GetRevenueByID(ctx, id)
}

func (s *Service) ListRevenues(ctx context.Context, countyID, limit, offset int32) ([]models.Revenue, error) {
	return s.repo.ListRevenues(ctx, models.ListRevenuesParams{
		CountyID: countyID,
		Limit: limit,
		Offset: offset,
	})
}

func (s *Service) ListRevenuesByTaxpayerID(ctx context.Context, taxpayerID string, limit, offset int32) ([]models.Revenue, error) {
	return s.repo.ListRevenuesByTaxpayerID(ctx, taxpayerID, limit, offset)
}

func (s *Service) UpdateRevenue(ctx context.Context, id string, req UpdateRevenueRequest) (models.Revenue, error) {
	if req.Amount != nil && *req.Amount <= 0 {
		return models.Revenue{}, errors.New("amount must be greater than 0 if provided")
	}
	if req.RevenueType != nil && (*req.RevenueType != "tax" && *req.RevenueType != "fee" && *req.RevenueType != "fine") {
		return models.Revenue{}, errors.New("revenue_type must be 'tax', 'fee', or 'fine' if provided")
	}
	if req.TransactionDate != nil && req.TransactionDate.IsZero() {
		return models.Revenue{}, errors.New("transaction_date must be a valid date if provided")
	}

	revenueID, err := uuid.Parse(id)
	if err != nil {
		return models.Revenue{}, err
	}
	params := models.UpdateRevenueParams{
		ID:              revenueID,
		Amount:          sql.NullString{Valid: req.Amount != nil, String: ""},
		RevenueType:     sql.NullString{Valid: req.RevenueType != nil, String: ""},
		TransactionDate: sql.NullTime{Valid: req.TransactionDate != nil, Time: time.Time{}},
		Description:     sql.NullString{Valid: req.Description != nil, String: ""},
	}
	if req.Amount != nil {
		params.Amount = sql.NullString{Valid: true, String: fmt.Sprintf("%.2f", *req.Amount)}
	}
	if req.RevenueType != nil {
		params.RevenueType = sql.NullString{Valid: true, String: *req.RevenueType}
	}
	if req.TransactionDate != nil {
		params.TransactionDate = sql.NullTime{Valid: true, Time: *req.TransactionDate}
	}
	if req.Description != nil {
		params.Description = sql.NullString{Valid: true, String: *req.Description}
	}

	return s.repo.UpdateRevenue(ctx, params)
}

func (s *Service) DeleteRevenue(ctx context.Context, id string) error {
	return s.repo.DeleteRevenue(ctx, id)
}


type CreateRevenueRequest struct {
	TaxpayerID      string    `json:"taxpayer_id"`
	CountyID        int32     `json:"county_id"`
	Amount          float64   `json:"amount"`
	RevenueType     string    `json:"revenue_type"`
	TransactionDate time.Time `json:"transaction_date"`
	Description     string    `json:"description"`
}

type UpdateRevenueRequest struct {
	Amount          *float64  `json:"amount,omitempty"`
	RevenueType     *string   `json:"revenue_type,omitempty"`
	TransactionDate *time.Time `json:"transaction_date,omitempty"`
	Description     *string   `json:"description,omitempty"`
}
