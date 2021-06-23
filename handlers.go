package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
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

func registrationHandler(w http.ResponseWriter, r *http.Request) {
	email := r.PostFormValue("email")
	pass := r.PostFormValue("password")

	if _, err := mail.ParseAddress(email); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Invalid Email Address")
		return
	}

	if len(pass) < 8 {
		writeJSONError(w, http.StatusBadRequest, "Invalid Password")
		return
	}
	db, err := sql.Open("postgres", dbstr)
	if err != nil || db.Ping() != nil {
		writeJSONError(w, http.StatusInternalServerError, "Unknown Errors")
		return
	}

	defer db.Close()

	newUserID := uuid.New()
	newPwHasher := sha256.New()
	newPwHasher.Write([]byte(pass))
	newPwHash := fmt.Sprintf("%x", newPwHasher.Sum(nil))

	regCommand := "INSERT INTO users (id, email, passhash) VALUES ($1, $2, $3)"

	if _, err := db.Exec(regCommand, newUserID, email, newPwHash); err != nil {

		if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
			writeJSONError(w, http.StatusBadRequest, "User already exists")
			return
		}

		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	payload, _ := json.Marshal(RegistrationResp{UserID: newUserID})
	w.WriteHeader(http.StatusCreated)
	w.Write(payload)
}
