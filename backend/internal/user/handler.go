package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/auth"
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/platform/httpresponse"
)

type Handler struct {
	service *Service
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponseData struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type currentUserResponseData struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type loginResponseData struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (handler *Handler) Register(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	var requestBody registerRequest

	err := json.NewDecoder(
		request.Body,
	).Decode(&requestBody)

	if err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusBadRequest,
			"INVALID_JSON",
			"invalid request body",
			nil,
		)

		return
	}

	registeredUser, err := handler.service.Register(
		request.Context(),
		RegisterInput{
			Name:     requestBody.Name,
			Username: requestBody.Username,
			Email:    requestBody.Email,
			Password: requestBody.Password,
		},
	)

	if err != nil {
		var validationErrors ValidationErrors

		if errors.As(
			err,
			&validationErrors,
		) {
			httpresponse.Error(
				responseWriter,
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				"invalid input",
				validationErrors,
			)

			return
		}

		if errors.Is(
			err,
			ErrEmailAlreadyExists,
		) {
			httpresponse.Error(
				responseWriter,
				http.StatusConflict,
				"EMAIL_ALREADY_EXISTS",
				"email already in use",
				nil,
			)

			return
		}

		if errors.Is(
			err,
			ErrUsernameAlreadyExists,
		) {
			httpresponse.Error(
				responseWriter,
				http.StatusConflict,
				"USERNAME_ALREADY_EXISTS",
				"username already in use",
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

	httpresponse.Success(
		responseWriter,
		http.StatusCreated,
		"user registered successfully",
		registerResponseData{
			ID:        registeredUser.ID,
			Name:      registeredUser.Name,
			Username:  registeredUser.Username,
			Email:     registeredUser.Email,
			CreatedAt: registeredUser.CreatedAt,
		},
	)
}

func (handler *Handler) Login(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	var requestBody loginRequest

	err := json.NewDecoder(
		request.Body,
	).Decode(&requestBody)

	if err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusBadRequest,
			"INVALID_JSON",
			"invalid request body",
			nil,
		)

		return
	}

	loginResult, err := handler.service.Login(
		request.Context(),
		LoginInput{
			Email:    requestBody.Email,
			Password: requestBody.Password,
		},
	)

	if err != nil {
		if errors.Is(
			err,
			ErrInvalidCredentials,
		) {
			httpresponse.Error(
				responseWriter,
				http.StatusUnauthorized,
				"INVALID_CREDENTIALS",
				"invalid email or password",
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

	http.SetCookie(
		responseWriter,
		&http.Cookie{
			Name:     "samatalk_session",
			Value:    loginResult.Token,
			Path:     "/",
			HttpOnly: true,
			Secure:   false,
			SameSite: http.SameSiteLaxMode,
			Expires:  loginResult.ExpiresAt,
		},
	)

	httpresponse.Success(
		responseWriter,
		http.StatusOK,
		"login successful",
		loginResponseData{
			ID:       loginResult.User.ID,
			Name:     loginResult.User.Name,
			Username: loginResult.User.Username,
			Email:    loginResult.User.Email,
		},
	)
}

func (handler *Handler) Me(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	cookie, err := request.Cookie(
		"samatalk_session",
	)

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

	currentUser, err := handler.service.CurrentUser(
		request.Context(),
		cookie.Value,
	)

	if err != nil {
		if errors.Is(
			err,
			auth.ErrInvalidSession,
		) {
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

	httpresponse.Success(
		responseWriter,
		http.StatusOK,
		"current user retrieved successfully",
		currentUserResponseData{
			ID:       currentUser.ID,
			Name:     currentUser.Name,
			Username: currentUser.Username,
			Email:    currentUser.Email,
		},
	)
}
