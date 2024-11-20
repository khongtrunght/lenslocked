package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) {
	// TODO: implement this
	email = strings.ToLower(email)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	passwordHashed := string(hashedBytes)
	user := User{
		Email:        email,
		PasswordHash: passwordHashed,
	}
	row := us.DB.QueryRow(
		`
    INSERT INTO users (email, password_hash)
    VALUES ($1,$2) RETURNING id;
    `, email, passwordHashed)
	err = row.Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &user, nil
}

func (us *UserService) Update(user *User) error {
	return nil
}