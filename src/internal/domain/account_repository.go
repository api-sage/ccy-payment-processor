package domain

import "context"

type AccountRepository interface {
	Create(ctx context.Context, account Account) (Account, error)
}
