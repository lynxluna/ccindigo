package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const dbstr = "user=postgres password=postgres dbname=indigo sslmode=disable"

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(http.StatusBadRequest)
	errJSON := JSONError{Message: "Invalid Email Address"}
	payload, _ := json.Marshal(errJSON)

	w.Write(payload)
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isValidPassword(pass string) bool {
	const maxPasswordLen = 8
	return len(pass) > maxPasswordLen
}

func saveUser(email, pass string) (*uuid.UUID, error) {

	db, err := sql.Open("postgres", dbstr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	newUserID := uuid.New()
	newPwHasher := sha256.New()
	newPwHasher.Write([]byte(pass))
	newPwHash := fmt.Sprintf("%x", newPwHasher.Sum(nil))

	regCommand := "INSERT INTO users (id, email, passhash) VALUES ($1, $2, $3)"

	if _, err := db.Exec(regCommand, newUserID, email, newPwHash); err != nil {

		if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
			return nil, errors.New("User already exists")
		}

		return nil, err
	}

	return &newUserID, nil
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	pass := r.PostFormValue("password")

	if !isValidEmail(email) {
		writeJSONError(w, http.StatusBadRequest, "Invalid Email Address")
		return
	}

	if !isValidPassword(pass) {
		writeJSONError(w, http.StatusBadRequest, "Invalid Password")
		return
	}

	userId, err := saveUser(email, pass)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	payload, _ := json.Marshal(RegistrationResp{UserID: *userId})
	w.WriteHeader(http.StatusCreated)
	w.Write(payload)
}
