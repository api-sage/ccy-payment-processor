package usecase

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/api-sage/ccy-payment-processor/src/internal/adapter/http/models"
	"github.com/api-sage/ccy-payment-processor/src/internal/domain"
)

type AccountService struct {
	accountRepo domain.AccountRepository
}

func NewAccountService(accountRepo domain.AccountRepository) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (s *AccountService) CreateAccount(ctx context.Context, req models.CreateAccountRequest) (domain.Account, error) {
	if err := req.Validate(); err != nil {
		return domain.Account{}, err
	}

	balance, err := parseBalance(req.InitialDeposit)
	if err != nil {
		return domain.Account{}, err
	}

	account := domain.Account{
		CustomerID:       strings.TrimSpace(req.CustomerID),
		AccountNumber:    generateAccountNumber(),
		Currency:         strings.ToUpper(strings.TrimSpace(req.Currency)),
		AvailableBalance: balance,
		LedgerBalance:    balance,
		Status:           domain.AccountStatusActive,
	}

	created, err := s.accountRepo.Create(ctx, account)
	if err != nil {
		return domain.Account{}, err
	}

	return created, nil
}

func parseBalance(raw string) (string, error) {
	if strings.TrimSpace(raw) == "" {
		return "0.00", nil
	}

	parsed, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return "", fmt.Errorf("initialDeposit must be a valid number: %w", err)
	}

	if parsed < 0 {
		return "", fmt.Errorf("initialDeposit cannot be negative")
	}

	return fmt.Sprintf("%.2f", parsed), nil
}

func generateAccountNumber() string {
	return fmt.Sprintf("%010d", time.Now().UnixNano()%10_000_000_000)
}
