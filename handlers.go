package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

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

func saveUser(user User) (*uuid.UUID, error) {

	db, err := sql.Open("postgres", dbstr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	defer db.Close()

	regCommand := "INSERT INTO users (id, email, passhash, active) VALUES ($1, $2, $3, $4)"

	if _, err := db.Exec(regCommand, user.ID, user.Email, user.PassHash, user.IsActive); err != nil {

		if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
			return nil, errors.New("User already exists")
		}

		return nil, err
	}

	return &user.ID, nil
}

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	pass := r.PostFormValue("password")

	newUser, err := GenerateUser(email, pass)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := saveUser(newUser)

	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	payload, _ := json.Marshal(RegistrationResp{UserID: *userId})
	w.WriteHeader(http.StatusCreated)
	w.Write(payload)
}
