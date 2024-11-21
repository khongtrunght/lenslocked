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

	// TODO: save to the database
	return &session, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	// TODO: Implement the User method
	return nil, nil
}

func (ss *SessionService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}
