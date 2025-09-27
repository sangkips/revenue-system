package counties

import (
	"context"

	"github.com/sangkips/revenue-system/internal/domain/counties/models"
)

type Repository interface {
	CreateCounty(ctx context.Context, county models.InsertCountyParams) error
}

type repository struct {
	q *models.Queries
}

func NewRepository(db models.DBTX) Repository {
	return &repository{q: models.New(db)}
}

func (r *repository) CreateCounty(ctx context.Context, county models.InsertCountyParams) error {
	return r.q.InsertCounty(ctx, county)
}
