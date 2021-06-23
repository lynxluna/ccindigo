package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"net/http"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/register", registrationHandler)

	http.ListenAndServe(":1337", r)
}
