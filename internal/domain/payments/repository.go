package payments

import (
	"context"

	"github.com/google/uuid"
	"github.com/sangkips/revenue-system/internal/domain/payments/models"
)

type Repository interface {
	CreatePayment(ctx context.Context, payment models.InsertPaymentParams) (models.Payment, error)
	GetPaymentByID(ctx context.Context, id string) (models.Payment, error)
	ListPayments(ctx context.Context, params models.ListPaymentsParams) ([]models.Payment, error)
	ListPaymentsByRevenueID(ctx context.Context, revenueID string) ([]models.Payment, error)
	UpdatePayment(ctx context.Context, params models.UpdatePaymentParams) (models.Payment, error)
	DeletePayment(ctx context.Context, id string) error

	// Payment Allocations
	CreatePaymentAllocation(ctx context.Context, allocation models.InsertPaymentAllocationParams) (models.PaymentAllocation, error)
	ListPaymentAllocations(ctx context.Context, paymentID string) ([]models.PaymentAllocation, error)
	DeletePaymentAllocation(ctx context.Context, id string) error

	// Receipts
	CreateReceipt(ctx context.Context, receipt models.InsertReceiptParams) error
	GetReceiptByID(ctx context.Context, id string) (models.Receipt, error)
	ListReceiptsByPayment(ctx context.Context, paymentID string) ([]models.Receipt, error)
	UpdateReceipt(ctx context.Context, params models.UpdateReceiptParams) error
	DeleteReceipt(ctx context.Context, id string) error
}

type repository struct {
	q *models.Queries
}

func NewRepository(db models.DBTX) Repository {
	return &repository{q: models.New(db)}
}

func (r *repository) CreatePayment(ctx context.Context, payment models.InsertPaymentParams) (models.Payment, error) {
	return r.q.InsertPayment(ctx, payment)
}

func (r *repository) GetPaymentByID(ctx context.Context, id string) (models.Payment, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.Payment{}, err
	}
	return r.q.GetPaymentByID(ctx, parsedID)
}

func (r *repository) ListPayments(ctx context.Context, params models.ListPaymentsParams) ([]models.Payment, error) {
	return r.q.ListPayments(ctx, params)
}

func (r *repository) ListPaymentsByRevenueID(ctx context.Context, revenueID string) ([]models.Payment, error) {
	var assessmentID uuid.NullUUID
	if revenueID != "" {
		parsedID, err := uuid.Parse(revenueID)
		if err != nil {
			return nil, err
		}
		assessmentID = uuid.NullUUID{UUID: parsedID, Valid: true}
	}
	return r.q.ListPaymentsByRevenueID(ctx, assessmentID)
}

func (r *repository) UpdatePayment(ctx context.Context, params models.UpdatePaymentParams) (models.Payment, error) {
	return r.q.UpdatePayment(ctx, params)
}

func (r *repository) DeletePayment(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.DeletePayment(ctx, parsedID)
}


// Payment Allocations
func (r *repository) CreatePaymentAllocation(ctx context.Context, allocation models.InsertPaymentAllocationParams) (models.PaymentAllocation, error) {
	return r.q.InsertPaymentAllocation(ctx, allocation)
}

func (r *repository) ListPaymentAllocations(ctx context.Context, paymentID string) ([]models.PaymentAllocation, error) {
	parseID, err := uuid.Parse(paymentID)
	if err != nil {
		return []models.PaymentAllocation{}, err
	}
	return r.q.ListPaymentAllocations(ctx, parseID)
}

func (r *repository) DeletePaymentAllocation(ctx context.Context, id string) error {
	parseID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.DeletePaymentAllocation(ctx, parseID)
}

// Receipts
func (r *repository) CreateReceipt(ctx context.Context, receipt models.InsertReceiptParams) error {
	return r.q.InsertReceipt(ctx, receipt)
}

func (r *repository) GetReceiptByID(ctx context.Context, id string) (models.Receipt, error) {
	parseID, err := uuid.Parse(id)
	if err != nil {
		return models.Receipt{}, err
	}
	return r.q.GetReceiptByID(ctx, parseID)
}

func (r *repository) ListReceiptsByPayment(ctx context.Context, paymentID string) ([]models.Receipt, error) {
	parsedID, err := uuid.Parse(paymentID)
	if err != nil {
		return nil, err
	}
	return r.q.ListReceiptsByPayment(ctx, parsedID)
}

func (r *repository) UpdateReceipt(ctx context.Context, params models.UpdateReceiptParams) error {
	return r.q.UpdateReceipt(ctx, params)
}

func (r *repository) DeleteReceipt(ctx context.Context, id string) error {
	parseID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.DeleteReceipt(ctx, parseID)
}