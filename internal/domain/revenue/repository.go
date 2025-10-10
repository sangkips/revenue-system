package revenue

import (
	"context"

	"github.com/google/uuid"
	"github.com/sangkips/revenue-system/internal/domain/revenue/models"
)

type Repository interface {
	CreateRevenue(ctx context.Context, revenue models.InsertRevenueParams) (models.Revenue, error)
	GetRevenueByID(ctx context.Context, id string) (models.Revenue, error)
	ListRevenues(ctx context.Context, params models.ListRevenuesParams) ([]models.Revenue, error)
	ListRevenuesByTaxpayerID(ctx context.Context, taxpayerID string, limit, offset int32) ([]models.Revenue, error)
	UpdateRevenue(ctx context.Context, params models.UpdateRevenueParams) (models.Revenue, error)
	DeleteRevenue(ctx context.Context, id string) error
}

type repository struct {
	q *models.Queries
}

func NewRepository(db models.DBTX) Repository {
	return &repository{q: models.New(db)}
}

func (r *repository) CreateRevenue(ctx context.Context, revenue models.InsertRevenueParams) (models.Revenue, error) {
	return r.q.InsertRevenue(ctx, revenue)
}

func (r *repository) GetRevenueByID(ctx context.Context, id string) (models.Revenue, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.Revenue{}, err
	}
	return r.q.GetRevenueByID(ctx, parsedID)
}

func (r *repository) ListRevenues(ctx context.Context, params models.ListRevenuesParams) ([]models.Revenue, error) {
	return r.q.ListRevenues(ctx, params)
}

func (r *repository) ListRevenuesByTaxpayerID(ctx context.Context, taxpayerID string, limit, offset int32) ([]models.Revenue, error) {
	parsedTaxpayerID, err := uuid.Parse(taxpayerID)
	if err != nil {
		return nil, err
	}

	params := models.ListRevenuesByTaxpayerIDParams{
		Limit:      limit,
		Offset:     offset,
		TaxpayerID: parsedTaxpayerID,
	}
	return r.q.ListRevenuesByTaxpayerID(ctx, params)
}

func (r *repository) UpdateRevenue(ctx context.Context, params models.UpdateRevenueParams) (models.Revenue, error) {
	return r.q.UpdateRevenue(ctx, params)
}

func (r *repository) DeleteRevenue(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil
	}
	return r.q.DeleteRevenue(ctx, parsedID)
}
