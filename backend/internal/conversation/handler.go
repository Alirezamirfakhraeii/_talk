package conversation

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

type startConversationRequest struct {
	UserID int64 `json:"user_id"`
}

type conversationResponseData struct {
	ID        int64 `json:"id"`
	UserOneID int64 `json:"user_one_id"`
	UserTwoID int64 `json:"user_two_id"`
}

type conversationSummaryResponseData struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"user_id"`
	Name            string     `json:"name"`
	Username        string     `json:"username"`
	AvatarPath      string     `json:"avatar_path"`
	LastMessage     *string    `json:"last_message"`
	LastMessageTime *time.Time `json:"last_message_time"`
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

	if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
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

		if errors.Is(err, ErrTargetUserNotFound) {
			httpresponse.Error(
				responseWriter,
				http.StatusNotFound,
				"USER_NOT_FOUND",
				"user not found",
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

func (handler *Handler) List(
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

	conversations, err := handler.service.ListForUser(
		request.Context(),
		currentUser.ID,
	)
	if err != nil {
		httpresponse.Error(
			responseWriter,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
			nil,
		)
		return
	}

	results := make([]conversationSummaryResponseData, 0, len(conversations))

	for _, conversation := range conversations {
		results = append(results, conversationSummaryResponseData{
			ID:              conversation.ID,
			UserID:          conversation.UserID,
			Name:            conversation.Name,
			Username:        conversation.Username,
			AvatarPath:      conversation.AvatarPath,
			LastMessage:     conversation.LastMessage,
			LastMessageTime: conversation.LastMessageTime,
		})
	}

	httpresponse.Success(
		responseWriter,
		http.StatusOK,
		"conversations retrieved successfully",
		results,
	)
}
