package conversation

import (
	"context"
	"errors"
	"fmt"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/user"
)

var (
	ErrCannotChatWithSelf   = errors.New("cannot start conversation with yourself")
	ErrTargetUserNotFound   = errors.New("target user not found")
	ErrConversationNotFound = errors.New("conversation not found")
)

type Service struct {
	repository     *Repository
	userRepository *user.Repository
}

func NewService(
	repository *Repository,
	userRepository *user.Repository,
) *Service {
	return &Service{
		repository:     repository,
		userRepository: userRepository,
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

	_, err := service.userRepository.FindByID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return nil, ErrTargetUserNotFound
		}

		return nil, fmt.Errorf("find target user: %w", err)
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

func (service *Service) ListForUser(
	ctx context.Context,
	currentUserID int64,
) ([]ConversationSummary, error) {
	conversations, err := service.repository.ListForUser(
		ctx,
		currentUserID,
	)
	if err != nil {
		return nil, fmt.Errorf("list user conversations: %w", err)
	}

	return conversations, nil
}
