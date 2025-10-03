package counties

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sangkips/revenue-system/internal/domain/counties/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCounty(ctx context.Context, req CreateCountyRequest) (models.County, error) {
	if req.Name == "" || req.Code == "" {
		return models.County{}, errors.New("name and code are required")
	}

	params := models.InsertCountyParams{
		Name:            req.Name,
		Code:            req.Code,
		TreasuryAccount: sql.NullString{String: req.TreasuryAccount, Valid: req.TreasuryAccount != ""},
	}
	return s.repo.CreateCounty(ctx, params)
}

type CreateCountyRequest struct {
	Name            string `json:"name"`
	Code            string `json:"code"`
	TreasuryAccount string `json:"treasury_account"`
}

type UpdateCountyRequest struct {
	Name            *string `json:"name,omitempty"`
	TreasuryAccount *string `json:"treasury_account,omitempty"`
}

func (s *Service) GetCounty(ctx context.Context, id int32) (models.County, error) {
	return s.repo.GetCountyByID(ctx, id)
}

func (s *Service) ListCounties(ctx context.Context, limit, offset int32) ([]models.County, error) {
	return s.repo.ListCounties(ctx, models.ListCountiesParams{Limit: limit, Offset: offset})
}

func (s *Service) UpdateCounty(ctx context.Context, id int32, req UpdateCountyRequest) (models.County, error) {
	params := models.UpdateCountyParams{
		ID:                    id,
		UpdateName:            req.Name != nil,
		UpdateTreasuryAccount: req.TreasuryAccount != nil,
	}

	if req.Name != nil {
		params.Name = *req.Name
	}

	if req.TreasuryAccount != nil {
		params.TreasuryAccount = sql.NullString{String: *req.TreasuryAccount, Valid: true}
	}

	return s.repo.UpdateCounty(ctx, params)
}

func (s *Service) DeleteCounty(ctx context.Context, id int32) error {
	return s.repo.DeleteCounty(ctx, id)
}
