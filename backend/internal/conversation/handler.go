package conversation

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/auth"
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/platform/httpresponse"
)

type Handler struct {
	service *Service
}

type startConversationRequest struct {
	UserID int64 `json:"user_id"`
}

type conversationResponseData struct {
	ID        int64 `json:"id"`
	UserOneID int64 `json:"user_one_id"`
	UserTwoID int64 `json:"user_two_id"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (handler *Handler) Start(
	responseWriter http.ResponseWriter,
	request *http.Request,
) {
	currentUser, ok := auth.CurrentUserFromContext(request.Context())
	if !ok {
		httpresponse.Error(
			responseWriter,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authentication required",
			nil,
		)
		return
	}

	var requestBody startConversationRequest

	err := json.NewDecoder(request.Body).Decode(&requestBody)
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

	if requestBody.UserID <= 0 {
		httpresponse.Error(
			responseWriter,
			http.StatusBadRequest,
			"INVALID_USER_ID",
			"invalid user id",
			nil,
		)
		return
	}

	conversation, err := handler.service.Start(
		request.Context(),
		currentUser.ID,
		requestBody.UserID,
	)

	if err != nil {
		if errors.Is(err, ErrCannotChatWithSelf) {
			httpresponse.Error(
				responseWriter,
				http.StatusBadRequest,
				"CANNOT_CHAT_WITH_SELF",
				"you cannot start a conversation with yourself",
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
		"conversation ready",
		conversationResponseData{
			ID:        conversation.ID,
			UserOneID: conversation.UserOneID,
			UserTwoID: conversation.UserTwoID,
		},
	)
}
