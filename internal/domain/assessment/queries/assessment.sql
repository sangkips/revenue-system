-- internal/domains/assessment/queries/assessment.sql
-- name: InsertAssessment :exec
INSERT INTO assessments (
    county_id, taxpayer_id, revenue_id, assessment_number, assessment_type,
    financial_year, base_amount, calculated_amount, total_amount,
    status, due_date, assessed_by, assessed_date
)
VALUES (
    @county_id, @taxpayer_id, @revenue_id, @assessment_number, @assessment_type,
    @financial_year, @base_amount, @calculated_amount, @total_amount,
    @status, @due_date, @assessed_by, @assessed_date
);

-- name: GetAssessmentByID :one
SELECT id, county_id, taxpayer_id, revenue_id, assessment_number, assessment_type,
       financial_year, base_amount, calculated_amount, total_amount, status, due_date,
       assessed_by, assessed_date, created_at, updated_at
FROM assessments
WHERE id = @id;

-- name: ListAssessments :many
SELECT id, county_id, taxpayer_id, revenue_id, assessment_number, assessment_type,
       financial_year, base_amount, calculated_amount, total_amount, status, due_date,
       assessed_by, assessed_date, created_at, updated_at
FROM assessments
WHERE county_id = @county_id
ORDER BY assessed_date DESC
LIMIT $1 OFFSET $2;

-- name: UpdateAssessment :exec
UPDATE assessments
SET 
    base_amount = COALESCE(@base_amount, base_amount),
    calculated_amount = COALESCE(@calculated_amount, calculated_amount),
    total_amount = COALESCE(@total_amount, total_amount),
    status = COALESCE(@status, status),
    due_date = COALESCE(@due_date, due_date),
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id;

-- name: DeleteAssessment :exec
DELETE FROM assessments WHERE id = @id;

-- Assessment Items Queries
-- name: InsertAssessmentItem :exec
INSERT INTO assessment_items (
    assessment_id, item_description, quantity, unit_amount, total_amount
)
VALUES (
    @assessment_id, @item_description, @quantity, @unit_amount, @total_amount
);

-- name: ListAssessmentItems :many
SELECT id, assessment_id, item_description, quantity, unit_amount, total_amount, created_at
FROM assessment_items
WHERE assessment_id = @assessment_id
ORDER BY created_at ASC;

-- name: DeleteAssessmentItem :exec
DELETE FROM assessment_items WHERE id = @id;