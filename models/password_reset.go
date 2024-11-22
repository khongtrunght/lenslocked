package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/khongtrunght/lenslocked/rand"
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

	Now func() time.Time
}

func (prs *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.ToLower(email)
	row := prs.DB.QueryRow(`
    SELECT id FROM users WHERE email = $1
  `, email)
	var userID int
	err := row.Scan(&userID)
	if err != nil {
		// consider not found an error
		return nil, fmt.Errorf("create: %w", err)
	}

	bytesPerToken := prs.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}

	token, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	duration := prs.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}

	reset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: prs.hash(token),
		ExpiresAt: prs.Now().Add(duration),
	}

	row = prs.DB.QueryRow(`
    INSERT INTO password_resets (user_id, token_hash, expires_at)
    VALUES ($1, $2, $3)
    ON CONFLICT (user_id) DO UPDATE
    SET token_hash = $2, expires_at = $3
    RETURNING id;
  `, reset.UserID, reset.TokenHash, reset.ExpiresAt)

	err = row.Scan(&reset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}

	return &reset, nil
}

func (prs *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := prs.hash(token)
	var user User
	var pwReset PasswordReset
	row := prs.DB.QueryRow(`
    SELECT password_resets.id,  password_resets.expires_at,
    users.id, users.email, users.password_hash
    FROM password_resets
    JOIN users
    ON password_resets.user_id = users.id
    WHERE password_resets.token_hash = $1;
  `, tokenHash)

	err := row.Scan(&pwReset.ID, &pwReset.ExpiresAt, &user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}

	if prs.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired")
	}

	err = prs.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	return &user, nil
}

func (prs *PasswordResetService) delete(id int) error {
	_, err := prs.DB.Exec(`
    DELETE FROM password_resets WHERE id = $1;
  `, id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (prs *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
