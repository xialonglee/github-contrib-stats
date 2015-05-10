package githubstat

import (
	"fmt"

	"github.com/mgutz/ansi"
)

type PullRequestMetrics struct {
	pull_request int
}

func (m *PullRequestMetrics) GetMetrics() {
	fmt.Printf(ansi.Color("this metrics has not been implemented yet.", "red"))
}

type PullRequestMetricsRequest struct {
	param *MetricsParameters
}

func (m *PullRequestMetricsRequest) express() {
	fmt.Printf("target repository: %s/%s\n", *m.param.OwnerName, *m.param.RepoName)
	fmt.Println("metrics: pull request count")
}

func (m *PullRequestMetricsRequest) SetParameters(param *MetricsParameters) {
	m.param = param
}

func (m *PullRequestMetricsRequest) validate() bool {
	return true
}

func (m *PullRequestMetricsRequest) FetchMetrics() Metrics {
	return &PullRequestMetrics{
		pull_request: -1,
	}
}
