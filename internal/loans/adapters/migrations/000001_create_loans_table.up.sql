-- 000001_create_loans_table.up.sql
CREATE TABLE loans (
    uuid              UUID PRIMARY KEY,
    borrower_uuid     UUID NOT NULL,
    amount_idr        BIGINT NOT NULL,
    term_months       INT NOT NULL,
    loan_type         TEXT NOT NULL,
    interest_rate_apr NUMERIC NOT NULL,
    status            TEXT NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL,
    updated_at        TIMESTAMPTZ NOT NULL
);
