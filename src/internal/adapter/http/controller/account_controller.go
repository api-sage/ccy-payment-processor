package controller

import (
	"context"
	"net/http"

	"github.com/api-sage/ccy-payment-processor/src/internal/adapter/http/models"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req models.CreateAccountRequest) (models.Response[models.CreateAccountResponse], error)
}

type AccountController struct {
	service AccountService
}

func NewAccountController(service AccountService) *AccountController {
	return &AccountController{service: service}
}

func (c *AccountController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/accounts", c.createAccount)
}

func (c *AccountController) createAccount(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}
