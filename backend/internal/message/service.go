package message

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/conversation"
)

var (
	ErrEmptyMessage       = errors.New("message cannot be empty")
	ErrMessageTooLong     = errors.New("message is too long")
	ErrConversationDenied = errors.New("conversation not found or access denied")
)

type Service struct {
	repository             *Repository
	conversationRepository *conversation.Repository
}

type SendResult struct {
	Message         *Message
	RecipientUserID int64
}

func NewService(
	repository *Repository,
	conversationRepository *conversation.Repository,
) *Service {
	return &Service{
		repository:             repository,
		conversationRepository: conversationRepository,
	}
}

func (service *Service) Send(
	ctx context.Context,
	currentUserID int64,
	conversationID int64,
	content string,
) (*SendResult, error) {
	content = strings.TrimSpace(content)

	if content == "" {
		return nil, ErrEmptyMessage
	}

	if utf8.RuneCountInString(content) > 4000 {
		return nil, ErrMessageTooLong
	}

	foundConversation, err := service.conversationRepository.FindByIDForUser(
		ctx,
		conversationID,
		currentUserID,
	)
	if err != nil {
		if errors.Is(err, conversation.ErrConversationNotFound) {
			return nil, ErrConversationDenied
		}

		return nil, fmt.Errorf("check conversation access: %w", err)
	}

	recipientUserID := foundConversation.UserOneID

	if recipientUserID == currentUserID {
		recipientUserID = foundConversation.UserTwoID
	}

	newMessage := &Message{
		ConversationID: conversationID,
		SenderID:       currentUserID,
		Content:        content,
	}

	if err := service.repository.Create(ctx, newMessage); err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}

	return &SendResult{
		Message:         newMessage,
		RecipientUserID: recipientUserID,
	}, nil
}

func (service *Service) History(
	ctx context.Context,
	currentUserID int64,
	conversationID int64,
) ([]Message, error) {
	_, err := service.conversationRepository.FindByIDForUser(
		ctx,
		conversationID,
		currentUserID,
	)
	if err != nil {
		if errors.Is(err, conversation.ErrConversationNotFound) {
			return nil, ErrConversationDenied
		}

		return nil, fmt.Errorf("check conversation access: %w", err)
	}

	messages, err := service.repository.ListByConversation(
		ctx,
		conversationID,
		50,
	)
	if err != nil {
		return nil, fmt.Errorf("get conversation history: %w", err)
	}

	return messages, nil
}
