package main

import (
	"flag"
	"strings"

	"./githubstat"
)

func main() {
	// repository name
	ownerName := flag.String("owner", "", "e.g. yshnb")
	repoName := flag.String("repository", "", "e.g. github-stat-script")
	target_metrics := flag.String("metrics", "", "available metrics: (star)")

	flag.Parse()

	var metricsRequest githubstat.MetricsRequest

	switch *target_metrics {
	case "issue":
		metricsRequest = &githubstat.IssueMetricsRequest{}
	case "star":
		metricsRequest = &githubstat.StarMetricsRequest{}
	case "pull_request":
		metricsRequest = &githubstat.PullRequestMetricsRequest{}
	default:
		metricsRequest = &githubstat.DefaultMetricsRequest{}
	}

	parameters := flag.Args()
	if len(parameters) > 0 {
		repo := strings.Split(parameters[0], "/")
		if *ownerName == "" && repo[0] != "" {
			ownerName = &repo[0]
		}
		if *repoName == "" && repo[1] != "" {
			repoName = &repo[1]
		}
	}

	metricsRequest.SetParameters(&githubstat.MetricsParameters{
		OwnerName: ownerName,
		RepoName:  repoName,
	})
	metrics := metricsRequest.FetchMetrics()
	metrics.GetMetrics()
}
