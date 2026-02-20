package domain

import "context"

type TransientAccountRepository interface {
	DebitSuspenseAccount(ctx context.Context, suspenseAccountNumber string, currency string, amount string) error
	CreditSuspenseAccount(ctx context.Context, suspenseAccountNumber string, currency string, amount string) error
}
