package domain

import (
	"errors"
	"net/mail"
	"time"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

func NewUser(email, passwordHash string) (*User, error) {
	if email == "" {
		return nil, errors.New("email boş olamaz")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("geçersiz email formatı")
	}

	if passwordHash == "" {
		return nil, errors.New("şifre hash'i boş olamaz")
	}

	return &User{
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}, nil

}
