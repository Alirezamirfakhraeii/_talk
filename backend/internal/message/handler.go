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

type RealtimeDeliverer interface {
	Deliver(userID int64, payload []byte)
}

type Handler struct {
	service  *Service
	realtime RealtimeDeliverer
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

func NewHandler(service *Service, realtime RealtimeDeliverer) *Handler {
	return &Handler{
		service:  service,
		realtime: realtime,
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
			"INVALID_REQUEST",
			"invalid request body",
			nil,
		)
		return
	}

	sendResult, err := handler.service.Send(
		request.Context(),
		currentUser.ID,
		conversationID,
		requestBody.Content,
	)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyMessage):
			httpresponse.Error(
				responseWriter,
				http.StatusUnprocessableEntity,
				"EMPTY_MESSAGE",
				"message content is required",
				nil,
			)

		case errors.Is(err, ErrMessageTooLong):
			httpresponse.Error(
				responseWriter,
				http.StatusUnprocessableEntity,
				"MESSAGE_TOO_LONG",
				"message content is too long",
				nil,
			)

		case errors.Is(err, ErrConversationDenied):
			httpresponse.Error(
				responseWriter,
				http.StatusNotFound,
				"CONVERSATION_NOT_FOUND",
				"conversation not found",
				nil,
			)

		default:
			httpresponse.Error(
				responseWriter,
				http.StatusInternalServerError,
				"INTERNAL_ERROR",
				"internal server error",
				nil,
			)
		}

		return
	}

	responseData := messageResponseData{
		ID:             sendResult.Message.ID,
		ConversationID: sendResult.Message.ConversationID,
		SenderID:       sendResult.Message.SenderID,
		Content:        sendResult.Message.Content,
		CreatedAt:      sendResult.Message.CreatedAt,
	}

	payload, err := json.Marshal(struct {
		Type string              `json:"type"`
		Data messageResponseData `json:"data"`
	}{
		Type: "message.created",
		Data: responseData,
	})

	if err == nil {
		handler.realtime.Deliver(
			sendResult.RecipientUserID,
			payload,
		)
	}

	httpresponse.Success(
		responseWriter,
		http.StatusCreated,
		"message sent successfully",
		responseData,
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

	responseData := make([]messageResponseData, 0, len(messages))

	for _, currentMessage := range messages {
		responseData = append(responseData, messageResponseData{
			ID:             currentMessage.ID,
			ConversationID: currentMessage.ConversationID,
			SenderID:       currentMessage.SenderID,
			Content:        currentMessage.Content,
			CreatedAt:      currentMessage.CreatedAt,
		})
	}

	httpresponse.Success(
		responseWriter,
		http.StatusOK,
		"messages retrieved successfully",
		responseData,
	)
}
