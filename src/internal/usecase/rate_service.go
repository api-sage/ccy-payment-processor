package usecase

import "github.com/api-sage/ccy-payment-processor/src/internal/domain"

type RateService struct {
	rateRepo domain.RateRepository
}

func NewRateService(rateRepo domain.RateRepository) *RateService {
	return &RateService{rateRepo: rateRepo}
}
