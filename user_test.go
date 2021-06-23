package main

import (
	"testing"
)

var testEmails = []struct {
	Mail   string
	Result bool
}{
	{"abcdef", false},
	{"abc@dateline.example.foo", true},
	{"", false},
}

func TestEmailValidation(t *testing.T) {
	for _, item := range testEmails {
		res := isValidEmail(item.Mail)
		if res != item.Result {
			t.Errorf("Validation error '%s' should be '%v', but got '%v'",
				item.Mail, item.Result, res)
		}
	}
}

var testPasswords = []struct {
	Password string
	Result   bool
}{
	{"", false},
	{"abced", false},
	{"password-valid", true},
}

func TestPasswordValidation(t *testing.T) {
	for _, item := range testPasswords {
		res := isValidPassword(item.Password)

		if res != item.Result {
			t.Errorf("Validation Error '%s' should be '%v', but got '%v'",
				item.Password, item.Result, res)
		}
	}
}

var blankUser User

var testUserInput = []struct {
	Email    string
	Password string
	Result   User
	Error    error
}{
	{"", "", blankUser, ErrorInvalidEmail},
	{"cilukba", "kekok", blankUser, ErrorInvalidEmail},
	{"cilukba@gmail.com", "inv", blankUser, ErrorInvalidPassword},
	{"cilukba@ggmail.gg", "zera4kbaureu", User{Email: "cilukba@ggmail.gg"}, nil},
}

func TestGenerateUser(t *testing.T) {
	for _, item := range testUserInput {
		u, err := GenerateUser(item.Email, item.Password)

		if err != nil && err != item.Error {
			t.Errorf("Error should be %v but got %v", item.Error.Error(), err.Error())
		}

		if err == nil && u.Email != item.Email {
			t.Errorf("Email doesn't match, should be %v, but got %v", item.Email, u.Email)
		}
	}
}
