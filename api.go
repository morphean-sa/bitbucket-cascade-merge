package main

import (
	"encoding/json"
	"github.com/samcontesse/bitbucket-cascade-merge/bitbucket"
	"net/http"
)

type EventHandler struct {
	channel chan<- bitbucket.PullRequestEvent
}

func (e EventHandler) Handle(writer http.ResponseWriter, request *http.Request) {

	var event bitbucket.PullRequestEvent

	err := json.NewDecoder(request.Body).Decode(&event)

	if err != nil || event.PullRequest == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// take only merged state
	if *event.PullRequest.State != bitbucket.Merged {
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
