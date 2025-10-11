package assessment

import (
	"context"

	"github.com/google/uuid"
	"github.com/sangkips/revenue-system/internal/domain/assessment/models"
)


type Repository interface {
	CreateAssessment(ctx context.Context, assessment models.InsertAssessmentParams) (models.Assessment, error)
	GetAssessmentByID(ctx context.Context, id string) (models.Assessment, error)
	ListAssessments(ctx context.Context, params models.ListAssessmentsParams) ([]models.Assessment, error)
	UpdateAssessment(ctx context.Context, params models.UpdateAssessmentParams) (models.Assessment, error)
	DeleteAssessment(ctx context.Context, id string) error

	CreateAssessmentItem(ctx context.Context, item models.InsertAssessmentItemParams) (models.AssessmentItem, error)
	ListAssessmentItems(ctx context.Context, asessmentID string) ([]models.AssessmentItem, error)
	DeleteAssessmentItem(ctx context.Context, id string) error
	GetAssessmentItemByID(ctx context.Context, id string) (models.AssessmentItem, error)
}

type repository struct {
	q *models.Queries
}

func NewRepository(db models.DBTX) Repository {
	return &repository{q: models.New(db)}
}

func (r *repository) CreateAssessment(ctx context.Context, assessment models.InsertAssessmentParams) (models.Assessment, error) {
	return r.q.InsertAssessment(ctx, assessment)
}

func (r *repository) GetAssessmentByID(ctx context.Context, id string) (models.Assessment, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.Assessment{}, err
	}
	return r.q.GetAssessmentByID(ctx, parsedID)
}

func (r *repository) ListAssessments(ctx context.Context, params models.ListAssessmentsParams) ([]models.Assessment, error) {
	return r.q.ListAssessments(ctx, params)
}

func (r *repository) UpdateAssessment(ctx context.Context, params models.UpdateAssessmentParams) (models.Assessment, error) {
	return r.q.UpdateAssessment(ctx, params)
}

func (r *repository) DeleteAssessment(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.DeleteAssessment(ctx, parsedID)
}

func (r *repository) CreateAssessmentItem(ctx context.Context, item models.InsertAssessmentItemParams) (models.AssessmentItem, error) {
	return r.q.InsertAssessmentItem(ctx, item)
}

func (r *repository) ListAssessmentItems(ctx context.Context, assessmentID string) ([]models.AssessmentItem, error) {
	parsedID, err := uuid.Parse(assessmentID)
	if err != nil {
		return nil, err
	}
	return r.q.ListAssessmentItems(ctx, parsedID)
}

func (r *repository) GetAssessmentItemByID(ctx context.Context, id string) (models.AssessmentItem, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return models.AssessmentItem{}, err
	}
	return r.q.GetAssessmentItemByID(ctx, parsedID)
}


func (r *repository) DeleteAssessmentItem(ctx context.Context, id string) error {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.q.DeleteAssessmentItem(ctx, parsedID)
}