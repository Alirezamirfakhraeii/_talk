package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	database *pgxpool.Pool
}

func NewRepository(database *pgxpool.Pool) *Repository {
	return &Repository{database: database}
}

func (repository *Repository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (
			name,
			username,
			email,
			password_hash
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			created_at,
			updated_at
	`

	err := repository.database.QueryRow(
		ctx,
		query,
		user.Name,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		var pgError *pgconn.PgError

		if errors.As(err, &pgError) &&
			pgError.Code == "23505" {

			switch pgError.ConstraintName {

			case "users_email_key":
				return ErrEmailAlreadyExists

			case "users_username_key":
				return ErrUsernameAlreadyExists
			}
		}

		return fmt.Errorf(
			"create user: %w",
			err,
		)
	}

	return nil
}

func (repository *Repository) FindByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	query := `
		SELECT
			id,
			name,
			username,
			email,
			password_hash,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	foundUser := &User{}

	err := repository.database.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&foundUser.ID,
		&foundUser.Name,
		&foundUser.Username,
		&foundUser.Email,
		&foundUser.PasswordHash,
		&foundUser.CreatedAt,
		&foundUser.UpdatedAt,
	)

	if err != nil {
		if errors.Is(
			err,
			pgx.ErrNoRows,
		) {
			return nil, ErrUserNotFound
		}

		return nil, fmt.Errorf(
			"find user by email: %w",
			err,
		)
	}

	return foundUser, nil
}
