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
	mux.HandleFunc(
		"POST /api/v1/auth/register",
		handler.Register,
	)

	mux.HandleFunc(
		"POST /api/v1/auth/login",
		handler.Login,
	)

	mux.Handle(
		"GET /api/v1/me",
		authMiddleware.Authenticate(
			http.HandlerFunc(handler.Me),
		),
	)

	mux.Handle(
		"PUT /api/v1/me/profile",
		authMiddleware.Authenticate(
			http.HandlerFunc(handler.UpdateProfile),
		),
	)

	mux.Handle(
		"GET /api/v1/users/search",
		authMiddleware.Authenticate(
			http.HandlerFunc(handler.Search),
		),
	)

	mux.Handle(
		"POST /api/v1/me/avatar",
		authMiddleware.Authenticate(
			http.HandlerFunc(handler.UploadAvatar),
		),
	)
}
