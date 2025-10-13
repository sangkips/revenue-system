package assessment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sangkips/revenue-system/internal/domain/assessment/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateAssessment(ctx context.Context, req CreateAssessmentRequest, userID string) (models.Assessment, error) {
	if req.CountyID == 0 || req.TaxpayerID == "" || req.AssessmentNumber == "" || req.AssessmentType == "" ||
	req.FinancialYear == "" || req.BaseAmount <= 0 || req.TotalAmount <= 0 {
		return models.Assessment{}, errors.New("required fields missing or invalid")
	}

	if req.Status != "" && !validStatus(req.Status) {
		return models.Assessment{}, errors.New("invalid status value: must be 'pending', 'approved', 'rejected', or 'paid'")
	}

	if userID == "" {
		return models.Assessment{}, errors.New("user ID is required")
	}

	assessedDate := req.AssessedDate
	if assessedDate.IsZero() {
		assessedDate = time.Now()
	}

	dueDate := req.DueDate
	if dueDate.IsZero() {
		dueDate = time.Now().AddDate(0, 1,0)
	}

	status := req.Status
	if status == "" {
		status = "pending"
	}

	taxpayerID, err := uuid.Parse(req.TaxpayerID)
	if err != nil {
		return models.Assessment{}, err
	}

	var revenueID uuid.NullUUID
	if req.RevenueID != "" {
		revenueID.UUID, err = uuid.Parse(req.RevenueID)
		if err != nil {
			return models.Assessment{}, err
		}
		revenueID.Valid = true
	}

	var assessedBy uuid.NullUUID
	if req.AssessedBy != "" {
		assessedBy.UUID, err = uuid.Parse(req.AssessedBy)
		if err != nil {
			return models.Assessment{}, err
		}
	}

	params := models.InsertAssessmentParams{
		CountyID:         req.CountyID,
		TaxpayerID:       taxpayerID,
		RevenueID:        revenueID,
		AssessmentNumber: req.AssessmentNumber,
		AssessmentType:   req.AssessmentType,
		FinancialYear:    req.FinancialYear,
		BaseAmount:       fmt.Sprintf("%.2f", req.BaseAmount),
		CalculatedAmount: fmt.Sprintf("%.2f", req.CalculatedAmount),
		TotalAmount:      fmt.Sprintf("%.2f", req.TotalAmount),
		Status:           status,
		DueDate:          dueDate,
		AssessedBy:       uuid.NullUUID{UUID: uuid.MustParse(userID), Valid: true},
		AssessedDate:     assessedDate,
	}

	return s.repo.CreateAssessment(ctx, params)
}

func (s *Service) GetAssessment(ctx context.Context, id string) (models.Assessment, error) {
	return s.repo.GetAssessmentByID(ctx, id)
}

func (s *Service) ListAssessments(ctx context.Context, countyID, limit, offset int32) ([]models.Assessment, error) {
	return s.repo.ListAssessments(ctx, models.ListAssessmentsParams{
		CountyID: countyID,
		Limit: limit,
		Offset: offset,
	})
}

func (s *Service) UpdateAssessment(ctx context.Context, id string, req UpdateAssessmentRequest) (models.Assessment, error) {
	if req.BaseAmount != nil && *req.BaseAmount <= 0 {
		return models.Assessment{}, errors.New("base_amount must be greater than 0")
	}

	if req.TotalAmount != nil && *req.TotalAmount <= 0 {
		return models.Assessment{}, errors.New("total_amount must be greater than 0")
	}

	if req.Status != nil && !validStatus(*req.Status) {
		return models.Assessment{}, errors.New("invalid status value")
	}

	if req.DueDate != nil && req.DueDate.IsZero() {
		return models.Assessment{}, errors.New("due_date must be a valid date")
	}

	assessmentID, err := uuid.Parse(id)
	if err != nil {
		return models.Assessment{}, err
	}

	params := models.UpdateAssessmentParams{
		ID:               assessmentID,
		BaseAmount:       "",
		CalculatedAmount: "",
		TotalAmount:      "",
		Status:           "",
		DueDate:          time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	if req.BaseAmount != nil {
		params.BaseAmount = fmt.Sprintf("%.2f", *req.BaseAmount)
	}
	if req.CalculatedAmount != nil {
		params.CalculatedAmount = fmt.Sprintf("%.2f", *req.CalculatedAmount)
	}
	if req.TotalAmount != nil {
		params.TotalAmount = fmt.Sprintf("%.2f", *req.TotalAmount)
	}
	if req.Status != nil {
		params.Status = *req.Status
	}
	if req.DueDate != nil {
		params.DueDate = *req.DueDate
	}

	return s.repo.UpdateAssessment(ctx, params)
}

func (s *Service) DeleteAssessment(ctx context.Context, id string) error {
	return s.repo.DeleteAssessment(ctx, id)
}

func (s *Service) CreateAssessmentItem(ctx context.Context, req CreateAssessmentItemRequest) (models.AssessmentItem, error) {
	if req.AssessmentID == "" || req.ItemDescription == "" || req.UnitAmount <= 0 || req.TotalAmount <= 0 {
		return models.AssessmentItem{}, errors.New("required fields missing or invalid")
	}
	assessmentUUID, err := uuid.Parse(req.AssessmentID)
	if err != nil {
		return models.AssessmentItem{}, errors.New("invalid assessment_id format")
	}
	params := models.InsertAssessmentItemParams{
		AssessmentID:     assessmentUUID,
		ItemDescription:  req.ItemDescription,
		Quantity:         sql.NullString{String: fmt.Sprintf("%.2f", req.Quantity), Valid: true},
		UnitAmount:       fmt.Sprintf("%.2f", req.UnitAmount),
		TotalAmount:      fmt.Sprintf("%.2f", req.TotalAmount),
	}
	return s.repo.CreateAssessmentItem(ctx, params)
}

func (s *Service) ListAssessmentItems(ctx context.Context, assessmentID string) ([]models.AssessmentItem, error) {
	return s.repo.ListAssessmentItems(ctx, assessmentID)
}

func (s *Service) DeleteAssessmentItem(ctx context.Context, assessmentID, itemID string) error {
	item, err := s.repo.GetAssessmentItemByID(ctx, itemID)
	if err != nil {
		return errors.New("assessment item not found")
	}
	parsedID, err := uuid.Parse(assessmentID)
	if err != nil {
		return err
	}

	if item.AssessmentID != parsedID {
		return errors.New("assessment item does not belong to the specified assessment")
	}
	return s.repo.DeleteAssessmentItem(ctx, itemID)
}

func validStatus(status string) bool {
	return status == "pending" || status == "approved" || status == "rejected" || status == "paid"
}


type CreateAssessmentRequest struct {
	CountyID        int32     `json:"county_id"`
	TaxpayerID      string    `json:"taxpayer_id"`
	RevenueID       string    `json:"revenue_id,omitempty"`
	AssessmentNumber string   `json:"assessment_number"`
	AssessmentType  string    `json:"assessment_type"`
	FinancialYear   string    `json:"financial_year"`
	BaseAmount      float64   `json:"base_amount"`
	CalculatedAmount float64  `json:"calculated_amount"`
	TotalAmount     float64   `json:"total_amount"`
	Status          string    `json:"status,omitempty"`
	DueDate         time.Time `json:"due_date,omitempty"`
	AssessedBy      string    `json:"assessed_by,omitempty"`
	AssessedDate    time.Time `json:"assessed_date,omitempty"`
}

type UpdateAssessmentRequest struct {
	BaseAmount      *float64  `json:"base_amount,omitempty"`
	CalculatedAmount *float64 `json:"calculated_amount,omitempty"`
	TotalAmount     *float64  `json:"total_amount,omitempty"`
	Status          *string   `json:"status,omitempty"`
	DueDate         *time.Time `json:"due_date,omitempty"`
}

type CreateAssessmentItemRequest struct {
	AssessmentID    string   `json:"assessment_id"`
	ItemDescription string   `json:"item_description"`
	Quantity        float64  `json:"quantity"`
	UnitAmount      float64  `json:"unit_amount"`
	TotalAmount     float64  `json:"total_amount"`
}
