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
	Name            string `json:"name"`
	TreasuryAccount string `json:"treasury_account"`
}
