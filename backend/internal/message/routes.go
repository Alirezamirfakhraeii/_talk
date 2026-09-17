package message

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
		"GET /api/v1/conversations/{conversationID}/messages",
		authMiddleware.Authenticate(
			http.HandlerFunc(handler.History),
		),
	)

	mux.Handle(
		"POST /api/v1/conversations/{conversationID}/messages",
		authMiddleware.Authenticate(
			http.HandlerFunc(handler.Send),
		),
	)
}
