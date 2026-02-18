package router

import "net/http"

type AccountRouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler)
}

func New(accountController AccountRouteRegistrar, authMiddleware func(http.Handler) http.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	if accountController != nil {
		accountController.RegisterRoutes(mux, authMiddleware)
	}

	return mux
}
