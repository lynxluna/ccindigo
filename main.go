package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/google/uuid"
	"net/http"
)

type JSONError struct {
	Message string `json:"message,omitempty"`
}

type RegistrationResp struct {
	UserID uuid.UUID `json:"user_id,omitempty"`
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/register", registrationHandler)

	http.ListenAndServe(":1337", r)
}
