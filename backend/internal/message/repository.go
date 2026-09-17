package message

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	database *pgxpool.Pool
}

func NewRepository(database *pgxpool.Pool) *Repository {
	return &Repository{
		database: database,
	}
}

func (repository *Repository) Create(
	ctx context.Context,
	message *Message,
) error {
	query := `
		INSERT INTO messages (
			conversation_id,
			sender_id,
			content
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			created_at
	`

	err := repository.database.QueryRow(
		ctx,
		query,
		message.ConversationID,
		message.SenderID,
		message.Content,
	).Scan(
		&message.ID,
		&message.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("create message: %w", err)
	}

	return nil
}

func (repository *Repository) ListByConversation(
	ctx context.Context,
	conversationID int64,
	limit int,
) ([]Message, error) {
	query := `
		SELECT
			id,
			conversation_id,
			sender_id,
			content,
			created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY id DESC
		LIMIT $2
	`

	rows, err := repository.database.Query(
		ctx,
		query,
		conversationID,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("list conversation messages: %w", err)
	}
	defer rows.Close()

	messages := make([]Message, 0)

	for rows.Next() {
		var foundMessage Message

		err := rows.Scan(
			&foundMessage.ID,
			&foundMessage.ConversationID,
			&foundMessage.SenderID,
			&foundMessage.Content,
			&foundMessage.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan conversation message: %w", err)
		}

		messages = append(messages, foundMessage)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate conversation messages: %w", err)
	}

	for left, right := 0, len(messages)-1; left < right; left, right = left+1, right-1 {
		messages[left], messages[right] = messages[right], messages[left]
	}

	return messages, nil
}
