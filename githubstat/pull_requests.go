package githubstat

import (
	"fmt"

	"github.com/google/go-github/github"
)

// Pull Request List
func pullRequestsList(repo github.Repository, page int) []github.PullRequest {
	opt := &github.PullRequestListOptions{
		State: "closed",
		// Head: ,
		// Base: ,
		// Sort: ,
		// Direction: ,
		ListOptions: github.ListOptions{
			PerPage: 100,
			Page:    page,
		},
	}

	proxyClient := &ProxyClient{}
	client := proxyClient.getClient()
	pr, _, err := client.PullRequests.List(*repo.Owner.Login, *repo.Name, opt)

	if err != nil {
		fmt.Println(err)
	}

	return pr
}
