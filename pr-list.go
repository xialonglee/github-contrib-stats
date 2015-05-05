package main

import (
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// ... how to insert that parameters from external configuration files?
const accessToken = "--- access token ---"
const orgName = "--- organization name ---"

type tokenSource struct {
	token *oauth2.Token
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

// create a github client only once.
// call getClient() instead of client() or create client each time.
func client() func() *github.Client {
	ts := &tokenSource{
		&oauth2.Token{AccessToken: accessToken},
	}

	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	return func() *github.Client {
		return client
	}
}

var getClient = client()

func orgRepositoriesList() []github.Repository {
	opt := &github.RepositoryListByOrgOptions{
		Type: "all",
		ListOptions: github.ListOptions{
			PerPage: 50,
		},
	}

	client := getClient()
	repos, _, err := client.Repositories.ListByOrg(orgName, opt)

	if err != nil {
		fmt.Println(err)
	}
	return repos
}

func pullRequestsList(repo github.Repository) []github.PullRequest {
	opt := &github.PullRequestListOptions{
		// State: ,
		// Head: ,
		// Base: ,
		// Sort: ,
		// Direction: ,
		ListOptions: github.ListOptions{
			PerPage: 10,
		},
	}

	client := getClient()
	pr, _, err := client.PullRequests.List(*repo.Owner.Login, *repo.Name, opt)

	if err != nil {
		fmt.Println(err)
	}

	return pr
}

func main() {
	repos := orgRepositoriesList()

	for i := 0; i < len(repos); i++ {
		fmt.Println("---------------------")
		fmt.Println(*repos[i].Name)
		fmt.Println("---------------------")

		pulls := pullRequestsList(repos[i])

		for j := 0; j < len(pulls); j++ {
			fmt.Println(*pulls[j].Title)
		}

	}
}
