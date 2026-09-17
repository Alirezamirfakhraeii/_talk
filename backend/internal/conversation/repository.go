package conversation

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func (repository *Repository) FindOrCreate(
	ctx context.Context,
	userOneID int64,
	userTwoID int64,
) (*DirectConversation, error) {
	query := `
		INSERT INTO direct_conversations (
			user_one_id,
			user_two_id
		)
		VALUES ($1, $2)
		ON CONFLICT (user_one_id, user_two_id)
		DO UPDATE SET
			updated_at = direct_conversations.updated_at
		RETURNING
			id,
			user_one_id,
			user_two_id,
			created_at,
			updated_at
	`

	foundConversation := &DirectConversation{}

	err := repository.database.QueryRow(
		ctx,
		query,
		userOneID,
		userTwoID,
	).Scan(
		&foundConversation.ID,
		&foundConversation.UserOneID,
		&foundConversation.UserTwoID,
		&foundConversation.CreatedAt,
		&foundConversation.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"find or create direct conversation: %w",
			err,
		)
	}

	return foundConversation, nil
}

func (repository *Repository) FindByIDForUser(
	ctx context.Context,
	conversationID int64,
	userID int64,
) (*DirectConversation, error) {
	query := `
		SELECT
			id,
			user_one_id,
			user_two_id,
			created_at,
			updated_at
		FROM direct_conversations
		WHERE id = $1
		  AND (user_one_id = $2 OR user_two_id = $2)
	`

	conversation := &DirectConversation{}

	err := repository.database.QueryRow(
		ctx,
		query,
		conversationID,
		userID,
	).Scan(
		&conversation.ID,
		&conversation.UserOneID,
		&conversation.UserTwoID,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrConversationNotFound
		}

		return nil, fmt.Errorf("find conversation for user: %w", err)
	}

	return conversation, nil
}

func (repository *Repository) ListForUser(
	ctx context.Context,
	userID int64,
) ([]ConversationSummary, error) {
	query := `
		SELECT
			conversation.id,
			other_user.id,
			other_user.name,
			other_user.username,
			COALESCE(other_user.avatar_path, ''),
			last_message.content,
			last_message.created_at
		FROM direct_conversations AS conversation
		JOIN users AS other_user
			ON other_user.id = CASE
				WHEN conversation.user_one_id = $1
					THEN conversation.user_two_id
				ELSE conversation.user_one_id
			END
		LEFT JOIN LATERAL (
			SELECT
				message.content,
				message.created_at
			FROM messages AS message
			WHERE message.conversation_id = conversation.id
			ORDER BY message.id DESC
			LIMIT 1
		) AS last_message ON TRUE
		WHERE
			conversation.user_one_id = $1
			OR conversation.user_two_id = $1
		ORDER BY
			COALESCE(
				last_message.created_at,
				conversation.updated_at
			) DESC
	`

	rows, err := repository.database.Query(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list conversations: %w",
			err,
		)
	}
	defer rows.Close()

	conversations := make(
		[]ConversationSummary,
		0,
	)

	for rows.Next() {
		var summary ConversationSummary
		var lastMessage pgtype.Text
		var lastMessageTime pgtype.Timestamptz

		err := rows.Scan(
			&summary.ID,
			&summary.UserID,
			&summary.Name,
			&summary.Username,
			&summary.AvatarPath,
			&lastMessage,
			&lastMessageTime,
		)
		if err != nil {
			return nil, fmt.Errorf(
				"scan conversation summary: %w",
				err,
			)
		}

		if lastMessage.Valid {
			summary.LastMessage =
				&lastMessage.String
		}

		if lastMessageTime.Valid {
			summary.LastMessageTime =
				&lastMessageTime.Time
		}

		conversations = append(
			conversations,
			summary,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate conversations: %w",
			err,
		)
	}

	return conversations, nil
}
