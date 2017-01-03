package githubstat

import (
	"fmt"
	"os"
	"time"

	"strings"

	"strconv"

	"github.com/google/go-github/github"
	"github.com/olekukonko/tablewriter"
)

type AllPullRequestMetrics struct {
	*WeekPullRequestMetrics
	*OverallPullRequestMetrics
}

func (a *AllPullRequestMetrics) Show() {
	a.WeekPullRequestMetrics.Show()
	a.OverallPullRequestMetrics.Show()
}

type WeekPullRequestMetrics struct {
	Week []*PullRequestMetrics
}

func (w *WeekPullRequestMetrics) merge() {

	w.Week = merge(w.Week)

}
func (w *WeekPullRequestMetrics) Show() {
	if !Config.StatEndTime.IsZero() {
		fmt.Println("Week statistics is disabled because statEndTime is specified")
		return
	}
	w.merge()
	data := [][]string{}
	var totalMerged int
	var totalMergedCommits int
	var totalLGTMed int
	var totalNonLGTMed int
	var totalCreated int

	for _, metrics := range w.Week {
		r := []string{metrics.User, strconv.Itoa(metrics.Merged),
			strconv.Itoa(metrics.MergedCommits), strconv.Itoa(metrics.LGTMed),
			strconv.Itoa(metrics.NonLGTMed), strconv.Itoa(metrics.Created)}
		data = append(data, r)
		totalMerged += metrics.Merged
		totalMergedCommits += metrics.MergedCommits
		totalLGTMed += metrics.LGTMed
		totalNonLGTMed += metrics.NonLGTMed
		totalCreated += metrics.Created

	}
	if len(data) != 0 {
		table := tablewriter.NewWriter(os.Stdout)
		fmt.Printf("\nStatistics for this Week ( week first day : %v)\n", Config.ThisWeekFirstDay)
		table.SetHeader([]string{"User Name", "Merged PRs", "Merged Commits",
			"LGTM'ed PRs", "NonLGTM'ed PRs", "Created PRs"})
		table.AppendBulk(data)
		table.Append([]string{
			"Total",
			strconv.Itoa(totalMerged),
			strconv.Itoa(totalMergedCommits),
			strconv.Itoa(totalLGTMed),
			strconv.Itoa(totalNonLGTMed),
			strconv.Itoa(totalCreated)},
		)
		table.Render() // Send output
	}
}

type OverallPullRequestMetrics struct {
	Overall []*PullRequestMetrics
}
type PullRequestMetrics struct {
	User                  string
	Merged                int // already merged PRs, PRs of this kind are also closed
	MergedCommits         int // the sum of commits number in merged PRs
	DeviatedMergedCommits int // stackalytics.com also do statistics analysis on kubernetes commits, but shamefully incorrectly. for report purpose, we calculate these through stackalyticsDeviation
	LGTMed                int // open PRs with LGTM label
	NonLGTMed             int //open PRs without LGTM label
	Created               int // created PRs including all open PRs and all merged closed PRs
}

func (m *OverallPullRequestMetrics) merge() {

	m.Overall = merge(m.Overall)

}

func (m *OverallPullRequestMetrics) deviate() {
	for _, metrics := range m.Overall {
		metrics.DeviatedMergedCommits = metrics.MergedCommits + getStackalyticsDeviation(metrics.User)
	}
}
func (m *OverallPullRequestMetrics) Show() {
	m.merge()
	//m.deviate()
	data := [][]string{}
	var totalMerged int
	var totalMergedCommits int
	var totalDeviatedMergedCommits int
	var totalLGTMed int
	var totalNonLGTMed int
	for _, metrics := range m.Overall {
		var mergedCommits string
		mergedCommits = fmt.Sprintf("%d / %d", metrics.MergedCommits, metrics.DeviatedMergedCommits)
		r := []string{metrics.User, strconv.Itoa(metrics.Merged), mergedCommits, strconv.Itoa(metrics.LGTMed), strconv.Itoa(metrics.NonLGTMed)}
		data = append(data, r)
		totalMerged += metrics.Merged
		totalMergedCommits += metrics.MergedCommits
		totalDeviatedMergedCommits += metrics.DeviatedMergedCommits
		totalLGTMed += metrics.LGTMed
		totalNonLGTMed += metrics.NonLGTMed
	}
	if len(data) != 0 {
		table := tablewriter.NewWriter(os.Stdout)
		fmt.Printf("\nOverall Statistics ( %v ~ %v)\n", Config.StatBeginTime, time.Now())
		mergedCommitsHeader := "Merged Commits(actual/stack)"
		table.SetHeader([]string{"User Name", "Merged PRs", mergedCommitsHeader, "LGTM'ed PRs", "NonLGTM'ed PRs"})
		table.AppendBulk(data)
		table.Append([]string{
			"Total",
			strconv.Itoa(totalMerged),
			fmt.Sprintf("%d / %d", totalMergedCommits, totalDeviatedMergedCommits),
			strconv.Itoa(totalLGTMed),
			strconv.Itoa(totalNonLGTMed),
		})
		table.Render() // Send output
	}

}
func getStackalyticsDeviation(userName string) int {
	for _, u := range Config.Users {
		if u.Name == userName {
			return u.StackalyticsDeviation
		}
	}
	return 0
}
func merge(toBeMerged []*PullRequestMetrics) []*PullRequestMetrics {
	// user name to slice index of the first occurence of user's metrics
	mapping := make(map[string]int)
	var merged []*PullRequestMetrics
	for _, metrics := range toBeMerged {
		if i, found := mapping[metrics.User]; found {
			prm := merged[i]
			prm.Merged += metrics.Merged
			prm.MergedCommits += metrics.MergedCommits
			prm.LGTMed += metrics.LGTMed
			prm.NonLGTMed += metrics.NonLGTMed
			prm.Created += metrics.Created
			prm.DeviatedMergedCommits += metrics.DeviatedMergedCommits
		} else {
			mapping[metrics.User] = len(merged)
			merged = append(merged, metrics)
		}
	}

	return merged
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
		State:       "open",
		Sort:        "created",
		Direction:   "desc",
	}
	var allPRs []*github.PullRequest

	page := 1
loop:
	for {
		prs, resp, err := client.PullRequests.List(owner, repo, opt)
		if err != nil {
			return nil, err
		}

		fmt.Printf("page:%d fin\n", page)
		for _, pr := range prs {
			t := pr.CreatedAt
			if !Config.StatEndTime.IsZero() && !t.Before(Config.StatEndTime) {
				continue
			}
			if !t.Before(Config.StatBeginTime) {
				allPRs = append(allPRs, pr)
			} else {
				break loop
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
		page++
	}
	return allPRs, nil
}
func listClosedPullRequests(client *github.Client, owner string, repo string) ([]*github.PullRequest, error) {
	opt := &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
		State:       "closed",
		Sort:        "updated",
		Direction:   "desc",
	}

	var allPRs []*github.PullRequest

	page := 1

loop:
	for {
		prs, resp, err := client.PullRequests.List(owner, repo, opt)
		if err != nil {
			return nil, err
		}

		fmt.Printf("page:%d fin\n", page)
		for _, pr := range prs {

			if pr.MergedAt == nil {
				continue
			}
			t := pr.UpdatedAt
			/*
				MergedAt is always before UpdatedAt, so if a PR is updated before stat begin time,
				this PR is absolutely merged before stat begin time.
				WARNING: UpdatedAt is sorted descendingly, but MergedAt is not. so we can break outer loop according to MergedAt
			*/
			if !t.Before(Config.StatBeginTime) {

				if !pr.MergedAt.Before(Config.StatBeginTime) {
					if !Config.StatEndTime.IsZero() && pr.MergedAt.Before(Config.StatEndTime) {
						allPRs = append(allPRs, pr)
					}else if Config.StatEndTime.IsZero() {
						allPRs = append(allPRs, pr)
					}
				}
			} else {
				break loop
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.ListOptions.Page = resp.NextPage
		page++
	}

	return allPRs, nil
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
func getPullRequestLatestLGTMEvent(client *github.Client, owner string, repo string, number int) (*github.IssueEvent, error) {
	//var allEvents []*github.IssueEvent
	page := 1
	opt := &github.ListOptions{PerPage: 100}
	for {
		events, resp, err := client.Issues.ListIssueEvents(owner, repo, number, opt)
		if err != nil {
			return nil, err
		}
		for _, evt := range events {
			//fmt.Printf("event created at : %v", evt.CreatedAt)
			if *evt.Event == "labeled" && (strings.EqualFold(*evt.Label.Name, "LGTM") || strings.EqualFold(*evt.Label.Name, "Docs LGTM")) {
				return evt, nil
			}
		}
		//allEvents = append(allEvents, events...)
		fmt.Printf("page:%d fin\n", page)
		if resp.NextPage == 0 || events[len(events)-1].CreatedAt.Before(Config.StatBeginTime) {
			break
		}
		opt.Page = resp.NextPage
		page++
	}
	return nil, fmt.Errorf("no LGTM event found")

}
func isLGTMed(client *github.Client, owner string, repo string, number int) bool {
	lnames := getPullRequestLabelNames(client, owner, repo, number)
	if StringSliceContainsAnyFold(lnames, "LGTM", "Docs LGTM") {
		return true
	}
	return false
}
func StringSliceContainsAnyFold(s []string, str ...string) bool {
	if len(str) == 0 {
		return false
	}
	for _, elem := range str {
		for _, e := range s {
			if strings.ToUpper(e) == strings.ToUpper(elem) {
				return true
			}
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
func inThisWeek(t *time.Time) bool {
	if t.Year() >= Config.ThisWeekFirstDay.Year() &&
		t.Month() >= Config.ThisWeekFirstDay.Month() &&
		t.Day() >= Config.ThisWeekFirstDay.Day() {
		return true
	}
	return false
}
func (m *PullRequestMetricsRequest) FetchMetrics() Metrics {

	m.express()

	proxyClient := &ProxyClient{}
	client := proxyClient.getClient()
	if m.validate() {

		var metrics OverallPullRequestMetrics = OverallPullRequestMetrics{Overall: []*PullRequestMetrics{}}
		var weekMetrics WeekPullRequestMetrics = WeekPullRequestMetrics{Week: []*PullRequestMetrics{}}
		var all AllPullRequestMetrics = AllPullRequestMetrics{WeekPullRequestMetrics: &weekMetrics, OverallPullRequestMetrics: &metrics}
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
				var overallMergedPRs []*github.PullRequest
				var overallLGTMedPRs []*github.PullRequest
				var overallNonLGTMedPRs []*github.PullRequest
				var overallStackalyticsCommits []*WarpRepositoryCommit
				var weekMergedPRs []*github.PullRequest
				var weekLGTMedPRs []*github.PullRequest
				var weekNonLGTMedPRs []*github.PullRequest
				var weekCreatedPRs []*github.PullRequest
				var weekStackalyticsCommits []*WarpRepositoryCommit
				userName := user.Name

				overallStackalyticsCommits = getStackalyticsCommits(client, ownerName, repoName, userName)
				filteredOpenPRs := filterByUserName(openPRs, userName)
				filteredClosedPRs := filterByUserName(closedPRs, userName)

				for _, c := range overallStackalyticsCommits {
					if inThisWeek(c.MergedAt) {
						weekStackalyticsCommits = append(weekStackalyticsCommits, c)
					}
				}

				for _, pr := range filteredOpenPRs {
					if inThisWeek(pr.CreatedAt) {
						weekCreatedPRs = append(weekCreatedPRs, pr)
					}
					if isLGTMed(client, ownerName, repoName, *pr.Number) {
						overallLGTMedPRs = append(overallLGTMedPRs, pr)

						if event, err := getPullRequestLatestLGTMEvent(client, ownerName, repoName, *pr.Number); err != nil {
							panic(err)
						} else {

							if inThisWeek(event.CreatedAt) {
								weekLGTMedPRs = append(weekLGTMedPRs, pr)
							}
						}

					} else {
						overallNonLGTMedPRs = append(overallNonLGTMedPRs, pr)
						if inThisWeek(pr.CreatedAt) {
							weekNonLGTMedPRs = append(weekNonLGTMedPRs, pr)
						}
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
						overallMergedPRs = append(overallMergedPRs, pr)
						if inThisWeek(pr.MergedAt) {
							weekMergedPRs = append(weekMergedPRs, pr)
							//fmt.Printf("pr title: %s, \npr merged at :%v\n", *pr.Title, *pr.MergedAt)
						}
						if inThisWeek(pr.CreatedAt) {
							weekCreatedPRs = append(weekCreatedPRs, pr)
						}
					}
				}

				lenMergedPRs := len(overallMergedPRs)
				lenLGTMedPRs := len(overallLGTMedPRs)
				lenNonLGTMed := len(overallNonLGTMedPRs)
				lenStackCommits := len(overallStackalyticsCommits)

				//fmt.Printf("User: %s, Merged: %d, LGTM'ed: %d, NonLGTM'ed: %d \n",
				//	user, lenMergedPRs, lenLGTMedPRs, lenNonLGTMed)

				metrics.Overall = append(metrics.Overall, &PullRequestMetrics{
					User:                  userName,
					Merged:                lenMergedPRs,
					MergedCommits:         sumCommits(overallMergedPRs),
					DeviatedMergedCommits: lenStackCommits,
					LGTMed:                lenLGTMedPRs,
					NonLGTMed:             lenNonLGTMed,
					Created:               -1,
				})

				weekMetrics.Week = append(weekMetrics.Week, &PullRequestMetrics{
					User:                  userName,
					Merged:                len(weekMergedPRs),
					MergedCommits:         sumCommits(weekMergedPRs),
					DeviatedMergedCommits: len(weekStackalyticsCommits),
					LGTMed:                len(weekLGTMedPRs),
					NonLGTMed:             len(weekNonLGTMedPRs),
					Created:               len(weekCreatedPRs),
				})

			}
		}

		return &all
	}
	return &AllPullRequestMetrics{
		OverallPullRequestMetrics: &OverallPullRequestMetrics{
			[]*PullRequestMetrics{
				&PullRequestMetrics{
					User:                  "",
					Merged:                -1,
					MergedCommits:         -1,
					DeviatedMergedCommits: -1,
					LGTMed:                -1,
					NonLGTMed:             -1,
					Created:               -1,
				},
			},
		},
		WeekPullRequestMetrics: &WeekPullRequestMetrics{
			[]*PullRequestMetrics{
				&PullRequestMetrics{
					User:                  "",
					Merged:                -1,
					MergedCommits:         -1,
					DeviatedMergedCommits: -1,
					LGTMed:                -1,
					NonLGTMed:             -1,
					Created:               -1,
				},
			},
		}}

}
