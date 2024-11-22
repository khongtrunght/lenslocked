package models

import (
	"database/sql"
	"errors"
	"time"
)

const DefaultResetDuration = 1 * time.Hour

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when PasswordReset is created.
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB *sql.DB

	// BytesPerToken is the number of bytes that will be used to generate the token.
	BytesPerToken int
	// Duration is the amount of time that a PasswordReset will be valid for.
	// Defaults to DefaultResetDuration.
	Duration time.Duration
}

func (prs *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, errors.New("not implemented")
}

func (prs *PasswordResetService) Consume(token string) (*User, error) {
	return nil, errors.New("not implemented")
}
