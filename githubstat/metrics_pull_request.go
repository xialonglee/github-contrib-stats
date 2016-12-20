package githubstat

import (
	"fmt"
	"os"

	"strings"

	"strconv"

	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
)

type OverallPullRequestMetrics struct {
	Overall []*PullRequestMetrics
}
type PullRequestMetrics struct {
	User          string
	Merged        int // already merged PRs, PRs of this kind are also closed
	MergedCommits int // the sum of commits number in merged PRs
	LGTMed        int // open PRs with LGTM label
	NonLGTMed     int //open PRs without LGTM label
}

func (m *OverallPullRequestMetrics) merge() {

	// user name to slice index of the first occurence of user's metrics
	mapping := make(map[string]int)
	var merged []*PullRequestMetrics
	for _, metrics := range m.Overall {
		if i, found := mapping[metrics.User]; found {
			prm := merged[i]
			prm.Merged += metrics.Merged
			prm.MergedCommits += metrics.MergedCommits
			prm.LGTMed += metrics.LGTMed
			prm.NonLGTMed += metrics.NonLGTMed
		} else {
			mapping[metrics.User] = len(merged)
			merged = append(merged, metrics)
		}
	}
	m.Overall = merged

}
func (m *OverallPullRequestMetrics) Show() {
	m.merge()
	data := [][]string{}
	for _, metrics := range m.Overall {
		r := []string{metrics.User, strconv.Itoa(metrics.Merged), strconv.Itoa(metrics.MergedCommits), strconv.Itoa(metrics.LGTMed), strconv.Itoa(metrics.NonLGTMed)}
		data = append(data, r)
	}
	if len(data) != 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"User Name", "Merged PRs", "Merged Commits", "LGTM'ed PRs", "NonLGTM'ed PRs"})

		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output
	}

}

type PullRequestMetricsRequest struct {
	param *MetricsParameters
}

func (m *PullRequestMetricsRequest) express() {
	//fmt.Printf("target repository: %s/%s\n", *m.param.OwnerName, *m.param.RepoName)
	fmt.Println("metrics: pull request stat analysis")
}

func (m *PullRequestMetricsRequest) SetParameters(param *MetricsParameters) {
	m.param = param
}

func (m *PullRequestMetricsRequest) validate() bool {
	for _, repo := range m.param.Repos {
		if *repo.OwnerName == "" {
			return false
		} else if *repo.RepoName == "" {
			return false
		}
	}

	return true
}
func getPullRequestCommits(client *github.Client, owner string, repo string, number int) (int, error) {

	pr, _, err := client.PullRequests.Get(owner, repo, number)
	if err != nil {
		return -1, err
	}
	return *pr.Commits, nil
}
func getPullRequest(client *github.Client, owner string, repo string, number int) (*github.PullRequest, error) {

	pr, _, err := client.PullRequests.Get(owner, repo, number)
	if err != nil {
		return nil, err
	}
	return pr, nil
}
func listPullRequests(client *github.Client, owner string, repo string, opt *github.PullRequestListOptions) ([]*github.PullRequest, error) {
	var allPRs []*github.PullRequest
	page := 1
	for {
		prs, resp, err := client.PullRequests.List(owner, repo, opt)
		if err != nil {
			return nil, err
		}
		allPRs = append(allPRs, prs...)
		fmt.Printf("page:%d fin\n", page)
		if resp.NextPage == 0 || prs[len(prs)-1].CreatedAt.Before(Config.StatBeginTime) {
			break
		}
		opt.ListOptions.Page = resp.NextPage
		page++
	}
	return allPRs, nil
}

func listRepositories(client *github.Client, owner string, opt *github.RepositoryListOptions) ([]*github.Repository, error) {
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(owner, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
		fmt.Printf("page:%d fin\n", resp.NextPage-1)
	}
	return allRepos, nil
}
func listOpenPullRequests(client *github.Client, owner string, repo string) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		//Head:        client.userName + ":",
	}
	return listPullRequests(client, owner, repo, opt)
}
func listClosedPullRequests(client *github.Client, owner string, repo string) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		State:       "closed",
		//Head:        client.userName + ":",
	}
	return listPullRequests(client, owner, repo, opt)
}
func getIssue(client *github.Client, owner string, repo string, number int) *github.Issue {
	issue, _, err := client.Issues.Get(owner, repo, number)
	if err != nil {
		panic(err)
	}
	return issue
}
func getPullRequestLabelNames(client *github.Client, owner string, repo string, number int) []string {
	issue := getIssue(client, owner, repo, number)
	var labelNames []string
	if issue.Labels != nil {
		for _, l := range issue.Labels {
			labelNames = append(labelNames, *l.Name)
		}
	}
	return labelNames

}
func isLGTMed(client *github.Client, owner string, repo string, number int) bool {
	lnames := getPullRequestLabelNames(client, owner, repo, number)
	if StringSliceContainsFold(lnames, "lgtm") {
		return true
	}
	return false
}
func StringSliceContainsFold(s []string, str string) bool {
	str = strings.ToUpper(str)
	for _, e := range s {
		if strings.ToUpper(e) == str {
			return true
		}
	}
	return false
}
func pullRequestOwnedBy(pr *github.PullRequest, userName string) bool {
	if pr != nil && pr.User != nil &&
		pr.User.Login != nil && *pr.User.Login == userName {
		return true
	}
	return false
}
func filterByUserName(prs []*github.PullRequest, userName string) []*github.PullRequest {
	var filtered []*github.PullRequest
	for _, pr := range prs {
		if pullRequestOwnedBy(pr, userName) {
			filtered = append(filtered, pr)
		}
	}
	return filtered
}
func (m *PullRequestMetricsRequest) expandRepos(client *github.Client) {
	var expanded []*RepoParameters
	for _, repo := range m.param.Repos {
		ownerName := *repo.OwnerName
		repoName := *repo.RepoName
		if repoName == "*" {
			repos, err := listRepositories(client, ownerName, &github.RepositoryListOptions{
				ListOptions: github.ListOptions{PerPage: 100}})

			if err != nil {
				panic(err)
			}
			for _, r := range repos {
				expanded = append(expanded, &RepoParameters{r.Owner.Login, r.Name})
			}

		} else {
			expanded = append(expanded, repo)
		}
	}

	m.param.Repos = expanded
}
func sumCommits(prs []*github.PullRequest) int {
	var sum int
	for _, pr := range prs {
		if pr.Commits != nil && *pr.Commits >= 0 {
			sum += *pr.Commits
		}
	}
	return sum
}
func (m *PullRequestMetricsRequest) FetchMetrics() Metrics {

	m.express()

	proxyClient := &ProxyClient{}
	client := proxyClient.getClient()
	if m.validate() {
		var metrics OverallPullRequestMetrics
		m.expandRepos(client)

		for _, repo := range m.param.Repos {
			ownerName := *repo.OwnerName
			repoName := *repo.RepoName
			fmt.Printf("%s/%s : listing open pull requests\n", ownerName, repoName)

			openPRs, err := listOpenPullRequests(client, ownerName, repoName)
			if err != nil {
				panic(err)
			}

			fmt.Printf("%s/%s : listing closed pull requests\n", ownerName, repoName)
			closedPRs, err := listClosedPullRequests(client, ownerName, repoName)
			if err != nil {
				panic(err)
			}

			for _, user := range Config.Users {
				var mergedPRs []*github.PullRequest
				var lgtmedPRs []*github.PullRequest
				var nonLGTMedPRs []*github.PullRequest
				filteredOpenPRs := filterByUserName(openPRs, user)
				filteredClosedPRs := filterByUserName(closedPRs, user)

				for _, pr := range filteredOpenPRs {
					if isLGTMed(client, ownerName, repoName, *pr.Number) {
						lgtmedPRs = append(lgtmedPRs, pr)
					} else {
						nonLGTMedPRs = append(nonLGTMedPRs, pr)
					}
				}

				for _, pr := range filteredClosedPRs {
					// Merged is always nil but MergedAt is not.
					if pr.MergedAt != nil {
						//get the specified pull request to fill in all other blank fields (such as Commits field)
						pr, err := getPullRequest(client, ownerName, repoName, *pr.Number)
						if err != nil {
							panic(err)
						}
						mergedPRs = append(mergedPRs, pr)
					}
				}
				if metrics.Overall == nil {
					metrics.Overall = []*PullRequestMetrics{}
				}

				lenMergedPRs := len(mergedPRs)
				lenLGTMedPRs := len(lgtmedPRs)
				lenNonLGTMed := len(nonLGTMedPRs)
				fmt.Printf("User: %s, Merged: %d, LGTM'ed: %d, NonLGTM'ed: %d \n",
					user, lenMergedPRs, lenLGTMedPRs, lenNonLGTMed)

				metrics.Overall = append(metrics.Overall, &PullRequestMetrics{
					User:          user,
					Merged:        lenMergedPRs,
					MergedCommits: sumCommits(mergedPRs),
					LGTMed:        lenLGTMedPRs,
					NonLGTMed:     lenNonLGTMed,
				})

			}
		}

		return &metrics
	}
	return &OverallPullRequestMetrics{
		[]*PullRequestMetrics{
			&PullRequestMetrics{
				User:          "",
				Merged:        -1,
				MergedCommits: -1,
				LGTMed:        -1,
				NonLGTMed:     -1,
			},
		},
	}

}
