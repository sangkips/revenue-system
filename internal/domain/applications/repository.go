package applications

import (
	"context"

	"github.com/google/uuid"
	"github.com/sangkips/revenue-system/internal/domain/applications/models"
)


type Repository interface {
	CreateApplication(ctx context.Context, application models.CreateApplicationParams) (models.Application, error)
	GetApplicationByID(ctx context.Context, id uuid.UUID) (models.GetApplicationByIDRow, error)
	UpdateApplicationStatus(ctx context.Context, params models.ListApplicationsByTaxpayerRow) error
	// DeleteApplication(ctx context.Context, id uuid.UUID) error
	CreateSingleBusinessPermit(ctx context.Context, params models.CreateSingleBusinessPermitParams) error
	CreateBuildingApproval(ctx context.Context, params models.CreateBuildingApprovalParams) error
	CreateSeasonalParkingTicket(ctx context.Context, params models.CreateSeasonalParkingTicketParams) error
	CreateHealthCertificate(ctx context.Context, params models.CreateHealthCertificateParams) error
	CreateApplicationDocument(ctx context.Context, params models.CreateApplicationDocumentParams) error
	ListApplicationsByTaxpayer(ctx context.Context, params models.ListApplicationsByTaxpayerRow) ([]models.Application, error)
	CreateApplicationAssessment(ctx context.Context, params models.CreateApplicationAssessmentParams) error

}

type repository struct {
	q *models.Queries
}

func NewRepository(db models.DBTX) Repository {
	return &repository{q: models.New(db)}
}

func (r *repository) CreateApplicationctx(ctx context.Context, application models.CreateApplicationParams) (models.Application, error) {
	return r.q.CreateApplication(ctx, application)
}

func (r *repository) GetApplicationByID(ctx context.Context, id uuid.UUID) (models.GetApplicationByIDRow, error) {
	return r.q.GetApplicationByID(ctx, id)
}

func (r *repository) UpdateApplicationStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.q.UpdateApplicationStatus(ctx, models.UpdateApplicationStatusParams{
		ID: id,
		Status: status,
	})
}

func (r *repository) CreateSingleBusinessPermit(ctx context.Context, params models.CreateSingleBusinessPermitParams) error {
	return r.q.CreateSingleBusinessPermit(ctx, params)
}

func (r *repository) CreateBuildingApproval(ctx context.Context, params models.CreateBuildingApprovalParams) error {
	return r.q.CreateBuildingApproval(ctx, params)
}

func (r *repository) CreateSeasonalParkingTicket(ctx context.Context, params models.CreateSeasonalParkingTicketParams) error {
	return r.q.CreateSeasonalParkingTicket(ctx, params)
}

func (r *repository) CreateHealthCertificate(ctx context.Context, params models.CreateHealthCertificateParams) error {
	return r.q.CreateHealthCertificate(ctx, params)
}

func (r *repository) CreateApplicationDocument(ctx context.Context, params models.CreateApplicationDocumentParams) error {
	return r.q.CreateApplicationDocument(ctx, params)
}

func (r *repository) ListApplicationsByTaxpayer(ctx context.Context, taxpayerID uuid.UUID) ([]models.ListApplicationsByTaxpayerRow, error) {
	return r.q.ListApplicationsByTaxpayer(ctx, taxpayerID)
}

func (r *repository) CreateApplicationAssessment(ctx context.Context, applicationID, assessmentID uuid.UUID) error {
	return r.q.CreateApplicationAssessment(ctx, models.CreateApplicationAssessmentParams{
		ApplicationID: applicationID,
		AssessmentID: assessmentID,
	})
}

// func (r *repository) DeleteApplication(ctx context.Context, id uuid.UUID) error {
// 	return r.q.DeleteApplication(ctx, id)
// }
