package user

import "time"

type User struct {
	ID           int64
	Name         string
	Username     string
	Email        string
	PasswordHash string
	Bio          string
	AvatarPath   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
