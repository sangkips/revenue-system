package counties

import (
	"context"

	"github.com/sangkips/revenue-system/internal/domain/counties/models"
)

type Repository interface {
	CreateCounty(ctx context.Context, county models.InsertCountyParams) (models.County, error)
	GetCountyByID(ctx context.Context, id int32) (models.County, error)
	ListCounties(ctx context.Context, params models.ListCountiesParams) ([]models.County, error)
	UpdateCounty(ctx context.Context, params models.UpdateCountyParams) (models.County, error)
	DeleteCounty(ctx context.Context, id int32) error
}

type repository struct {
	q *models.Queries
}

func NewRepository(db models.DBTX) Repository {
	return &repository{q: models.New(db)}
}

func (r *repository) CreateCounty(ctx context.Context, county models.InsertCountyParams) (models.County, error) {
	return r.q.InsertCounty(ctx, county)
}

func (r *repository) ListCounties(ctx context.Context, params models.ListCountiesParams) ([]models.County, error) {
	return r.q.ListCounties(ctx, params)
}

func (r *repository) GetCountyByID(ctx context.Context, id int32) (models.County, error) {
	return r.q.GetCountyByID(ctx, id)
}

func (r *repository) UpdateCounty(ctx context.Context, params models.UpdateCountyParams) (models.County, error) {
	return r.q.UpdateCounty(ctx, params)
}

func (r *repository) DeleteCounty(ctx context.Context, id int32) error {
	
	return r.q.DeleteCounty(ctx, id)
}
