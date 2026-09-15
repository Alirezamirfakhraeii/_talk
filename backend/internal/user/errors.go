package user

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")

	ErrProfileNameRequired     = errors.New("profile name is required")
	ErrProfileUsernameRequired = errors.New("profile username is required")
	ErrProfileBioTooLong       = errors.New("profile bio is too long")
)
