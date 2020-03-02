package main

import (
	"fmt"
	"github.com/ktrysmt/go-bitbucket"
)

type Bitbucket struct {
	Client   *bitbucket.Client
	Owner    string
	RepoSlug string
}

func NewBitbucket(username, password, owner, repoSlug string) *Bitbucket {
	return &Bitbucket{
		Client:   bitbucket.NewBasicAuth(username, password),
		Owner:    owner,
		RepoSlug: repoSlug,
	}
}

func (c *Bitbucket) GetCloneURL(protocols ...string) (string, error) {
	opt := &bitbucket.RepositoryOptions{
		Owner:    c.Owner,
		RepoSlug: c.RepoSlug,
	}

	r, err := c.Client.Repositories.Repository.Get(opt)
	if err != nil {
		return "", err
	}

	cloneLinks := r.Links["clone"]
	if cloneLinks != nil {
		for _, v := range cloneLinks.([]interface{}) {
			vv := v.(map[string]interface{})
			href := vv["href"].(string)
			name := vv["name"].(string)

			// no given protocol, return the first available
			if len(protocols) == 0 {
				return href, nil
			}

			// try protocols in the given order
			for _, p := range protocols {
				if p == name {
					return href, nil
				}
			}

		}
	}

	return "", fmt.Errorf("cannot determine clone url of %s", r.Full_name)
}

func (c *Bitbucket) GetCascadeOptions(owner, repo string) (*CascadeOptions, error) {
	opt := &bitbucket.RepositoryBranchingModelOptions{
		Owner:    c.Owner,
		RepoSlug: c.RepoSlug,
	}

	model, err := c.Client.Repositories.Repository.BranchingModel(opt)
	if err != nil {
		return nil, err
	}

	for _, bt := range model.Branch_Types {
		if bt.Kind == "release" {
			return &CascadeOptions{
				DevelopmentName: model.Development.Name,
				ReleasePrefix:   bt.Prefix,
			}, nil
		}
	}

	return nil, fmt.Errorf("cannot inspect branching model on %s", repo)
}

func (c *Bitbucket) CreatePullRequest(title, description, sourceBranch, destinationBranch string, reviewers ...string) error {
	opt := &bitbucket.PullRequestsOptions{
		Owner:             c.Owner,
		RepoSlug:          c.RepoSlug,
		Title:             title,
		Description:       description,
		SourceBranch:      sourceBranch,
		DestinationBranch: destinationBranch,
		// TODO reviewers by uuid
	}

	_, err := c.Client.Repositories.PullRequests.Create(opt)
	if err != nil {
		return err
	}
	return nil
}
