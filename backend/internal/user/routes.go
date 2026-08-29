package user

import (
	"net/http"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/auth"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	authMiddleware *auth.Middleware,
) {
	mux.HandleFunc("/api/v1/auth/register", handler.Register)
	mux.HandleFunc("/api/v1/auth/login", handler.Login)

	mux.Handle(
		"/api/v1/me",
		authMiddleware.Authenticate(
			http.HandlerFunc(handler.Me),
		),
	)

	mux.Handle(
		"/api/v1/users/search",
		authMiddleware.Authenticate(
			http.HandlerFunc(handler.Search),
		),
	)
}
