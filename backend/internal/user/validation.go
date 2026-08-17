package user

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var usernamePattern = regexp.MustCompile(`^[a-z0-9_]+$`)

var emailPattern = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

type ValidationErrors map[string]string

func (validationErrors ValidationErrors) Error() string {
	return "validation failed"
}

func ValidateRegisterInput(input *RegisterInput) error {
	errors := ValidationErrors{}

	input.Name = strings.TrimSpace(input.Name)
	input.Username = strings.ToLower(strings.TrimSpace(input.Username))

	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	// Name
	nameLength := utf8.RuneCountInString(input.Name)

	if nameLength < 2 {
		errors["name"] = "name must be at least 2 characters"
	} else if nameLength > 100 {
		errors["name"] = "name must be at most 100 characters"
	}

	// Username
	usernameLength := len(input.Username)

	if usernameLength < 3 {
		errors["username"] = "username must be at least 3 characters"
	} else if usernameLength > 50 {
		errors["username"] = "username must be at most 50 characters"
	} else if !usernamePattern.MatchString(input.Username) {
		errors["username"] = "username may contain only lowercase letters, numbers, and underscore"
	}

	// Email
	if input.Email == "" {
		errors["email"] = "email is required"
	} else if len(input.Email) > 255 {
		errors["email"] = "email must be at most 255 characters"
	} else if !emailPattern.MatchString(input.Email) {
		errors["email"] = "email format is invalid"
	}

	// Password
	passwordLength := len([]byte(input.Password))

	if passwordLength < 8 {
		errors["password"] = "password must be at least 8 characters"
	} else if passwordLength > 72 {
		errors["password"] = "password must be at most 72 bytes"
	}

	if len(errors) > 0 {
		return errors
	}

	return nil
}
