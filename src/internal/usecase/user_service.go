package usecase

import "github.com/api-sage/ccy-payment-processor/src/internal/domain"

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}
