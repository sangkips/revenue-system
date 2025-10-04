package taxpayers

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/sangkips/revenue-system/internal/domain/taxpayers/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTaxpayer(ctx context.Context, req CreateTaxpayerRequest) (models.Taxpayer, error) {
	if req.CountyID == 0 || req.TaxpayerType == "" || req.NationalID == "" || req.PhoneNumber == "" {
		return models.Taxpayer{}, errors.New("county_id, taxpayer_type, national_id, and phone number are required")
	}
	if req.TaxpayerType != "individual" && req.TaxpayerType != "business" {
		return models.Taxpayer{}, errors.New("taxpayer_type must be 'individual' or 'business'")
	}
	if req.TaxpayerType == "individual" && (req.FirstName == "" || req.LastName == "") {
		return models.Taxpayer{}, errors.New("first_name and last_name are required for individual taxpayers")
	}
	if req.TaxpayerType == "business" && req.BusinessName == "" {
		return models.Taxpayer{}, errors.New("business_name is required for business taxpayers")
	}

	params := models.InsertTaxpayerParams{
		CountyID:     req.CountyID,
		TaxpayerType: req.TaxpayerType,
		NationalID:   req.NationalID,
		Email:        req.Email,
		PhoneNumber:  sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
		FirstName:    sql.NullString{String: req.FirstName, Valid: req.FirstName != ""},
		LastName:     sql.NullString{String: req.LastName, Valid: req.LastName != ""},
		BusinessName: sql.NullString{String: req.BusinessName, Valid: req.BusinessName != ""},
	}

	return s.repo.CreateTaxpayer(ctx, params)
}

func (s *Service) ListTaxpayers(ctx context.Context, countyID, limit, offset int32) ([]models.Taxpayer, error) {
	return s.repo.ListTaxpayers(ctx, models.ListTaxpayersParams{
		CountyID: countyID,
		Limit:    limit,
		Offset:   offset,
	})
}


func (s *Service) GetTaxpayer(ctx context.Context, id string) (models.Taxpayer, error) {
	return s.repo.GetTaxpayerByID(ctx, id)
}


func (s *Service) UpdateTaxpayer(ctx context.Context, id string, req UpdateTaxpayerRequest) (models.Taxpayer, error) {
	taxpayerID, err := uuid.Parse(id)
	if err != nil {
		return models.Taxpayer{}, err
	}

	params := models.UpdateTaxpayerParams{
		ID:           taxpayerID,
		Email:        req.Email,
		PhoneNumber:  sql.NullString{String: req.PhoneNumber, Valid: req.PhoneNumber != ""},
		FirstName:    sql.NullString{String: req.FirstName, Valid: req.FirstName != ""},
		LastName:     sql.NullString{String: req.LastName, Valid: req.LastName != ""},
		BusinessName: sql.NullString{String: req.BusinessName, Valid: req.BusinessName != ""},
	}

	return s.repo.UpdateTaxpayer(ctx, params)
}


func (s *Service) DeleteTaxpayer(ctx context.Context, id string) error {
	return s.repo.DeleteTaxpayer(ctx, id)
}


type CreateTaxpayerRequest struct {
	CountyID     int32  `json:"county_id"`
	TaxpayerType string `json:"taxpayer_type"`
	NationalID   string `json:"national_id"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	BusinessName string `json:"business_name"`
}

type UpdateTaxpayerRequest struct {
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	BusinessName string `json:"business_name"`
}
