package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	database *pgxpool.Pool
}

type AuthenticatedUser struct {
	ID       int64
	Name     string
	Username string
	Email    string
}

func NewRepository(
	database *pgxpool.Pool,
) *Repository {
	return &Repository{
		database: database,
	}
}

func (repository *Repository) CreateSession(
	ctx context.Context,
	userID int64,
	tokenHash string,
	expiresAt time.Time,
) error {
	query := `
		INSERT INTO auth_sessions (
			user_id,
			token_hash,
			expires_at
		)
		VALUES ($1, $2, $3)
	`

	_, err := repository.database.Exec(
		ctx,
		query,
		userID,
		tokenHash,
		expiresAt,
	)

	if err != nil {
		return fmt.Errorf(
			"create auth session: %w",
			err,
		)
	}

	return nil
}

func (repository *Repository) FindUserByTokenHash(
	ctx context.Context,
	tokenHash string,
) (*AuthenticatedUser, error) {
	query := `
		SELECT
			users.id,
			users.name,
			users.username,
			users.email
		FROM auth_sessions
		INNER JOIN users
			ON users.id = auth_sessions.user_id
		WHERE auth_sessions.token_hash = $1
		  AND auth_sessions.expires_at > NOW()
	`

	authenticatedUser := &AuthenticatedUser{}

	err := repository.database.QueryRow(
		ctx,
		query,
		tokenHash,
	).Scan(
		&authenticatedUser.ID,
		&authenticatedUser.Name,
		&authenticatedUser.Username,
		&authenticatedUser.Email,
	)

	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return nil, ErrInvalidSession
		}

		return nil, fmt.Errorf(
			"find user by auth token: %w",
			err,
		)
	}

	return authenticatedUser, nil
}
