package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

type RegisterInput struct {
	Name     string
	Username string
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
