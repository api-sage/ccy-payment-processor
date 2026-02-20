package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/api-sage/ccy-payment-processor/src/internal/domain"
	"github.com/api-sage/ccy-payment-processor/src/internal/logger"
)

type TransientAccountRepository struct {
	db *sql.DB
}

func NewTransientAccountRepository(db *sql.DB) *TransientAccountRepository {
	return &TransientAccountRepository{db: db}
}

func (r *TransientAccountRepository) DebitSuspenseAccount(ctx context.Context, suspenseAccountNumber string, currency string, amount string) error {
	logger.Info("transient account repository debit", logger.Fields{
		"accountNumber": suspenseAccountNumber,
		"currency":      currency,
		"amount":        amount,
	})

	const existsQuery = `
SELECT 1
FROM transient_accounts
WHERE account_number = $1
  AND UPPER(currency) = UPPER($2)`

	var exists int
	if err := r.db.QueryRowContext(ctx, existsQuery, suspenseAccountNumber, currency).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Info("transient account repository record not found", logger.Fields{
				"accountNumber": suspenseAccountNumber,
				"currency":      currency,
			})
			return domain.ErrRecordNotFound
		}
		logger.Error("transient account repository check failed", err, logger.Fields{
			"accountNumber": suspenseAccountNumber,
			"currency":      currency,
		})
		return fmt.Errorf("check transient account: %w", err)
	}

	const debitQuery = `
UPDATE transient_accounts
SET available_balance = available_balance - $2::numeric,
    updated_at = NOW()
WHERE account_number = $1
  AND available_balance >= $2::numeric`

	result, err := r.db.ExecContext(ctx, debitQuery, suspenseAccountNumber, amount)
	if err != nil {
		logger.Error("transient account repository debit failed", err, logger.Fields{
			"accountNumber": suspenseAccountNumber,
			"currency":      currency,
			"amount":        amount,
		})
		return fmt.Errorf("debit transient account: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("transient account repository debit rows affected failed", err, logger.Fields{
			"accountNumber": suspenseAccountNumber,
			"currency":      currency,
		})
		return fmt.Errorf("debit transient account rows affected: %w", err)
	}
	if rows == 0 {
		logger.Info("transient account repository insufficient balance", logger.Fields{
			"accountNumber": suspenseAccountNumber,
			"currency":      currency,
			"amount":        amount,
		})
		return domain.ErrInsufficientBalance
	}

	logger.Info("transient account repository debit success", logger.Fields{
		"accountNumber": suspenseAccountNumber,
		"currency":      currency,
		"amount":        amount,
	})
	return nil
}

func (r *TransientAccountRepository) CreditSuspenseAccount(ctx context.Context, suspenseAccountNumber string, currency string, amount string) error {
	logger.Info("transient account repository credit", logger.Fields{
		"accountNumber": suspenseAccountNumber,
		"currency":      currency,
		"amount":        amount,
	})

	const query = `
UPDATE transient_accounts
SET available_balance = available_balance + $2::numeric,
    updated_at = NOW()
WHERE account_number = $1
  AND UPPER(currency) = UPPER($3)`

	result, err := r.db.ExecContext(ctx, query, suspenseAccountNumber, amount, currency)
	if err != nil {
		logger.Error("transient account repository credit failed", err, logger.Fields{
			"accountNumber": suspenseAccountNumber,
			"currency":      currency,
			"amount":        amount,
		})
		return fmt.Errorf("credit transient account: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		logger.Error("transient account repository credit rows affected failed", err, logger.Fields{
			"accountNumber": suspenseAccountNumber,
			"currency":      currency,
		})
		return fmt.Errorf("credit transient account rows affected: %w", err)
	}
	if rows == 0 {
		logger.Info("transient account repository record not found", logger.Fields{
			"accountNumber": suspenseAccountNumber,
			"currency":      currency,
		})
		return domain.ErrRecordNotFound
	}

	logger.Info("transient account repository credit success", logger.Fields{
		"accountNumber": suspenseAccountNumber,
		"currency":      currency,
		"amount":        amount,
	})
	return nil
}
