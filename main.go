package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"net/mail"
)

type JSONError struct {
	Message string `json:"message,omitempty"`
}

func main() {
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

		w.WriteHeader(http.StatusCreated)
	})

	http.ListenAndServe(":1337", r)
}
