package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"net/http"
	"net/mail"
)

type JSONError struct {
	Message string `json:"message,omitempty"`
}

type RegistrationResp struct {
	UserID uuid.UUID `json:"user_id,omitempty"`
}

func main() {
	dbstr := "user=postgres password=postgres dbname=indigo sslmode=disable"
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		email := r.PostFormValue("email")
		pass := r.PostFormValue("password")

		if _, err := mail.ParseAddress(email); err != nil {
			w.WriteHeader(http.StatusBadRequest)

			errJSON := JSONError{Message: "Invalid Email Address"}
			payload, _ := json.Marshal(errJSON)

			w.Write(payload)
			return
		}

		if len(pass) < 8 {
			w.WriteHeader(http.StatusBadRequest)

			errJSON := JSONError{Message: "Invalid Password"}
			payload, _ := json.Marshal(errJSON)

			w.Write(payload)
			return
		}
		db, err := sql.Open("postgres", dbstr)
		if err != nil || db.Ping() != nil {
			w.WriteHeader(http.StatusInternalServerError)

			errJSON := JSONError{Message: "Unknown Errors"}
			payload, _ := json.Marshal(errJSON)

			w.Write(payload)
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
				w.WriteHeader(http.StatusBadRequest)
				payload, _ := json.Marshal(JSONError{Message: "User already exists"})
				w.Write(payload)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)

			errJSON := JSONError{Message: err.Error()}
			payload, _ := json.Marshal(errJSON)

			pqError, _ := err.(*pq.Error)

			fmt.Println(pqError.Code)

			w.Write(payload)
			return
		}
		payload, _ := json.Marshal(RegistrationResp{UserID: newUserID})
		w.WriteHeader(http.StatusCreated)
		w.Write(payload)
	})

	http.ListenAndServe(":1337", r)
}
