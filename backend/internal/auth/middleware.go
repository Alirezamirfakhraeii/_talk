package auth

import (
	"context"
	"errors"
	"net/http"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/platform/httpresponse"
)

type contextKey string

const authenticatedUserKey contextKey = "authenticated_user"

type Middleware struct {
	repository *Repository
}

func NewMiddleware(repository *Repository) *Middleware {
	return &Middleware{
		repository: repository,
	}
}

func (middleware *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(
		responseWriter http.ResponseWriter,
		request *http.Request,
	) {
		cookie, err := request.Cookie("samatalk_session")
		if err != nil {
			httpresponse.Error(
				responseWriter,
				http.StatusUnauthorized,
				"UNAUTHORIZED",
				"authentication required",
				nil,
			)
			return
		}

		tokenHash := HashToken(cookie.Value)

		user, err := middleware.repository.FindUserByTokenHash(
			request.Context(),
			tokenHash,
		)

		if err != nil {
			if errors.Is(err, ErrInvalidSession) {
				httpresponse.Error(
					responseWriter,
					http.StatusUnauthorized,
					"UNAUTHORIZED",
					"invalid or expired session",
					nil,
				)
				return
			}

			httpresponse.Error(
				responseWriter,
				http.StatusInternalServerError,
				"INTERNAL_ERROR",
				"internal server error",
				nil,
			)
			return
		}

		ctx := context.WithValue(
			request.Context(),
			authenticatedUserKey,
			user,
		)

		next.ServeHTTP(
			responseWriter,
			request.WithContext(ctx),
		)
	})
}

func CurrentUserFromContext(ctx context.Context) (*AuthenticatedUser, bool) {
	user, ok := ctx.Value(authenticatedUserKey).(*AuthenticatedUser)

	return user, ok
}
