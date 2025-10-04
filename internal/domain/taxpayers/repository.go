package taxpayers

import (
	"context"

	"github.com/google/uuid"
	"github.com/sangkips/revenue-system/internal/domain/taxpayers/models"
)

type Repository interface {
	CreateTaxpayer(ctx context.Context, taxpayer models.InsertTaxpayerParams) (models.Taxpayer, error)
	GetTaxpayerByID(ctx context.Context, id string) (models.Taxpayer, error)
	GetTaxpayerByNationalID(ctx context.Context, nationalID string) (models.Taxpayer, error)
	ListTaxpayers(ctx context.Context, params models.ListTaxpayersParams) ([]models.Taxpayer, error)
	UpdateTaxpayer(ctx context.Context, params models.UpdateTaxpayerParams) (models.Taxpayer, error)
	DeleteTaxpayer(ctx context.Context, id string) error
}

type repository struct {
	q *models.Queries
}

func NewRepository(db models.DBTX) Repository {
	return &repository{q: models.New(db)}
}

func (r *repository) CreateTaxpayer(ctx context.Context, taxpayer models.InsertTaxpayerParams) (models.Taxpayer, error) {
	return r.q.InsertTaxpayer(ctx, taxpayer)
}

func (r *repository) GetTaxpayerByID(ctx context.Context, id string) (models.Taxpayer, error){
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.Taxpayer{}, err
	}
	return r.q.GetTaxpayerByID(ctx, parsedID)
}

func (r *repository) GetTaxpayerByNationalID(ctx context.Context, nationalID string) (models.Taxpayer, error) {
	return r.q.GetTaxpayerByNationalID(ctx, nationalID)
}

func (r *repository) ListTaxpayers(ctx context.Context, params models.ListTaxpayersParams) ([]models.Taxpayer, error) {
	return r.q.ListTaxpayers(ctx, params)
}

func (r *repository) UpdateTaxpayer(ctx context.Context, params models.UpdateTaxpayerParams) (models.Taxpayer, error) {
	return r.q.UpdateTaxpayer(ctx, params)
}

func (r *repository) DeleteTaxpayer(ctx context.Context, id string) error {
	parseID, err := uuid.Parse(id)
	if err != nil {
		return nil
	}
	return r.q.DeleteTaxpayer(ctx, parseID)
}
