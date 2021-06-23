package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/mail"
)

var (
	ErrorInvalidEmail    = errors.New("Invalid Email")
	ErrorInvalidPassword = errors.New("Invalid Password")
)

type User struct {
	ID       uuid.UUID
	Email    string
	PassHash string
	IsActive bool
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPassword(pass string) bool {
	const maxPasswordLen = 8
	return len(pass) > maxPasswordLen
}

func GenerateUser(email, pass string) (User, error) {
	var blankUser User

	if !isValidEmail(email) {
		return blankUser, ErrorInvalidEmail
	}

	if !isValidPassword(pass) {
		return blankUser, ErrorInvalidPassword
	}

	newPwHasher := sha256.New()
	newPwHasher.Write([]byte(pass))
	newPwHash := fmt.Sprintf("%x", newPwHasher.Sum(nil))

	return User{
		ID:       uuid.New(),
		Email:    email,
		PassHash: newPwHash,
		IsActive: false,
	}, nil
}
