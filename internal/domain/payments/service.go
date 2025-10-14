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
		CollectedBy:          func() uuid.NullUUID {
			if userID != "" {
				if parsedUserID, err := uuid.Parse(userID); err == nil {
					return uuid.NullUUID{UUID: parsedUserID, Valid: true}
				}
			}
			return uuid.NullUUID{Valid: false}
		}(),
		MpesaReceiptNumber:   sql.NullString{String: req.MpesaReceiptNumber, Valid: req.MpesaReceiptNumber != ""},
		BankReference:        sql.NullString{String: req.BankReference, Valid: req.BankReference != ""},
		ChequeNumber:         sql.NullString{String: req.ChequeNumber, Valid: req.ChequeNumber != ""},
		FailureReason:        sql.NullString{String: req.FailureReason, Valid: req.FailureReason != ""},
		CollectionPoint:      sql.NullString{String: req.CollectionPoint, Valid: req.CollectionPoint != ""},
		GpsCoordinates:       func() interface{} {
			if req.GPSCoordinates == "" {
				return nil
			}
			return req.GPSCoordinates
		}(),
		BlockchainHash:       sql.NullString{String: req.BlockchainHash, Valid: req.BlockchainHash != ""},
		BlockNumber:          sql.NullInt64{Int64: req.BlockNumber, Valid: req.BlockNumber != 0},
		Reconciled:           sql.NullBool{Bool: req.Reconciled, Valid: true},
		ReconciliationDate:   sql.NullTime{Time: req.ReconciliationDate, Valid: !req.ReconciliationDate.IsZero()},
		ReconciledBy:         func() uuid.NullUUID {
			if req.ReconciledBy != "" {
				if parsedReconciledBy, err := uuid.Parse(req.ReconciledBy); err == nil {
					return uuid.NullUUID{UUID: parsedReconciledBy, Valid: true}
				}
			}
			return uuid.NullUUID{Valid: false}
		}(),
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
		MpesaReceiptNumber:    current.MpesaReceiptNumber,
		BankReference:         current.BankReference,
		ChequeNumber:          current.ChequeNumber,
		FailureReason:         current.FailureReason,
		CollectionPoint:       current.CollectionPoint,
		GpsCoordinates:        current.GpsCoordinates,
		BlockchainHash:        current.BlockchainHash,
		BlockNumber:           current.BlockNumber,
		Reconciled:            current.Reconciled,
		ReconciliationDate:    current.ReconciliationDate,
		ReconciledBy:          current.ReconciledBy,
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
	if req.GPSCoordinates != nil {
		if *req.GPSCoordinates == "" {
			params.GpsCoordinates = nil
		} else {
			params.GpsCoordinates = *req.GPSCoordinates
		}
	}

	return s.repo.UpdatePayment(ctx, params)
}

func (s *Service) DeletePayment(ctx context.Context, id string) error {
	return s.repo.DeletePayment(ctx, id)
}


// Payment Allocations
func (s *Service) CreatePaymentAllocation(ctx context.Context, req CreatePaymentAllocationRequest) (models.PaymentAllocation, error) {
	if req.PaymentID == "" || req.AssessmentID == "" || req.AllocatedAmount <= 0 {
		return models.PaymentAllocation{}, errors.New("required fields missing or invalid")
	}
	if req.AllocationType != "" && !validAllocationType(req.AllocationType) {
		return models.PaymentAllocation{}, errors.New("invalid allocation_type")
	}
	paymentUUID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		return models.PaymentAllocation{}, err
	}
	assessmentUUID, err := uuid.Parse(req.AssessmentID)
	if err != nil {
		return models.PaymentAllocation{}, err
	}
	
	params := models.InsertPaymentAllocationParams{
		PaymentID:        paymentUUID,
		AssessmentID:     assessmentUUID,
		AllocatedAmount:  fmt.Sprintf("%.2f", req.AllocatedAmount),
		AllocationType:   sql.NullString{String: req.AllocationType, Valid: req.AllocationType != ""},
	}
	return s.repo.CreatePaymentAllocation(ctx, params)
}

func (s *Service) ListPaymentAllocations(ctx context.Context, paymentID string) ([]models.PaymentAllocation, error) {
	return s.repo.ListPaymentAllocations(ctx, paymentID)
}

func (s *Service) DeletePaymentAllocation(ctx context.Context, id string, paymentID string) error {
	return s.repo.DeletePaymentAllocation(ctx, id)
}


// Receipts
func (s *Service) CreateReceipt(ctx context.Context, req CreateReceiptRequest) error {
	if req.PaymentID == "" || req.ReceiptNumber == "" || req.ReceiptType == "" {
		return errors.New("required fields missing or invalid")
	}
	if !validReceiptType(req.ReceiptType) {
		return errors.New("invalid receipt_type")
	}
	paymentUUID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		return err
	}
	
	params := models.InsertReceiptParams{
		PaymentID:          paymentUUID,
		ReceiptNumber:      req.ReceiptNumber,
		ReceiptType:        sql.NullString{String: req.ReceiptType, Valid: req.ReceiptType != ""},
		PdfFilePath:        sql.NullString{String: req.PDFFilePath, Valid: req.PDFFilePath != ""},
		PdfFileSize:        sql.NullInt32{Int32: int32(req.PDFFileSize), Valid: req.PDFFileSize != 0},
		PdfGenerated:       sql.NullBool{Bool: req.PDFGenerated, Valid: true},
		SmsSent:            sql.NullBool{Bool: req.SMSSent, Valid: true},
		SmsSentAt:          sql.NullTime{Time: req.SMSSentAt, Valid: !req.SMSSentAt.IsZero()},
		EmailSent:          sql.NullBool{Bool: req.EmailSent, Valid: true},
		EmailSentAt:        sql.NullTime{Time: req.EmailSentAt, Valid: !req.EmailSentAt.IsZero()},
		BlockchainHash:     req.BlockchainHash,  // Required, non-null per schema
		BlockNumber:        sql.NullInt64{Int64: req.BlockNumber, Valid: req.BlockNumber != 0},
		BlockchainVerified: sql.NullBool{Bool: req.BlockchainVerified, Valid: true},
		QrCodeData:         sql.NullString{String: req.QRCodeData, Valid: req.QRCodeData != ""},
	}
	return s.repo.CreateReceipt(ctx, params)
}

func (s *Service) GetReceipt(ctx context.Context, id string) (models.Receipt, error) {
	return s.repo.GetReceiptByID(ctx, id)
}

func (s *Service) ListReceiptsByPayment(ctx context.Context, paymentID string) ([]models.Receipt, error) {
	return s.repo.ListReceiptsByPayment(ctx, paymentID)
}

func (s *Service) UpdateReceipt(ctx context.Context, id string, req UpdateReceiptRequest) error {
	if req.ReceiptType != nil && *req.ReceiptType != "" && !validReceiptType(*req.ReceiptType) {
		return errors.New("invalid receipt_type if provided")
	}
	receiptUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	
	params := models.UpdateReceiptParams{
		ID:                 receiptUUID,
		ReceiptType:        sql.NullString{Valid: req.ReceiptType != nil, String: ""},
		PdfFilePath:        sql.NullString{Valid: req.PDFFilePath != nil, String: ""},
		PdfFileSize:        sql.NullInt32{Valid: req.PDFFileSize != nil, Int32: 0},
		PdfGenerated:       sql.NullBool{Valid: req.PDFGenerated != nil, Bool: false},
		SmsSent:            sql.NullBool{Valid: req.SMSSent != nil, Bool: false},
		SmsSentAt:          sql.NullTime{Valid: req.SMSSentAt != nil, Time: time.Time{}},
		EmailSent:          sql.NullBool{Valid: req.EmailSent != nil, Bool: false},
		EmailSentAt:        sql.NullTime{Valid: req.EmailSentAt != nil, Time: time.Time{}},
		BlockchainVerified: sql.NullBool{Valid: req.BlockchainVerified != nil, Bool: false},
	}
	// Set values if provided (similar to other fields)
	if req.ReceiptType != nil {
		params.ReceiptType = sql.NullString{Valid: true, String: *req.ReceiptType}
	}
	if req.PDFFilePath != nil {
		params.PdfFilePath = sql.NullString{Valid: true, String: *req.PDFFilePath}
	}
	if req.PDFFileSize != nil {
		params.PdfFileSize = sql.NullInt32{Valid: true, Int32: int32(*req.PDFFileSize)}
	}
	if req.PDFGenerated != nil {
		params.PdfGenerated = sql.NullBool{Valid: true, Bool: *req.PDFGenerated}
	}
	if req.SMSSent != nil {
		params.SmsSent = sql.NullBool{Valid: true, Bool: *req.SMSSent}
	}
	if req.SMSSentAt != nil {
		params.SmsSentAt = sql.NullTime{Valid: true, Time: *req.SMSSentAt}
	}
	if req.EmailSent != nil {
		params.EmailSent = sql.NullBool{Valid: true, Bool: *req.EmailSent}
	}

	return s.repo.UpdateReceipt(ctx, params)
}

func (s *Service) DeleteReceipt(ctx context.Context, id string) error {
	return s.repo.DeleteReceipt(ctx, id)
}

func validPaymentMethod(method string) bool {
	return method == "mpesa" || method == "bank_transfer" || method == "card" || method == "cheque" || method == "cash"
}

func validStatus(status string) bool {
	return status == "pending" || status == "processing" || status == "completed" || status == "failed" || status == "cancelled"
}

func validReceiptType(typ string) bool {
	return typ == "payment" || typ == "provisional" || typ == "official"
}

func validAllocationType(typ string) bool {
	return typ == "principal" || typ == "penalty" || typ == "interest"
}

type CreatePaymentRequest struct {
	CountyID               int32   `json:"county_id"`
	TaxpayerID             string  `json:"taxpayer_id"`
	AssessmentID           string  `json:"assessment_id,omitempty"`
	PaymentNumber          string  `json:"payment_number"`
	Amount                 float64 `json:"amount"`
	PaymentMethod          string  `json:"payment_method"`
	PaymentChannel         string  `json:"payment_channel,omitempty"`
	ExternalTransactionID  string  `json:"external_transaction_id,omitempty"`
	MpesaReceiptNumber     string  `json:"mpesa_receipt_number,omitempty"`
	BankReference          string  `json:"bank_reference,omitempty"`
	ChequeNumber           string  `json:"cheque_number,omitempty"`
	PayerPhoneNumber       string  `json:"payer_phone_number,omitempty"`
	PayerName              string  `json:"payer_name,omitempty"`
	Status                 string  `json:"status,omitempty"`
	PaymentDate            time.Time `json:"payment_date,omitempty"`
	FailureReason          string  `json:"failure_reason,omitempty"`
	CollectionPoint        string  `json:"collection_point,omitempty"`
	GPSCoordinates         string  `json:"gps_coordinates,omitempty"`
	BlockchainHash         string  `json:"blockchain_hash,omitempty"`
	BlockNumber            int64   `json:"block_number,omitempty"`
	Reconciled             bool    `json:"reconciled,omitempty"`
	ReconciliationDate     time.Time `json:"reconciliation_date,omitempty"`
	ReconciledBy           string  `json:"reconciled_by,omitempty"`
}

type UpdatePaymentRequest struct {
	Amount                 *float64 `json:"amount,omitempty"`
	PaymentMethod          *string  `json:"payment_method,omitempty"`
	PaymentChannel         *string  `json:"payment_channel,omitempty"`
	ExternalTransactionID  *string  `json:"external_transaction_id,omitempty"`
	MpesaReceiptNumber     *string  `json:"mpesa_receipt_number,omitempty"`
	BankReference          *string  `json:"bank_reference,omitempty"`
	ChequeNumber           *string  `json:"cheque_number,omitempty"`
	PayerPhoneNumber       *string  `json:"payer_phone_number,omitempty"`
	PayerName              *string  `json:"payer_name,omitempty"`
	Status                 *string  `json:"status,omitempty"`
	FailureReason          *string  `json:"failure_reason,omitempty"`
	CollectionPoint        *string  `json:"collection_point,omitempty"`
	GPSCoordinates         *string  `json:"gps_coordinates,omitempty"`
	BlockchainHash         *string  `json:"blockchain_hash,omitempty"`
	BlockNumber            *int64   `json:"block_number,omitempty"`
	Reconciled             *bool    `json:"reconciled,omitempty"`
	ReconciliationDate     *time.Time `json:"reconciliation_date,omitempty"`
	ReconciledBy           *string  `json:"reconciled_by,omitempty"`
	CollectedBy            *string  `json:"collected_by,omitempty"`
}

type CreatePaymentAllocationRequest struct {
	PaymentID        string  `json:"payment_id"`
	AssessmentID     string  `json:"assessment_id"`
	AllocatedAmount  float64 `json:"allocated_amount"`
	AllocationType   string  `json:"allocation_type,omitempty"`
}

type CreateReceiptRequest struct {
	PaymentID          string  `json:"payment_id"`
	ReceiptNumber      string  `json:"receipt_number"`
	ReceiptType        string  `json:"receipt_type"`
	PDFFilePath        string  `json:"pdf_file_path,omitempty"`
	PDFFileSize        int     `json:"pdf_file_size,omitempty"`
	PDFGenerated       bool    `json:"pdf_generated,omitempty"`
	SMSSent            bool    `json:"sms_sent,omitempty"`
	SMSSentAt          time.Time `json:"sms_sent_at,omitempty"`
	EmailSent          bool    `json:"email_sent,omitempty"`
	EmailSentAt        time.Time `json:"email_sent_at,omitempty"`
	BlockchainHash     string  `json:"blockchain_hash"`
	BlockNumber        int64   `json:"block_number,omitempty"`
	BlockchainVerified bool    `json:"blockchain_verified,omitempty"`
	QRCodeData         string  `json:"qr_code_data,omitempty"`
}

type UpdateReceiptRequest struct {
	ReceiptType        *string  `json:"receipt_type,omitempty"`
	PDFFilePath        *string  `json:"pdf_file_path,omitempty"`
	PDFFileSize        *int     `json:"pdf_file_size,omitempty"`
	PDFGenerated       *bool    `json:"pdf_generated,omitempty"`
	SMSSent            *bool    `json:"sms_sent,omitempty"`
	SMSSentAt          *time.Time `json:"sms_sent_at,omitempty"`
	EmailSent          *bool    `json:"email_sent,omitempty"`
	EmailSentAt        *time.Time `json:"email_sent_at,omitempty"`
	BlockchainVerified *bool    `json:"blockchain_verified,omitempty"`
}