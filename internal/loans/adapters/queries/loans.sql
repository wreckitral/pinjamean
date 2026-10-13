-- name: GetLoanByID :one
SELECT * FROM loans WHERE uuid = $1;

-- name: SaveLoan :exec
INSERT INTO loans (uuid, borrower_uuid, officer_uuid, amount_idr, term_months, loan_type, interest_rate_apr, status, created_at, updated_at, assigned_to_uuid, analyst_uuid, rejection_reason)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
ON CONFLICT (uuid) DO UPDATE
SET assigned_to_uuid = EXCLUDED.assigned_to_uuid,
    analyst_uuid = EXCLUDED.analyst_uuid,
    rejection_reason = EXCLUDED.rejection_reason,
    status = EXCLUDED.status,
    updated_at = EXCLUDED.updated_at;
