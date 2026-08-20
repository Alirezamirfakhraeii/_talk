package user

import "errors"

var (
	ErrEmailAlreadyExists = errors.New(
		"email already exists",
	)

	ErrUsernameAlreadyExists = errors.New(
		"username already exists",
	)

	ErrUserNotFound = errors.New(
		"user not found",
	)

	ErrInvalidCredentials = errors.New(
		"invalid credentials",
	)
)
