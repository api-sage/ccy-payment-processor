package domain

import "context"

type TransientAccountRepository interface {
	DebitTransientAccount(ctx context.Context, transientAccountNumber string, amount string) error
	CreditTransientAccount(ctx context.Context, transientAccountNumber string, amount string) error
}
