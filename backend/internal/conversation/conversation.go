package conversation

import "time"

type DirectConversation struct {
	ID        int64
	UserOneID int64
	UserTwoID int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ConversationSummary struct {
	ID              int64
	UserID          int64
	Name            string
	Username        string
	AvatarPath      string
	LastMessage     *string
	LastMessageTime *time.Time
}
