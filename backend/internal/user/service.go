package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository        *Repository
	sessionRepository *auth.Repository
}

func NewService(
	repository *Repository,
	sessionRepository *auth.Repository,
) *Service {
	return &Service{
		repository:        repository,
		sessionRepository: sessionRepository,
	}
}

type LoginResult struct {
	User      *User
	Token     string
	ExpiresAt time.Time
}

type RegisterInput struct {
	Name     string
	Username string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

func (service *Service) Register(
	ctx context.Context,
	input RegisterInput,
) (*User, error) {
	if err := ValidateRegisterInput(&input); err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"hash password: %w",
			err,
		)
	}

	newUser := &User{
		Name:         input.Name,
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: string(passwordHash),
	}

	err = service.repository.Create(
		ctx,
		newUser,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"register user: %w",
			err,
		)
	}

	return newUser, nil
}

func (service *Service) Login(
	ctx context.Context,
	input LoginInput,
) (*LoginResult, error) {
	email := strings.ToLower(
		strings.TrimSpace(input.Email),
	)

	foundUser, err := service.repository.FindByEmail(
		ctx,
		email,
	)

	if err != nil {
		if errors.Is(
			err,
			ErrUserNotFound,
		) {
			return nil, ErrInvalidCredentials
		}

		return nil, fmt.Errorf(
			"login user: %w",
			err,
		)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(foundUser.PasswordHash),
		[]byte(input.Password),
	)

	if err != nil {
		return nil, ErrInvalidCredentials
	}

	token, tokenHash, err := auth.GenerateToken()
	if err != nil {
		return nil, fmt.Errorf(
			"generate login token: %w",
			err,
		)
	}

	expiresAt := time.Now().
		UTC().
		Add(7 * 24 * time.Hour)

	err = service.sessionRepository.CreateSession(
		ctx,
		foundUser.ID,
		tokenHash,
		expiresAt,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"create login session: %w",
			err,
		)
	}

	return &LoginResult{
		User:      foundUser,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (service *Service) CurrentUser(
	ctx context.Context,
	token string,
) (*auth.AuthenticatedUser, error) {
	if token == "" {
		return nil, auth.ErrInvalidSession
	}

	tokenHash := auth.HashToken(token)

	authenticatedUser, err :=
		service.sessionRepository.FindUserByTokenHash(
			ctx,
			tokenHash,
		)

	if err != nil {
		if errors.Is(
			err,
			auth.ErrInvalidSession,
		) {
			return nil, auth.ErrInvalidSession
		}

		return nil, fmt.Errorf(
			"get current user: %w",
			err,
		)
	}

	return authenticatedUser, nil
}
