-- name: CreateApplication :one
INSERT INTO applications (
    id, taxpayer_id, type, notes, status
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, taxpayer_id, type, notes, status, submission_date, approval_date, created_at, updated_at;

-- name: CreateSingleBusinessPermit :exec
INSERT INTO single_business_permits (
    application_id, business_name, kra_pin, business_type, business_location, number_of_employees
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: CreateBuildingApproval :exec
INSERT INTO building_approvals (
    application_id, project_name, plot_parcel_number, project_type, estimated_project_cost, contact_email, contact_phone
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: CreateSeasonalParkingTicket :exec
INSERT INTO seasonal_parking_tickets (
    application_id, vehicle_registration_number, preferred_parking_zone, duration, contact_email, contact_phone
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: CreateHealthCertificate :exec
INSERT INTO health_certificates (
    application_id, applicant_name, business_name, contact_email, contact_phone
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: CreateApplicationDocument :exec
INSERT INTO application_documents (
    id, application_id, file_path, file_type, uploaded_at
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: GetApplicationByID :one
SELECT 
    a.id, a.taxpayer_id, a.type, a.notes, a.status, a.submission_date, a.approval_date, a.created_at, a.updated_at,
    sbp.business_name, sbp.kra_pin, sbp.business_type, sbp.business_location, sbp.number_of_employees,
    ba.project_name, ba.plot_parcel_number, ba.project_type, ba.estimated_project_cost, ba.contact_email, ba.contact_phone,
    spt.vehicle_registration_number, spt.preferred_parking_zone, spt.duration, spt.contact_email, spt.contact_phone,
    hc.applicant_name, hc.business_name, hc.contact_email, hc.contact_phone,
    t.email AS taxpayer_email, t.phone_number AS taxpayer_phone
FROM applications a
LEFT JOIN single_business_permits sbp ON a.id = sbp.application_id
LEFT JOIN building_approvals ba ON a.id = ba.application_id
LEFT JOIN seasonal_parking_tickets spt ON a.id = spt.application_id
LEFT JOIN health_certificates hc ON a.id = hc.application_id
JOIN taxpayers t ON a.taxpayer_id = t.id
WHERE a.id = $1;

-- name: ListApplicationsByTaxpayer :many
SELECT 
    a.id, a.taxpayer_id, a.type, a.notes, a.status, a.submission_date, a.approval_date, a.created_at, a.updated_at,
    sbp.business_name, sbp.kra_pin, sbp.business_type, sbp.business_location, sbp.number_of_employees,
    ba.project_name, ba.plot_parcel_number, ba.project_type, ba.estimated_project_cost, ba.contact_email, ba.contact_phone,
    spt.vehicle_registration_number, spt.preferred_parking_zone, spt.duration, spt.contact_email, spt.contact_phone,
    hc.applicant_name, hc.business_name, hc.contact_email, hc.contact_phone,
    t.email AS taxpayer_email, t.phone_number AS taxpayer_phone
FROM applications a
LEFT JOIN single_business_permits sbp ON a.id = sbp.application_id
LEFT JOIN building_approvals ba ON a.id = ba.application_id
LEFT JOIN seasonal_parking_tickets spt ON a.id = spt.application_id
LEFT JOIN health_certificates hc ON a.id = hc.application_id
JOIN taxpayers t ON a.taxpayer_id = t.id
WHERE a.taxpayer_id = $1
ORDER BY a.created_at DESC;

-- name: UpdateApplicationStatus :exec
UPDATE applications
SET status = $2, approval_date = CASE WHEN $2 = 'approved' THEN CURRENT_TIMESTAMP ELSE approval_date END, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CreateApplicationAssessment :exec
INSERT INTO application_assessments (
    application_id, assessment_id
) VALUES (
    $1, $2
);