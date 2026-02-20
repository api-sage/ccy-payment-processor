package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/api-sage/ccy-payment-processor/src/internal/domain"
	"github.com/api-sage/ccy-payment-processor/src/internal/logger"
)

type TransferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

func (r *TransferRepository) Create(ctx context.Context, transfer domain.Transfer) (domain.Transfer, error) {
	logger.Info("transfer repository create", logger.Fields{
		"transactionReference": transfer.TransactionReference,
		"debitAccountNumber":   transfer.DebitAccountNumber,
		"creditAccountNumber":  transfer.CreditAccountNumber,
		"status":               transfer.Status,
	})

	const query = `
INSERT INTO transfers (
	external_refernece,
	transaction_reference,
	debit_account_number,
	credit_account_number,
	beneficiary_bank_code,
	debit_currency,
	credit_currency,
	debit_amount,
	credit_amount,
	fcy_rate,
	charge_amount,
	vat_amount,
	narration,
	status,
	audit_payload
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
RETURNING id, created_at, updated_at, processed_at`

	var (
		id          string
		createdAt   time.Time
		updatedAt   time.Time
		processedAt sql.NullTime
	)

	if err := r.db.QueryRowContext(
		ctx,
		query,
		transfer.ExternalRefernece,
		transfer.TransactionReference,
		transfer.DebitAccountNumber,
		transfer.CreditAccountNumber,
		transfer.BeneficiaryBankCode,
		transfer.DebitCurrency,
		transfer.CreditCurrency,
		transfer.DebitAmount,
		transfer.CreditAmount,
		transfer.FCYRate,
		transfer.ChargeAmount,
		transfer.VATAmount,
		transfer.Narration,
		transfer.Status,
		transfer.AuditPayload,
	).Scan(&id, &createdAt, &updatedAt, &processedAt); err != nil {
		logger.Error("transfer repository create failed", err, logger.Fields{
			"transactionReference": transfer.TransactionReference,
		})
		return domain.Transfer{}, fmt.Errorf("create transfer: %w", err)
	}

	transfer.ID = id
	transfer.CreatedAt = createdAt
	transfer.UpdatedAt = updatedAt
	if processedAt.Valid {
		value := processedAt.Time
		transfer.ProcessedAt = &value
	}

	logger.Info("transfer repository create success", logger.Fields{
		"transferId":           transfer.ID,
		"transactionReference": transfer.TransactionReference,
	})

	return transfer, nil
}

func (r *TransferRepository) Update(ctx context.Context, transfer domain.Transfer) (domain.Transfer, error) {
	return domain.Transfer{}, fmt.Errorf("not implemented")
}

func (r *TransferRepository) Get(ctx context.Context, id string, transactionReference string, externalRefernece string) (domain.Transfer, error) {
	return domain.Transfer{}, fmt.Errorf("not implemented")
}

func (r *TransferRepository) ProcessInternalTransfer(ctx context.Context, debitAccountNumber string, debitAmount string, suspenseAccountNumber string, creditAccountNumber string, creditAmount string) error {
	return fmt.Errorf("not implemented")
}

func (r *TransferRepository) UpdateStatus(ctx context.Context, transferID string, status domain.TransferStatus) error {
	return fmt.Errorf("not implemented")
}
