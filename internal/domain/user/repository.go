package user

import (
	"context"

	"github.com/sangkips/revenue-system/internal/domain/user/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.InsertUserParams) (models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
}

type repository struct {
	q *models.Queries
}

func NewRepository(db models.DBTX) Repository {
	return &repository{q: models.New(db)}
}

func (r *repository) CreateUser(ctx context.Context, user models.InsertUserParams) (models.User, error) {
	return r.q.InsertUser(ctx, user)
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	return r.q.GetUserByUsername(ctx, username)
}
