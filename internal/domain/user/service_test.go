package user

import (
	"context"

	"github.com/sangkips/revenue-system/internal/domain/user/models"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, user models.InsertUserParams) (models.User, error) {
	args := m.Called(ctx, user)

	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)

	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockRepository) ListUsers(ctx context.Context, params models.ListUsersParams) ([]models.ListUsersRow, error) {
	args := m.Called(ctx, params)

	return args.Get(0).([]models.ListUsersRow), args.Error(1)
}

func (m *MockRepository) ListAllUsers(ctx context.Context, params models.ListAllUsersParams) ([]models.ListAllUsersRow, error) {
	args := m.Called(ctx, params)

	return args.Get(0).([]models.ListAllUsersRow), args.Error(1)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id string) (models.GetUserByIDRow, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(models.GetUserByIDRow), args.Error(1)
}

func (m *MockRepository) UpdateUser(ctx context.Context, params models.UpdateUserParams) error {
	args := m.Called(ctx, params)

	return args.Error(0)
}

func (m *MockRepository) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}
