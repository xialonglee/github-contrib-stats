package githubstat

import "fmt"

type Result struct {
	Name              string
	ClosedPullRequest int
	OpenPullRequest   int
}

func (r *Result) StatPullRequests() []Result {
	repos := orgRepositoriesList()
	stats := make([]Result, len(repos), len(repos))

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
