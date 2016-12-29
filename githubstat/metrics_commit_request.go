package githubstat

import (
	"fmt"
	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
	"os"
	"strconv"
	"strings"
	"time"
)

type WarpRepositoryCommit struct {
	RepositoryCommit *github.RepositoryCommit
	Owner            string
	Repo             string
	MergedAt         *time.Time
}

// sameCommitter means the pr author was same with the commit author
func (m *WarpRepositoryCommit) findMergedTime(client *github.Client, sameCommitter bool) bool {
	SHA := m.RepositoryCommit.SHA
	owner := m.Owner
	repo := m.Repo
	author := m.RepositoryCommit.Author.Login
	if !sameCommitter {
		author = m.RepositoryCommit.Committer.Login
	}
	// TODO consider using template
	query := *SHA + " repo:" + owner + "/" + repo + " type:pr" + " author:" + *author
	pr := findPRfromCommit(client, query)
	if pr == nil {
		return false
	}
	m.MergedAt = pr.ClosedAt
	return true
}

type OverallCommitMetrics struct {
	Overall []*CommitMetrics
}

func (m *OverallCommitMetrics) Show() {
	m.merge()
	data := [][]string{}
	for _, metrics := range m.Overall {
		r := []string{metrics.User, strconv.Itoa(metrics.Commits), strconv.Itoa(metrics.FilteredCommits)}
		data = append(data, r)
	}
	if len(data) != 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"User Name", "Commits", "stackalytics style"})

		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output
	}
}

func (m *OverallCommitMetrics) merge() {
	// user name to slice index of the first occurence of user's metrics
	mapping := make(map[string]int)
	var merged []*CommitMetrics
	for _, metrics := range m.Overall {
		if i, found := mapping[metrics.User]; found {
			cm := merged[i]
			cm.Commits += metrics.Commits
			cm.FilteredCommits += metrics.FilteredCommits
		} else {
			mapping[metrics.User] = len(merged)
			merged = append(merged, metrics)
		}
	}
	m.Overall = merged
}

type CommitMetrics struct {
	User            string
	Commits         int
	FilteredCommits int // commits filtered by some rules, equal or less than Commits
}

type CommitMetricsRequest struct {
	prReq PullRequestMetricsRequest
}

func (cr *CommitMetricsRequest) express() {
	fmt.Println("metrics: commit request stat analysis")
}

func (cr *CommitMetricsRequest) validate() bool {
	return cr.prReq.validate()
}

func (cr *CommitMetricsRequest) SetParameters(param *MetricsParameters) {
	cr.prReq.SetParameters(param)
}

func (cr *CommitMetricsRequest) FetchMetrics() Metrics {
	cr.express()

	proxyClient := &ProxyClient{}
	client := proxyClient.getClient()

	if cr.validate() {

		m := cr.prReq

		var metrics OverallCommitMetrics
		cr.prReq.expandRepos(client)

		for _, user := range Config.Users {
			fmt.Printf("counts commits of %s\n", user)

			for _, repo := range m.param.Repos {
				ownerName := *repo.OwnerName
				repoName := *repo.RepoName
				fmt.Printf("%s/%s : listing commits\n", ownerName, repoName)

				commits, err := listCommits(client, ownerName, repoName, user.Name)
				if err != nil {
					panic(err)
				}

				if metrics.Overall == nil {
					metrics.Overall = []*CommitMetrics{}
				}

				lenCommits := len(commits)
				lenFilteredCommits := len(filterCommits(commits))
				fmt.Printf("User: %s, commits: %d, filterCommits: %d\n",
					user, lenCommits, lenFilteredCommits)

				metrics.Overall = append(metrics.Overall, &CommitMetrics{
					User:            user.Name,
					Commits:         lenCommits,
					FilteredCommits: lenFilteredCommits,
				})
			}

		}
		return &metrics
	}
	return &OverallCommitMetrics{
		Overall: []*CommitMetrics{
			{
				User:            "",
				Commits:         -1,
				FilteredCommits: -1,
			},
		},
	}

}

func listCommits(client *github.Client, owner string, repo string, author string) ([]*github.RepositoryCommit, error) {
	opt := &github.CommitsListOptions{
		Author:      author,
		ListOptions: github.ListOptions{PerPage: 100, Page: 1},
		//Head:        client.userName + ":",
	}

	page := 1
	var allCommits []*github.RepositoryCommit
	for {
		commits, resp, err := client.Repositories.ListCommits(owner, repo, opt)
		if err != nil {
			return nil, err
		}

		allCommits = append(allCommits, commits...)
		fmt.Printf("page:%d fin\n", page)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
		page++
	}

	return allCommits, nil
}

func filterCommits(commits []*github.RepositoryCommit) []*github.RepositoryCommit {
	var filteredCommits []*github.RepositoryCommit
	prefixToFiltered := "Merge branch 'master' into"
	for _, c := range commits {
		if strings.HasPrefix(*c.Commit.Message, prefixToFiltered) {
			continue
		}
		filteredCommits = append(filteredCommits, c)
	}
	return filteredCommits
}

// findPRfromCommit finds a pull request from SHA of commit
// Note: every pull request is an issue, but not every issue is a pull request.
// make sure string "+type:pr" was included in query string
func findPRfromCommit(client *github.Client, query string) *github.Issue {
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 100,
			// 100 is so big, it should not return so much results
			// so no need to handle the "NextPage" attr of the response.
		},
	}
	results, _, err := client.Search.Issues(query, opt)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return nil
	} else if *results.Total > 1 {
		fmt.Printf("warning: find multiple pull requests from a commit's SHA, we take the first one.\n")
		fmt.Printf("query string is %s \n", query)
		// TODO to sort the issues according the pulls num and earliest pr is the first item.
		return &results.Issues[0]
	} else if len(results.Issues) == 0 {
		fmt.Printf("warning: find no pull requests from a commit's SHA, the commit should come from commiter.\n")
		fmt.Printf("query string is %s \n", query)
		return nil
	}
	return &results.Issues[0]
}

func getStackalyticsCommits(client *github.Client, owner string, repo string, author string) []*WarpRepositoryCommit {
	fmt.Printf("%s/%s : listing commits of stackalytics.com style\n", owner, repo)
	overallCommits, err := listCommits(client, owner, repo, author)
	if err != nil {
		panic(err)
	}
	overallCommits = filterCommits(overallCommits)

	var warpOverallCommits []*WarpRepositoryCommit
	for _, commit := range overallCommits {
		warpCommit := &WarpRepositoryCommit{commit, owner, repo, &time.Time{}}
		if !warpCommit.findMergedTime(client, true) {
			fmt.Printf("could not find pull request which include this commit:%s \n", *commit.Commit.Message)
			fmt.Printf("change author to the the committer:%s \n", *commit.Committer.Login)
			warpCommit.findMergedTime(client, false)
		}
		warpOverallCommits = append(warpOverallCommits, warpCommit)
	}
	return warpOverallCommits
}
