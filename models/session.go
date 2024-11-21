package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"

	"github.com/khongtrunght/lenslocked/rand"
)

const MinBytesPerToken = 32

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session. When look up a session
	// this will be left empty, as we only store the hashed version.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB

	BytesPerToken int
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	// TODO: Create the session token
	// TODO: Implement the Create method
	bytePerToken := ss.BytesPerToken
	if bytePerToken < MinBytesPerToken {
		bytePerToken = MinBytesPerToken
	}

	token, err := rand.String(bytePerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	session := Session{
		UserID:    userID,
		Token:     token,
		TokenHash: ss.hash(token),
	}

	row := ss.DB.QueryRow(`
      UPDATE sessions
      SET token_hash = $1
      WHERE user_id = $2
      RETURNING id;
    `, session.TokenHash, session.UserID)
	err = row.Scan(&session.ID)
	if err == sql.ErrNoRows {
		row = ss.DB.QueryRow(`
    INSERT INTO sessions (user_id, token_hash)
    VALUES ($1, $2) RETURNING id;
  `, session.UserID, session.TokenHash)
		err = row.Scan(&session.ID)
	}
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	tokenHash := ss.hash(token)
	var user User
	row := ss.DB.QueryRow(`
    SELECT users.id, users.email, users.password_hash
    FROM users
    JOIN sessions
    ON users.id = sessions.user_id
    WHERE sessions.token_hash = $1;
  `, tokenHash)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("user: %w", err)
	}
	return &user, nil
}

func (ss *SessionService) Delete(token string) error {
	tokenHash := ss.hash(token)
	_, err := ss.DB.Exec(`
    DELETE FROM sessions
    WHERE token_hash = $1;
  `, tokenHash)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
