package githubstat

import "fmt"

type PullRequestMetrics struct {
	pull_request int
}

func (m *PullRequestMetrics) GetMetrics() {
}

type PullRequestMetricsRequest struct {
}

func (m *PullRequestMetricsRequest) express() {
	fmt.Println("pull request selected")
}

func (m *PullRequestMetricsRequest) SetParameters() {
}

func (m *PullRequestMetricsRequest) FetchMetrics() {
}
