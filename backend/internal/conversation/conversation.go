package conversation

import "time"

type DirectConversation struct {
	ID        int64
	UserOneID int64
	UserTwoID int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
