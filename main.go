package main

import (
	"flag"
	"fmt"

	"./githubstat"

	"github.com/mgutz/ansi"
)

func main() {
	// repository name
	repoName := flag.String("repository", "github-stat-script", "e.g. yshnb")
	ownerName := flag.String("owner", "yshnb", "")
	check_metrics := flag.String("metrics", "", "available metrics: (issue|star|pull_request)")

	flag.Parse()

	var metricsRequest githubstat.MetricsRequest

	switch *check_metrics {
	case "issue":
		//		metricsRequest = &githubstat.IssueMetricsRequest{}
	case "star":
		metricsRequest = &githubstat.StarMetricsRequest{}
	case "pull_request":
		//		metricsRequest = &githubstat.PullRequestMetricsRequest{}
	default:
		fmt.Println(ansi.Color("you must select available metrics at least one.", "red"))
		metricsRequest = &githubstat.StarMetricsRequest{}
	}

	metricsRequest.SetParameters(&githubstat.MetricsParameters{
		OwnerName: ownerName,
		RepoName:  repoName,
	})
	metrics := metricsRequest.FetchMetrics()
	metrics.GetMetrics()

	//	marshaling, _ := json.Marshal(metrics)
	//	fmt.Println(string(marshaling))
	//
	//	stats := githubstat.StatPullRequests()
	//	format := &githubstat.Format{}
	//	format.FormatOutput(stats)
}
