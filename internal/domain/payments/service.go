package payments

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sangkips/revenue-system/internal/domain/payments/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePayment(ctx context.Context, req CreatePaymentRequest, userID string) (models.Payment, error) {
	if req.CountyID == 0 || req.TaxpayerID == "" || req.PaymentNumber == "" || req.Amount <= 0 || req.PaymentMethod == "" {
		return models.Payment{}, errors.New("required fields missing or invalid")
	}

	if req.Status != "" && !validStatus(req.Status) {
		return models.Payment{}, errors.New("invalid status value: must be 'pending', 'processing', 'completed', 'failed', or 'cancelled'")
	}

	if req.PaymentMethod != "" && !validPaymentMethod(req.PaymentMethod) {
		return models.Payment{}, errors.New("invalid payment_method: must be 'mpesa', 'bank_transfer', 'card', 'cheque', or 'cash'")
	}

	if userID == "" {
		return models.Payment{}, errors.New("user ID is required")
	}

	status := req.Status
	if status == "" {
		status = "pending"
	}

	taxpayerID, err := uuid.Parse(req.TaxpayerID)
	if err != nil {
		return models.Payment{}, err
	}

	paymentDate := req.PaymentDate
	if paymentDate.IsZero() {
		paymentDate = time.Now()
	}


	var assessmentID uuid.NullUUID
	if req.AssessmentID != "" {
		parsedAssessmentID, err := uuid.Parse(req.AssessmentID)
		if err != nil {
			return models.Payment{}, err
		}
		assessmentID = uuid.NullUUID{UUID: parsedAssessmentID, Valid: true}
	}

	params := models.InsertPaymentParams{
		CountyID:             req.CountyID,
		TaxpayerID:           taxpayerID,
		AssessmentID:         assessmentID,
		PaymentNumber:        req.PaymentNumber,
		Amount:               fmt.Sprintf("%.2f",req.Amount),
		PaymentMethod:        req.PaymentMethod,
		PaymentChannel:       sql.NullString{String: req.PaymentChannel, Valid: req.PaymentChannel != ""},
		ExternalTransactionID: sql.NullString{String: req.ExternalTransactionID, Valid: req.ExternalTransactionID != ""},
		PayerPhoneNumber:     sql.NullString{String: req.PayerPhoneNumber, Valid: req.PayerPhoneNumber != ""},
		PayerName:            sql.NullString{String: req.PayerName, Valid: req.PayerName != ""},
		Status:               status,
		CollectedBy:          uuid.NullUUID{UUID: uuid.MustParse(userID), Valid: true},
	}
	return s.repo.CreatePayment(ctx, params)
}


func (s *Service) GetPayment(ctx context.Context, id string) (models.Payment, error) {
	return s.repo.GetPaymentByID(ctx, id)
}

func (s *Service) ListPayments(ctx context.Context, countyID int32, limit int32, offset int32) ([]models.Payment, error) {
	return s.repo.ListPayments(ctx, models.ListPaymentsParams{
		CountyID: countyID,
		Limit:    limit,
		Offset:   offset,
	})
}

func (s *Service) ListPaymentsByRevenueID(ctx context.Context, revenueID string) ([]models.Payment, error) {
	return s.repo.ListPaymentsByRevenueID(ctx, revenueID)
}

func (s *Service) UpdatePayment(ctx context.Context, id string, req UpdatePaymentRequest, userID string) (models.Payment, error) {
	if req.Status != nil && !validStatus(*req.Status) {
		return models.Payment{}, errors.New("invalid status value: must be 'pending', 'processing', 'completed', 'failed', or 'cancelled'")
	}

	if req.PaymentMethod != nil && !validPaymentMethod(*req.PaymentMethod) {
		return models.Payment{}, errors.New("invalid payment_method: must be 'mpesa', 'bank_transfer', 'card', 'cheque', or 'cash'")
	}

	if userID == "" {
		return models.Payment{}, errors.New("user ID is required")
	}

	paymentID, err := uuid.Parse(id)
	if err != nil {
		return models.Payment{}, err
	}

	// Get current payment to use existing values for fields not being updated
	current, err := s.repo.GetPaymentByID(ctx, id)
	if err != nil {
		return models.Payment{}, err
	}

	params := models.UpdatePaymentParams{
		ID:                    paymentID,
		Amount:                current.Amount,
		PaymentMethod:         current.PaymentMethod,
		PaymentChannel:        current.PaymentChannel,
		ExternalTransactionID: current.ExternalTransactionID,
		PayerPhoneNumber:      current.PayerPhoneNumber,
		PayerName:             current.PayerName,
		Status:                current.Status,
		CollectedBy:           current.CollectedBy,
	}

	if req.Amount != nil {
		params.Amount = fmt.Sprintf("%.2f", *req.Amount)
	}
	if req.PaymentMethod != nil {
		params.PaymentMethod = *req.PaymentMethod
	}
	if req.PaymentChannel != nil {
		params.PaymentChannel = sql.NullString{String: *req.PaymentChannel, Valid: true}
	}
	if req.ExternalTransactionID != nil {
		params.ExternalTransactionID = sql.NullString{String: *req.ExternalTransactionID, Valid: true}
	}
	if req.PayerPhoneNumber != nil {
		params.PayerPhoneNumber = sql.NullString{String: *req.PayerPhoneNumber, Valid: true}
	}
	if req.PayerName != nil {
		params.PayerName = sql.NullString{String: *req.PayerName, Valid: true}
	}
	if req.Status != nil {
		params.Status = *req.Status
	}
	if req.CollectedBy != nil {
		collectedByUUID, err := uuid.Parse(*req.CollectedBy)
		if err != nil {
			return models.Payment{}, err
		}
		params.CollectedBy = uuid.NullUUID{UUID: collectedByUUID, Valid: true}
	}

	return s.repo.UpdatePayment(ctx, params)
}

func (s *Service) DeletePayment(ctx context.Context, id string) error {
	return s.repo.DeletePayment(ctx, id)
}

func validPaymentMethod(method string) bool {
	return method == "mpesa" || method == "bank_transfer" || method == "card" || method == "cheque" || method == "cash"
}

func validStatus(status string) bool {
	return status == "pending" || status == "processing" || status == "completed" || status == "failed" || status == "cancelled"
}


type CreatePaymentRequest struct {
	CountyID             int32   `json:"county_id"`
	TaxpayerID           string  `json:"taxpayer_id"`
	AssessmentID         string  `json:"assessment_id,omitempty"`
	PaymentNumber        string  `json:"payment_number"`
	Amount               float64 `json:"amount"`
	PaymentMethod        string  `json:"payment_method"`
	PaymentChannel       string  `json:"payment_channel,omitempty"`
	ExternalTransactionID string `json:"external_transaction_id,omitempty"`
	PayerPhoneNumber     string  `json:"payer_phone_number,omitempty"`
	PayerName            string  `json:"payer_name,omitempty"`
	Status               string  `json:"status,omitempty"`
	PaymentDate          time.Time `json:"payment_date,omitempty"`
}

type UpdatePaymentRequest struct {
	Amount               *float64 `json:"amount,omitempty"`
	PaymentMethod        *string  `json:"payment_method,omitempty"`
	PaymentChannel       *string  `json:"payment_channel,omitempty"`
	ExternalTransactionID *string  `json:"external_transaction_id,omitempty"`
	PayerPhoneNumber     *string  `json:"payer_phone_number,omitempty"`
	PayerName            *string  `json:"payer_name,omitempty"`
	Status               *string  `json:"status,omitempty"`
	CollectedBy          *string  `json:"collected_by,omitempty"`
}