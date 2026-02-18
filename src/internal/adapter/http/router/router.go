package router

import "net/http"

type AccountRouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler)
}

type UserRouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler)
}

type ParticipantBankRouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux, authMiddleware func(http.Handler) http.Handler)
}

func New(
	accountController AccountRouteRegistrar,
	userController UserRouteRegistrar,
	participantBankController ParticipantBankRouteRegistrar,
	authMiddleware func(http.Handler) http.Handler,
) *http.ServeMux {
	mux := http.NewServeMux()
	registerSwaggerRoutes(mux)
	mux.Handle("/account", http.RedirectHandler("/get-account", http.StatusMovedPermanently))
	mux.Handle("/verify-user-pin", http.RedirectHandler("/verify-pin", http.StatusMovedPermanently))
	mux.Handle("/getparticipantbanks", http.RedirectHandler("/get-participant-banks", http.StatusMovedPermanently))

	if accountController != nil {
		accountController.RegisterRoutes(mux, authMiddleware)
	}
	if userController != nil {
		userController.RegisterRoutes(mux, authMiddleware)
	}
	if participantBankController != nil {
		participantBankController.RegisterRoutes(mux, authMiddleware)
	}

	return mux
}
