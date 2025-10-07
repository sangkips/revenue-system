package revenue

import (
	"context"

	"github.com/sangkips/revenue-system/internal/domain/revenue/models"
)

type Repository interface {
	CreateRevenue(ctx context.Context, revenue models.InsertRevenueParams) (models.Revenue, error)
	GetRevenueByID(ctx context.Context, id string) (models.Revenue, error)
	ListRevenues(ctx context.Context, params models.ListRevenuesParams) ([]models.Revenue, error)
	UpdateRevenue(ctx context.Context, params models.UpdateRevenueParams) (models.Revenue, error)
	DeleteRevenue(ctx context.Context, id string) error
}

type repository struct {
	q *models.Queries
}

func NewRepository(db models.DBTX) Repository {
	return &repository{q: models.New(db)}
}

func (r *repository) CreateRevenue(ctx context.Context, revenue models.InsertRevenueParams) (models.Revenue, error) {}

func (r *repository) GetRevenueByID(ctx context.Context, id string) (models.Revenue, error) {}

func (r *repository) ListRevenues(ctx context.Context, params models.ListRevenuesParams) (models.Revenue, error) {}

func (r *repository) UpdateRevenue(ctx context.Context, params models.UpdateRevenueParams) (models.Revenue, error) {}

func (r *repository) DeleteRevenue(ctx context.Context, id string) error {}