package message

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/auth"
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/platform/httpresponse"
)

type Handler struct {
	service *Service
}

type sendMessageRequest struct {
	Content string `json:"content"`
}

type messageResponseData struct {
	ID             int64     `json:"id"`
	ConversationID int64     `json:"conversation_id"`
	SenderID       int64     `json:"sender_id"`
	Content        string    `json:"content"`
	CreatedAt      time.Time `json:"created_at"`
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (handler *Handler) Send(
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

	conversationID, err := strconv.ParseInt(
		request.PathValue("conversationID"),
		10,
		64,
	)
	if err != nil || conversationID <= 0 {
		httpresponse.Error(
			responseWriter,
			http.StatusBadRequest,
			"INVALID_CONVERSATION_ID",
			"invalid conversation id",
			nil,
		)
		return
	}

	var requestBody sendMessageRequest

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

	sentMessage, err := handler.service.Send(
		request.Context(),
		currentUser.ID,
		conversationID,
		requestBody.Content,
	)

	if err != nil {
		if errors.Is(err, ErrEmptyMessage) {
			httpresponse.Error(
				responseWriter,
				http.StatusBadRequest,
				"EMPTY_MESSAGE",
				"message cannot be empty",
				nil,
			)
			return
		}

		if errors.Is(err, ErrMessageTooLong) {
			httpresponse.Error(
				responseWriter,
				http.StatusBadRequest,
				"MESSAGE_TOO_LONG",
				"message must not exceed 4000 characters",
				nil,
			)
			return
		}

		if errors.Is(err, ErrConversationDenied) {
			httpresponse.Error(
				responseWriter,
				http.StatusNotFound,
				"CONVERSATION_NOT_FOUND",
				"conversation not found",
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
		"message sent successfully",
		messageResponseData{
			ID:             sentMessage.ID,
			ConversationID: sentMessage.ConversationID,
			SenderID:       sentMessage.SenderID,
			Content:        sentMessage.Content,
			CreatedAt:      sentMessage.CreatedAt,
		},
	)
}

func (handler *Handler) History(
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

	conversationID, err := strconv.ParseInt(
		request.PathValue("conversationID"),
		10,
		64,
	)
	if err != nil || conversationID <= 0 {
		httpresponse.Error(
			responseWriter,
			http.StatusBadRequest,
			"INVALID_CONVERSATION_ID",
			"invalid conversation id",
			nil,
		)
		return
	}

	messages, err := handler.service.History(
		request.Context(),
		currentUser.ID,
		conversationID,
	)
	if err != nil {
		if errors.Is(err, ErrConversationDenied) {
			httpresponse.Error(
				responseWriter,
				http.StatusNotFound,
				"CONVERSATION_NOT_FOUND",
				"conversation not found",
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

	results := make([]messageResponseData, 0, len(messages))

	for _, foundMessage := range messages {
		results = append(results, messageResponseData{
			ID:             foundMessage.ID,
			ConversationID: foundMessage.ConversationID,
			SenderID:       foundMessage.SenderID,
			Content:        foundMessage.Content,
			CreatedAt:      foundMessage.CreatedAt,
		})
	}

	httpresponse.Success(
		responseWriter,
		http.StatusOK,
		"messages retrieved successfully",
		results,
	)
}
