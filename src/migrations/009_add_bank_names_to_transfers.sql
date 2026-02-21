ALTER TABLE transfers
    ADD COLUMN IF NOT EXISTS debit_bank_name VARCHAR(255),
    ADD COLUMN IF NOT EXISTS credit_bank_name VARCHAR(255);
