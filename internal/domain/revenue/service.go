package revenue

import (
	"context"
	"time"

	"github.com/sangkips/revenue-system/internal/domain/revenue/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}


func (s *Service) CreateRevenue(ctx context.Context, req CreateRevenueRequest) (models.Revenue, error) {}

func (s *Service) GetRevenue(ctx context.Context, id string) (models.Revenue, error) {}

func (s *Service) ListRevenues(ctx context.Context, countyID, limit, offset int32) ([]models.Revenue, error) {}

func (s *Service) UpdateRevenue(ctx context.Context, id string, req UpdateRevenueRequest) (models.Revenue, error) {}

func (s *Service) DeleteRevenue(ctx context.Context, id string) error {}


type CreateRevenueRequest struct {
	TaxpayerID      string    `json:"taxpayer_id"`
	CountyID        int32     `json:"county_id"`
	Amount          float64   `json:"amount"`
	RevenueType     string    `json:"revenue_type"`
	TransactionDate time.Time `json:"transaction_date"`
	Description     string    `json:"description"`
}

type UpdateRevenueRequest struct {
	Amount          float64   `json:"amount"`
	RevenueType     string    `json:"revenue_type"`
	TransactionDate time.Time `json:"transaction_date"`
	Description     string    `json:"description"`
}