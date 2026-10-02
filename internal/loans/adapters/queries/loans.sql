-- name: GetLoanByID :one
SELECT * FROM loans WHERE uuid = $1;

-- name: SaveLoan :exec
INSERT INTO loans (uuid, borrower_uuid, officer_uuid, amount_idr, term_months, loan_type, interest_rate_apr, status, created_at, updated_at, assigned_to_uuid)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (uuid) DO UPDATE
SET assigned_to_uuid = EXCLUDED.assigned_to_uuid,
    updated_at = EXCLUDED.updated_at;
