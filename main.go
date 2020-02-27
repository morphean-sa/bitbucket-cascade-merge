package main

import (
	"fmt"
	"github.com/ktrysmt/go-bitbucket"
	"log"
	"time"
)

func main() {
	/*c := make(chan Event, 100) // will queue at maximum 100 requests
	go worker(c)

	handler := NewEventHandler(c)

	// TODO add middleware to validate origin (source ip, headers, ...)
	http.HandleFunc("/hook", handler.Handle)
	http.ListenAndServe(fmt.Sprintf(":%s", getEnv("PORT", "5000")), nil)*/
	/*

		branches := listBranches(Repository{
			Name:    "audit",
			Links:   nil,
			Project: nil,
			Owner:   &Owner{Username: "morphean-sa"},
		})
		if len(branches) > 0 {
			next := findNextBranch("release/47", branches)
			if next != nil {
				log.Println(next)
				// TODO checkout next branch
				// TODO merge source
				// TODO if no conflict, cascade merge until develop

			} else {
				// nothing to do
				// TODO log ?
			}
		}*/
}

func worker(event <-chan PullRequestEvent) {
	for e := range event {
		log.Printf("push event on: %s", e.Repository.Links["self"])
		time.Sleep(5 * time.Second)
	}
}

func listBranches(r Repository) []bitbucket.RepositoryBranch {
	client := bitbucket.NewBasicAuth("samcontesse", "")

	opt := &bitbucket.RepositoryBranchingModelOptions{
		Owner:    r.Owner.Username,
		RepoSlug: r.Name,
	}

	model, _ := client.Repositories.Repository.BranchingModel(opt)
	for _, bt := range model.Branch_Types {
		if bt.Kind == "release" {
			refs, _ := client.Repositories.Repository.ListBranches(&bitbucket.RepositoryBranchOptions{
				Owner:    r.Owner.Username,
				RepoSlug: r.Name,
				Query:    fmt.Sprintf("(name~\"%s\" or name~\"%s\")", bt.Prefix, model.Development.Name),
			})
			return refs.Branches
		}
	}
	return []bitbucket.RepositoryBranch{}
}
