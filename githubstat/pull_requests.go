package githubstat

import (
	"fmt"

	"github.com/google/go-github/github"
)

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

	client := client()
	pr, _, err := client.PullRequests.List(*repo.Owner.Login, *repo.Name, opt)

	if err != nil {
		fmt.Println(err)
	}

	return pr
}

type ReposStat struct {
	Name              string
	ClosedPullRequest int
	OpenPullRequest   int
}

func StatPullRequests() []ReposStat {
	repos := orgRepositoriesList()
	stats := make([]ReposStat, len(repos), len(repos))

	for i := 0; i < len(repos); i++ {
		stats[i].Name = *repos[i].Name

		closed_pr := 0
		cont := true
		for page := 1; cont; page++ {
			fmt.Print(".")
			pulls := pullRequestsList(repos[i], page)

			if len(pulls) > 0 {
				closed_pr += len(pulls)
			} else {
				cont = false
			}
		}

		stats[i].ClosedPullRequest = closed_pr
	}

	return stats
}
