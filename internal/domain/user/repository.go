package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/sangkips/revenue-system/internal/domain/user/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.InsertUserParams) (models.User, error)
	GetUserByUsername(ctx context.Context, username string) (models.User, error)
	ListUsers(ctx context.Context, params models.ListUsersParams) ([]models.ListUsersRow, error)
	ListAllUsers(ctx context.Context, params models.ListAllUsersParams) ([]models.ListAllUsersRow, error)
	GetUserByID(ctx context.Context, id string) (models.GetUserByIDRow, error)
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

func (r *repository) ListUsers(ctx context.Context, params models.ListUsersParams) ([]models.ListUsersRow, error) {
	return r.q.ListUsers(ctx, params)
}

func (r *repository) ListAllUsers(ctx context.Context, params models.ListAllUsersParams) ([]models.ListAllUsersRow, error) {
	return r.q.ListAllUsers(ctx, params)
}

func (r *repository) GetUserByID(ctx context.Context, id string) (models.GetUserByIDRow, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.GetUserByIDRow{}, err
	}
	return r.q.GetUserByID(ctx, parsedID)
}
