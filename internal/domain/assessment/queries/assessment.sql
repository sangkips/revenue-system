-- internal/domains/assessment/queries/assessment.sql
-- name: InsertAssessment :one
INSERT INTO assessments (
    county_id, taxpayer_id, revenue_id, assessment_number, assessment_type,
    financial_year, base_amount, calculated_amount, total_amount,
    status, due_date, assessed_by, assessed_date
)
VALUES (
    @county_id, @taxpayer_id, @revenue_id, @assessment_number, @assessment_type,
    @financial_year, @base_amount, @calculated_amount, @total_amount,
    @status, @due_date, @assessed_by, @assessed_date
)
RETURNING id, county_id, taxpayer_id, revenue_id, assessment_number, assessment_type,
    financial_year, base_amount, calculated_amount, total_amount, status, due_date,
    assessed_by, assessed_date, created_at, updated_at;

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

-- name: UpdateAssessment :one
UPDATE assessments
SET
    base_amount = CASE WHEN @base_amount::text = '' THEN base_amount ELSE @base_amount::decimal END,
    calculated_amount = CASE WHEN @calculated_amount::text = '' THEN calculated_amount ELSE @calculated_amount::decimal END,
    total_amount = CASE WHEN @total_amount::text = '' THEN total_amount ELSE @total_amount::decimal END,
    status = CASE WHEN @status = '' THEN status ELSE @status END,
    due_date = CASE WHEN @due_date = '1970-01-01T00:00:00Z'::timestamptz THEN due_date ELSE @due_date END,
    updated_at = CURRENT_TIMESTAMP
WHERE id = @id
RETURNING id, county_id, taxpayer_id, revenue_id, assessment_number, assessment_type,
    financial_year, base_amount, calculated_amount, total_amount, status, due_date,
    assessed_by, assessed_date, created_at, updated_at;

-- name: DeleteAssessment :exec
DELETE FROM assessments WHERE id = @id;

-- Assessment Items Queries
-- name: InsertAssessmentItem :one
INSERT INTO assessment_items (
    assessment_id, item_description, quantity, unit_amount, total_amount
)
VALUES (
    @assessment_id, @item_description, @quantity, @unit_amount, @total_amount
)
RETURNING id, assessment_id, item_description, quantity, unit_amount, total_amount, created_at;

-- name: ListAssessmentItems :many
SELECT id, assessment_id, item_description, quantity, unit_amount, total_amount, created_at
FROM assessment_items
WHERE assessment_id = @assessment_id
ORDER BY created_at ASC;

-- name: GetAssessmentItemByID :one
SELECT id, assessment_id, item_description, quantity, unit_amount, total_amount, created_at
FROM assessment_items
WHERE id = @id;

-- name: DeleteAssessmentItem :exec
DELETE FROM assessment_items WHERE id = @id;
