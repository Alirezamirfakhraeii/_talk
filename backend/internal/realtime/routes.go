package realtime

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
		"GET /api/v1/ws",
		authMiddleware.Authenticate(
			http.HandlerFunc(handler.Connect),
		),
	)
}
