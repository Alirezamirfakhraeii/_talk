package conversation

import (
	"net/http"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/auth"
)

func RegisterRoutes(
	mux *http.ServeMux,
	handler *Handler,
	authMiddleware *auth.Middleware,
) {
	mux.Handle(
		"/api/v1/conversations",
		authMiddleware.Authenticate(
			http.HandlerFunc(handler.Start),
		),
	)
}
