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
			PerPage: 50,
		},
	}

	client := client()
	repos, _, err := client.Repositories.ListByOrg(config.OrgName, opt)

	if err != nil {
		fmt.Println(err)
	}
	return repos
}

func pullRequestsList(repo github.Repository) []github.PullRequest {
	opt := &github.PullRequestListOptions{
		// State: "closed",
		// Head: ,
		// Base: ,
		// Sort: ,
		// Direction: ,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	client := client()
	pr, _, err := client.PullRequests.List(*repo.Owner.Login, *repo.Name, opt)

	if err != nil {
		fmt.Println(err)
	}

	return pr
}

func main() {
	repos := orgRepositoriesList()

	for i := 0; i < len(repos); i++ {
		fmt.Println("|---------------------|")
		fmt.Printf("| %s\n", *repos[i].Name)
		fmt.Println("|---------------------|")

		pulls := pullRequestsList(repos[i])

		for j := 0; j < len(pulls); j++ {
			fmt.Printf("#%d - %s\n", *pulls[j].Number, *pulls[j].Title)
		}

	}
}
