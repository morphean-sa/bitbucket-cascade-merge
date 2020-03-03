package main

import (
	"encoding/json"
	"net/http"
)

type EventHandler struct {
	channel chan<- PullRequestEvent
}

func (e EventHandler) Handle(writer http.ResponseWriter, request *http.Request) {

	var event PullRequestEvent

	err := json.NewDecoder(request.Body).Decode(&event)

	if err != nil || event.PullRequest == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// take only merged state
	if event.PullRequest.State != Merged {
		writer.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	// notify the channel
	select {
	case e.channel <- event:
		writer.WriteHeader(http.StatusCreated)
	default:
		writer.WriteHeader(http.StatusTooManyRequests)
	}

}

func NewEventHandler(c chan PullRequestEvent) *EventHandler {
	return &EventHandler{channel: c}
}
