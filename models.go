package main

type PullRequestEvent struct {
	Repository  *Repository `json:"repository"`
	Actor       *User       `json:"actor"`
	PullRequest *PullRequest
}

type PullRequest struct {
	Id          int               `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	State       *PullRequestState `json:"state"`
	Author      *Author           `json:"author"`
	Source      *PullRequestRef   `json:"source"`
	Destination *PullRequestRef   `json:"destination"`
}

type PullRequestState string

const (
	Merged PullRequestState = "MERGED"
)

type PullRequestRef struct {
	Branch     *PullRequestBranch     `json:"branch"`
	Commit     *PullRequestCommit     `json:"commit"`
	Repository *PullRequestRepository `json:"repository"`
}

type PullRequestBranch struct {
	Name string `json:"name"`
}

type PullRequestCommit struct {
	Hash  string          `json:"name"`
	Links map[string]Link `json:"links"`
}

type PullRequestRepository struct {
	Name     string          `json:"name"`
	Fullname string          `json:"full_name"`
	Uuid     string          `json:"uuid"`
	Links    map[string]Link `json:"links"`
}

type Repository struct {
	Uuid    string          `json:"uuid"`
	Name    string          `json:"name"`
	Links   map[string]Link `json:"links"`
	Project *Project        `json:"project"`
	Owner   *Owner          `json:"owner"`
}

type Project struct {
	Name  string          `json:"name"`
	Links map[string]Link `json:"links"`
}

type Link struct {
	Href string `json:"href"`
}

type Author struct {
	Raw  string `json:"raw"`
	User *User  `json:"user,omitempty"`
}

type Owner struct {
	Username string `json:"username"`
}

type User struct {
	DisplayName string          `json:"display_name"`
	Links       map[string]Link `json:"links"`
}
