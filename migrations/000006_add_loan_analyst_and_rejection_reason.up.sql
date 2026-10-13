ALTER TABLE loans
    ADD COLUMN analyst_uuid UUID,
    ADD COLUMN rejection_reason TEXT;
