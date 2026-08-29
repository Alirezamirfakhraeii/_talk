package conversation

import (
	"context"
	"errors"
	"fmt"
)

var ErrCannotChatWithSelf = errors.New("cannot start conversation with yourself")

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (service *Service) Start(
	ctx context.Context,
	currentUserID int64,
	targetUserID int64,
) (*DirectConversation, error) {
	if currentUserID == targetUserID {
		return nil, ErrCannotChatWithSelf
	}

	userOneID := currentUserID
	userTwoID := targetUserID

	if userOneID > userTwoID {
		userOneID, userTwoID = userTwoID, userOneID
	}

	conversation, err := service.repository.FindOrCreate(
		ctx,
		userOneID,
		userTwoID,
	)
	if err != nil {
		return nil, fmt.Errorf("start direct conversation: %w", err)
	}

	return conversation, nil
}
