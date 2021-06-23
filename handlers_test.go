package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

const emptyCommand = "TRUNCATE TABLE USERS"

var regTests = []struct {
	Username     string
	Password     string
	ExpectedCode int
}{
	{"", "", 400},
	{"teramakuro", "zeroma", http.StatusBadRequest},
	{"abc@def.com", "zeroma", http.StatusBadRequest},
	{"abc@def@com", "zeromakuro", http.StatusBadRequest},
	{"me@example.com", "lamatakturun", http.StatusCreated},
	{"me@example.com", "zeromakuro", http.StatusBadRequest}, // second request
}

func TestRegistration(t *testing.T) {

	db, err := sql.Open("postgres", dbstr)

	if err != nil || db.Ping() != nil {
		t.Skip("SQL Server not online, skipping registration test")
	}

	_, err = db.Exec(emptyCommand)
	db.Close()

	if err != nil {
		t.Skipf("Error occured when testing registration: %s", err.Error())
	}

	for i, item := range regTests {
		formValues := url.Values{}
		formValues.Add("email", item.Username)
		formValues.Add("password", item.Password)

		payload := formValues.Encode()

		req, _ := http.NewRequest(http.MethodPost, "/register", strings.NewReader(payload))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("Content-Length", strconv.Itoa(len(payload)))

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(registrationHandler)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != item.ExpectedCode {
			t.Errorf("Handler returned wrong status code: got %v want %v.\nFor item %v username %s password %s",
				status, item.ExpectedCode, i, item.Username, item.Password)
		}
	}
}