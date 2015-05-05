package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var config Config

type Config struct {
	AccessToken string
	OrgName     string
}

// read config file
var _ = func() int {
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(file, &config)
	return 0
}()

type tokenSource struct {
	token *oauth2.Token
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

// create a github client only once.
// call client() and create client only once.
var client = func() func() *github.Client {
	ts := &tokenSource{
		&oauth2.Token{AccessToken: config.AccessToken},
	}

	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	return func() *github.Client {
		return client
	}
}()

func orgRepositoriesList() []github.Repository {
	opt := &github.RepositoryListByOrgOptions{
		Type: "all",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	client := client()
	repos, _, err := client.Repositories.ListByOrg(config.OrgName, opt)

	if err != nil {
		fmt.Println(err)
	}
	return repos
}

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

func statPullRequests() []ReposStat {
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

func formatOutput(stats []ReposStat) {
	fmt.Print("\n")
	for i := 0; i < len(stats); i++ {
		fmt.Println("------------------------------------------")
		fmt.Printf(" repo   : %s\n", stats[i].Name)
		fmt.Printf(" open   : %d\n", stats[i].OpenPullRequest)
		fmt.Printf(" closed : %d\n", stats[i].ClosedPullRequest)
	}
	fmt.Println("------------------------------------------")
}

func main() {
	stats := statPullRequests()
	formatOutput(stats)
}
