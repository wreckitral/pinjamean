ALTER TABLE loans
ADD COLUMN assigned_to_uuid UUID;

UPDATE loans
SET assigned_to_uuid = officer_uuid;


ALTER TABLE loans
ALTER COLUMN assigned_to_uuid SET NOT NULL;
