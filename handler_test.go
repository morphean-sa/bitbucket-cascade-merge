package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// The body is a pull request event. Happy path. We expect a status 201
func TestEventHandler_HandleValidPayload(t *testing.T) {
	c := make(chan PullRequestEvent, 1)
	eh := EventHandler{channel: c}

	rr, err := request("test/fixtures/hook-pull-request-fulfilled.json", eh.Handle())
	if err != nil {
		t.Fatal(err)
	}

	event := <-c
	if &event.Repository == nil {
		t.Error("repository must be defined")
	}

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

// The body is a pull request event but state is not set to MERGED. We expect a status 422
func TestEventHandler_HandleUnsupportedState(t *testing.T) {
	c := make(chan PullRequestEvent, 1)
	eh := EventHandler{channel: c}

	rr, err := request("test/fixtures/hook-pull-request-created.json", eh.Handle())
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnprocessableEntity)
	}
}

// The body is a push event (not a pull request event). We expect a status 400
func TestEventHandler_HandleUnsupportedEventType(t *testing.T) {
	c := make(chan PullRequestEvent, 1)
	eh := EventHandler{channel: c}

	rr, err := request("test/fixtures/hook-pr-merged-develop.json", eh.Handle())
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

// The body contains some random model that we cannot deserialize. We expect a status 400.
func TestEventHandler_HandleBadRequest(t *testing.T) {
	c := make(chan PullRequestEvent, 1)
	eh := EventHandler{channel: c}

	rr, err := request("test/fixtures/hook-bad-request.json", eh.Handle())
	if err != nil {
		t.Fatal(err)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func request(filename string, hf http.Handler) (*httptest.ResponseRecorder, error) {
	file, _ := os.Open(filename)
	req, err := http.NewRequest("POST", "/hook", file)
	if err != nil {
		return nil, err
	}

	rr := httptest.NewRecorder()
	hf.ServeHTTP(rr, req)

	return rr, nil
}
