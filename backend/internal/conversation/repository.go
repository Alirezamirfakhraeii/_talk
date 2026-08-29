package conversation

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
