package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/platform/httpresponse"
)

type Handler struct {
	service *Service
}

func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service: service,
	}
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
