package main

import (
	"errors"
	"math"
	"sort"
	"strconv"
	"strings"
)

type PullRequestEvent struct {
	Repository  *Repository `json:"repository"`
	Actor       *User       `json:"actor"`
	PullRequest *PullRequest
}

type PullRequest struct {
	Id          int              `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	State       PullRequestState `json:"state"`
	Author      *Author          `json:"author"`
	Source      *PullRequestRef  `json:"source"`
	Destination *PullRequestRef  `json:"destination"`
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
	Uuid    string   `json:"uuid"`
	Name    string   `json:"name"`
	Links   Links    `json:"links"`
	Project *Project `json:"project"`
	Owner   *Owner   `json:"owner"`
}

type Project struct {
	Name  string          `json:"name"`
	Links map[string]Link `json:"links"`
}

type Links struct {
	Self  *Link
	Clone []*Link `json:"clone,omitempty"`
}

type Link struct {
	Name string `json:"name,omitempty"`
	Href string `json:"href"`
}

type Author struct {
	Raw   string `json:"raw"`
	User  *User  `json:"user,omitempty"`
	Name  string
	Email string
}

type Owner struct {
	UUID string `json:"uuid"`
}

type User struct {
	UUID  string          `json:"uuid"`
	Links map[string]Link `json:"links"`
}

type CascadeMergeState struct {
	Source string
	Target string
	error
}

type Cascade struct {
	Branches []string
	Current  int
}

type CascadeOptions struct {
	DevelopmentName string
	ReleasePrefix   string
}

// It returns the next branch in the cascade or an empty string if it reached the end.
func (c *Cascade) Next() string {
	if len(c.Branches) > c.Current+1 {
		c.Current += 1
		return c.Branches[c.Current]
	}
	return ""
}

// Add a branch to the cascade and sort branches. If the cascade already contains a branch named identically,
// the cascade will remain unmodified.
func (c *Cascade) Append(branchName string) {
	for _, b := range c.Branches {
		if b == branchName {
			return
		}
	}
	c.Branches = append(c.Branches, branchName)
	sort.Sort(ByVersion(c.Branches))
}

// Slice cascade branches to have only the target branch and its following branches.
func (c *Cascade) Slice(startBranch string) {
	for _, branch := range c.Branches {
		if branch != startBranch {
			c.Branches = c.Branches[1:]
		} else {
			break
		}
	}
}

// Extract an int representation of the version found in the given branch. Branch must be named accordingly to the
// following format :
//     <kind>/<version>
// The part following the slash must be an int.
// It returns the version or MaxInt32 if it comply to the format.
func extractVersion(branch string) int {
	parts := strings.Split(branch, "/")
	if len(parts) > 0 {
		version, err := strconv.Atoi(parts[len(parts)-1])
		if err == nil {
			return version
		}
	}
	return math.MaxInt32
}

type ByVersion []string

func (b ByVersion) Len() int {
	return len(b)
}

func (b ByVersion) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByVersion) Less(i, j int) bool {
	return extractVersion(b[i]) < extractVersion(b[j])
}

func (r *Repository) URL(protocols ...string) (string, error) {
	links := r.Links.Clone
	if links == nil {
		return "", errors.New("missing clone link")
	}

	for _, cloneLink := range links {
		for _, p := range protocols {
			if len(p) == 0 || p == cloneLink.Name {
				return cloneLink.Href, nil
			}
		}
	}

	return "", errors.New("no matching clone link")
}
