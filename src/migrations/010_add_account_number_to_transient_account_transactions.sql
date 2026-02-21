ALTER TABLE transient_account_transactions
    ADD COLUMN IF NOT EXISTS account_number VARCHAR(32);
