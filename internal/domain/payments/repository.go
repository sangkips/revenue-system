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