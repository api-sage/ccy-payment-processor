package controller

import (
	"context"
	"net/http"

	"github.com/api-sage/ccy-payment-processor/src/internal/adapter/http/models"
)

type ChargesService interface {
	GetChargesSummary(ctx context.Context, req models.GetChargesRequest) (models.Response[models.GetChargesResponse], error)
}

type ChargesController struct {
	service ChargesService
}

func NewChargesController(service ChargesService) *ChargesController {
	return &ChargesController{service: service}
}

func (c *ChargesController) RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler) {
	handler := http.HandlerFunc(c.getCharges)
	if authMiddleware != nil {
		handler = authMiddleware(handler).ServeHTTP
	}

	mux.Handle("/get-charges", http.HandlerFunc(handler))
}

func (c *ChargesController) getCharges(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, models.ErrorResponse[models.GetChargesResponse]("method not allowed"))
		return
	}

	req := models.GetChargesRequest{
		Amount:       r.URL.Query().Get("amount"),
		FromCurrency: r.URL.Query().Get("fromCurrency"),
	}

	response, err := c.service.GetChargesSummary(r.Context(), req)
	if err != nil {
		status := http.StatusInternalServerError
		if response.Message == "validation failed" {
			status = http.StatusBadRequest
		}
		writeJSON(w, status, response)
		return
	}

	writeJSON(w, http.StatusOK, response)
}
