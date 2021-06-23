package main

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type UserRepository interface {
	Save(user User) error
}

type UserRepositoryPostgreSQL struct {
	db *sql.DB
}

func (repo *UserRepositoryPostgreSQL) SaveUser(user User) error {

	db, err := sql.Open("postgres", dbstr)

	if err != nil {
		return err
	}

	err = db.Ping()

	if err != nil {
		return err
	}

	defer db.Close()

	regCommand := "INSERT INTO users (id, email, passhash, active) VALUES ($1, $2, $3, $4)"

	if _, err := db.Exec(regCommand, user.ID, user.Email, user.PassHash, user.IsActive); err != nil {

		if pqError, ok := err.(*pq.Error); ok && pqError.Code == "23505" {
			return errors.New("User already exists")
		}

		return err
	}

	return nil
}
